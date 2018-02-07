package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/session"
	"github.com/chanxuehong/util"
	"github.com/chanxuehong/util/security"
	mchcore "github.com/chanxuehong/wechat.v2/mch/core"
	mchpay "github.com/chanxuehong/wechat.v2/mch/pay"
	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/chanxuehong/wechat.v2/oauth2"
	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/ext"
	"github.com/reechou/real-reading/models"
	"github.com/reechou/real-reading/proto"
)

const (
	ReadingPrefix = "/reading"
)

const (
	READING_URI_SIGN_UP    = "signup"
	READING_URI_ENROLL     = "enroll"
	READING_URI_GO_ENROLL  = "goenroll"
	READING_URI_PAY        = "pay"
	READING_URI_PAY_NOTIFY = "paynotify"
	READING_URI_SUCCESS    = "success"
)

// manager uri
const (
	READING_URI_TODAY_ORDER = "today_order"
)

type ShareTpl struct {
	Title string
	Img   string
}

type HandlerRequest struct {
	Method string
	Path   string
	Val    []byte
	Params []string
}

type ReadingHandler struct {
	l *Logic

	smsExt *ext.SMSNotifyExt
	gm     *GraduationManager

	sessionStorage *session.Storage
	oauth2Endpoint oauth2.Endpoint
	oauth2Client   *oauth2.Client
	mchClient      *mchcore.Client
}

func NewReadingHandler(l *Logic) *ReadingHandler {
	lh := &ReadingHandler{l: l}

	lh.smsExt = ext.NewSMSNotifyExt(lh.l.cfg)
	var err error
	lh.gm, err = NewGraduationManager(lh.l.cfg)
	if err != nil {
		holmes.Fatal("new graduation manager error: %v", err)
	}

	lh.sessionStorage = session.New(20*60, 60*60)
	lh.oauth2Endpoint = mpoauth2.NewEndpoint(
		lh.l.cfg.ReadingOauth.ReadingWxAppId,
		lh.l.cfg.ReadingOauth.ReadingWxAppSecret)
	lh.oauth2Client = &oauth2.Client{
		Endpoint: lh.oauth2Endpoint,
	}
	lh.mchClient = mchcore.NewClient(
		lh.l.cfg.ReadingOauth.ReadingWxAppId,
		lh.l.cfg.ReadingOauth.ReadingMchId,
		lh.l.cfg.ReadingOauth.ReadingMchApiKey,
		nil,
	)

	return lh
}

func (self *ReadingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rr, err := parseRequest(r)
	if err != nil {
		holmes.Error("parse request error: %v", err)
		writeRsp(w, &proto.Response{Code: proto.RESPONSE_ERR})
		return
	}

	if rr.Path == "" {
		return
	}

	if strings.HasSuffix(rr.Path, "txt") {
		http.ServeFile(w, r, self.l.cfg.ReadingOauth.MpVerifyDir+rr.Path)
		return
	}
	if strings.HasSuffix(rr.Path, "css") {
		http.ServeFile(w, r, "./views/css/"+rr.Path)
		return
	}
	if strings.HasSuffix(rr.Path, "js") {
		http.ServeFile(w, r, "./views/js/"+rr.Path)
		return
	}
	if strings.HasSuffix(rr.Path, "jpg") {
		http.ServeFile(w, r, "./tmp/"+rr.Path)
		return
	}
	if strings.HasPrefix(rr.Path, READING_COURSE_MANAGER_URI_PREFIX) {
		self.courseManagerHandle(rr, w, r)
		return
	}
	if strings.HasPrefix(rr.Path, READING_COURSE_URI_PREFIX) {
		self.courseHandle(rr, w, r)
		return
	}
	if strings.HasPrefix(rr.Path, REGISTER_URI_PREFIX) {
		self.registerHandle(rr, w, r)
		return
	}
	if strings.HasPrefix(rr.Path, DATA_STATISTICS_URI_PREFIX) {
		self.dataStatisticsHandle(rr, w, r)
		return
	}

	switch rr.Path {
	case READING_URI_SIGN_UP:
		self.readingSignup(rr, w, r)
	case READING_URI_ENROLL:
		self.readingEnroll(rr, w, r)
	case READING_URI_GO_ENROLL:
		self.readingGoEnroll(rr, w, r)
	case READING_URI_PAY:
		self.readingPay(rr, w, r)
	case READING_URI_PAY_NOTIFY:
		self.readingPayNotify(rr, w, r)
	case READING_URI_SUCCESS:
		self.readingSuccess(rr, w, r)
	case READING_URI_TODAY_ORDER:
		self.readingPayToday(rr, w, r)
	default:
		http.ServeFile(w, r, self.l.cfg.ReadingOauth.MpVerifyDir+rr.Path)
	}
}

func (self *ReadingHandler) readingSignup(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	course := &models.Course{
		CourseType: READING_COURSE_TYPE_GD,
	}
	has, err := models.GetCourseMaxNum(course)
	if err != nil {
		holmes.Error("get course max num error: %v", err)
		io.WriteString(w, MSG_ERROR_SYSTEM)
		return
	}
	if !has {
		holmes.Error("get course max num has none")
		io.WriteString(w, MSG_ERROR_COURSE_NOT_FOUND)
		return
	}

	readingUserInfo := &ReadingEnrollUserInfo{
		CourseType:     course.CourseType,
		CourseNum:      course.CourseNum,
		IndexStartTime: time.Unix(course.StartTime, 0).Format("2006年01月02日"),
		StartTime:      time.Unix(course.StartTime, 0).Format("2006.01"),
		EndTime:        time.Unix(course.EndTime, 0).Format("2006.01"),
	}

	renderView(w, "./views/reading_sign.html", readingUserInfo)
}

func (self *ReadingHandler) readingEnroll(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		io.WriteString(w, err.Error())
		holmes.Error("url parse query error: %v", err)
		return
	}

	code := queryValues.Get("code")
	if code == "" {
		state := string(rand.NewHex())
		redirectUrl := fmt.Sprintf("http://%s%s", r.Host, r.URL.String())
		AuthCodeURL := mpoauth2.AuthCodeURL(self.l.cfg.ReadingOauth.ReadingWxAppId,
			redirectUrl,
			self.l.cfg.ReadingOauth.ReadingOauth2ScopeUser, state)
		http.Redirect(w, r, AuthCodeURL, http.StatusFound)
		return
	}

	token, err := self.oauth2Client.ExchangeToken(code)
	if err != nil {
		//holmes.Error("exchange token error: %v", err)
		http.Redirect(w, r, fmt.Sprintf("http://%s%s", r.Host, r.URL.Path), http.StatusFound)
		return
	}

	userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		io.WriteString(w, MSG_ERROR_SYSTEM)
		holmes.Error("get user info error: %v", err)
		return
	}
	holmes.Debug("user info: %+v", userinfo)

	readingUserInfo := &ReadingEnrollUserInfo{
		NickName:  userinfo.Nickname,
		AvatarUrl: userinfo.HeadImageURL,
		OpenId:    token.OpenId,
	}

	// --- new logic ---
	user := &models.User{
		OpenId: token.OpenId,
	}
	has, err := models.GetUserFromOpenid(user)
	if err != nil {
		holmes.Error("get user from openid error: %v", err)
		io.WriteString(w, MSG_ERROR_SYSTEM)
		return
	}
	if has {
		if user.Name != userinfo.Nickname || user.AvatarUrl != userinfo.HeadImageURL {
			user.Name = userinfo.Nickname
			user.AvatarUrl = userinfo.HeadImageURL
			err = models.UpdateUserWxInfo(user)
			if err != nil {
				holmes.Error("update user wxinfo error: %v", err)
			}
		}
		courseList, err := models.GetUserCourseList(user.ID)
		if err != nil {
			holmes.Error("get user course list error: %v", err)
			io.WriteString(w, MSG_ERROR_SYSTEM)
			return
		}
		for _, v := range courseList {
			if v.Course.CourseType == READING_COURSE_TYPE_GD {
				readingUserInfo.CourseType = v.Course.CourseType
				readingUserInfo.CourseNum = v.Course.CourseNum
				readingUserInfo.StartTime = time.Unix(v.Course.StartTime, 0).Format("2006.01.02")
				readingUserInfo.EndTime = time.Unix(v.Course.EndTime, 0).Format("2006.01.02")
				if v.Course.StartTime <= time.Now().Unix() {
					readingUserInfo.IfCourseStart = 1
				}
				renderView(w, "./views/reading_sign_success.html", readingUserInfo)
				return
			}
		}
	}

	if !has {
		user.AppId = self.l.cfg.ReadingOauth.ReadingWxAppId
		user.Name = userinfo.Nickname
		user.AvatarUrl = userinfo.HeadImageURL
		err = models.CreateUser(user)
		if err != nil {
			holmes.Error("create user error: %v", err)
			return
		}
	}

	// check older version
	readingUser := &models.ReadingPay{
		OpenId: token.OpenId,
	}
	has, err = models.GetReadingPayWithStatus(readingUser)
	if err != nil {
		holmes.Error("get reading pay error: %v", err)
		io.WriteString(w, MSG_ERROR_SYSTEM)
		return
	}
	if has {
		course := &models.Course{
			CourseType: readingUser.Course,
			CourseNum:  readingUser.Number,
		}
		has, err = models.GetCourseFromCourse(course)
		if err != nil {
			holmes.Error("get course from course error: %v", err)
			return
		}
		userCourse := &models.UserCourse{
			UserId:     user.ID,
			CourseId:   course.ID,
			CourseType: READING_COURSE_TYPE_GD,
			Money:      readingUser.Money,
			Status:     readingUser.Status,
			PayTime:    readingUser.UpdatedAt,
		}
		err = models.CreateUserCourse(userCourse)
		if err != nil {
			holmes.Error("create user course error: %v", err)
			return
		}
		readingUserInfo.CourseType = course.CourseType
		readingUserInfo.CourseNum = course.CourseNum
		readingUserInfo.StartTime = time.Unix(course.StartTime, 0).Format("2006.01.02")
		readingUserInfo.EndTime = time.Unix(course.EndTime, 0).Format("2006.01.02")
		renderView(w, "./views/reading_sign_success.html", readingUserInfo)
		return
	}

	course := &models.Course{
		CourseType: READING_COURSE_TYPE_GD,
	}
	has, err = models.GetCourseMaxNum(course)
	if err != nil {
		holmes.Error("get course max num error: %v", err)
		io.WriteString(w, MSG_ERROR_SYSTEM)
		return
	}
	if !has {
		holmes.Error("get course max num has none")
		io.WriteString(w, MSG_ERROR_COURSE_NOT_FOUND)
		return
	}
	if user.RealName != "" {
		readingUserInfo.EnrollName = user.RealName
	}
	if user.Phone != "" {
		readingUserInfo.EnrollMobile = user.Phone
	}
	if user.Wechat != "" {
		readingUserInfo.EnrollWechat = user.Wechat
	}
	readingUserInfo.CourseType = course.CourseType
	readingUserInfo.CourseNum = course.CourseNum
	readingUserInfo.StartTime = time.Unix(course.StartTime, 0).Format("2006-01-02")
	readingUserInfo.EndTime = time.Unix(course.EndTime, 0).Format("2006-01-02")
	renderView(w, "./views/reading_enroll.html", readingUserInfo)
	// --- end ---

	//readingUser := &models.ReadingPay{
	//	OpenId: token.OpenId,
	//}
	//has, err = models.GetReadingPay(readingUser)
	//if err != nil {
	//	holmes.Error("get reading pay error: %v", err)
	//	return
	//}
	//if has {
	//	if readingUser.Name != userinfo.Nickname || readingUser.AvatarUrl != userinfo.HeadImageURL {
	//		readingUser.Name = userinfo.Nickname
	//		readingUser.AvatarUrl = userinfo.HeadImageURL
	//		err = models.UpdateReadingPayWxInfo(readingUser)
	//		if err != nil {
	//			holmes.Error("update reading wx info error: %v", err)
	//		}
	//	}
	//	if readingUser.Status == READING_COURSE_STATUS_PAIED {
	//		renderView(w, "./views/reading_sign_success.html", readingUserInfo)
	//		return
	//	}
	//	if readingUser.RealName != "" {
	//		readingUserInfo.EnrollName = readingUser.RealName
	//	}
	//	if readingUser.Phone != "" {
	//		readingUserInfo.EnrollMobile = readingUser.Phone
	//	}
	//	if readingUser.Wechat != "" {
	//		readingUserInfo.EnrollWechat = readingUser.Wechat
	//	}
	//}
	//
	//if !has {
	//	readingUser = &models.ReadingPay{
	//		OpenId:    token.OpenId,
	//		AppId:     self.l.cfg.ReadingOauth.ReadingWxAppId,
	//		Name:      userinfo.Nickname,
	//		AvatarUrl: userinfo.HeadImageURL,
	//		Course:    READING_COURSE_TYPE_GD,
	//	}
	//	err = models.CreateReadingPay(readingUser)
	//	if err != nil {
	//		holmes.Error("create reading pay error: %v", err)
	//		return
	//	}
	//}
	//
	//renderView(w, "./views/reading_enroll.html", readingUserInfo)
}

func (self *ReadingHandler) readingGoEnroll(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &proto.ReadingEnrollReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	// --- new logic ---
	user := &models.User{
		OpenId: req.OpenId,
	}
	has, err := models.GetUserFromOpenid(user)
	if err != nil {
		holmes.Error("get user from open error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	if !has {
		holmes.Error("cannot found this openid")
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	user.RealName = req.Realname
	user.Phone = req.Mobile
	user.Wechat = req.Weixin
	err = models.UpdateUserInfo(user)
	if err != nil {
		holmes.Error("update user info error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	// --- end ---

	//readingUser := &models.ReadingPay{
	//	OpenId: req.OpenId,
	//}
	//has, err := models.GetReadingPay(readingUser)
	//if err != nil {
	//	holmes.Error("get reading pay error: %v", err)
	//	return
	//}
	//if !has {
	//	holmes.Error("cannot found this openid")
	//	rsp.Code = proto.RESPONSE_ERR
	//	return
	//}
	//readingUser.RealName = req.Realname
	//readingUser.Phone = req.Mobile
	//readingUser.Wechat = req.Weixin
	//err = models.UpdateReadingPayUserInfo(readingUser)
	//if err != nil {
	//	holmes.Error("update reading pay userinfo error: %v", err)
	//	rsp.Code = proto.RESPONSE_ERR
	//	return
	//}
}

func (self *ReadingHandler) readingPay(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		io.WriteString(w, err.Error())
		holmes.Error("url parse query error: %v", err)
		return
	}
	openid := queryValues.Get("openid")

	// --- new logic ---
	user := &models.User{
		OpenId: openid,
	}
	has, err := models.GetUserFromOpenid(user)
	if err != nil {
		holmes.Error("get user from openid error: %v", err)
		io.WriteString(w, MSG_ERROR_SYSTEM)
		return
	}
	if !has {
		io.WriteString(w, MSG_ERROR_USER_NOT_FOUND)
		return
	}
	course := &models.Course{
		CourseType: READING_COURSE_TYPE_GD,
	}
	has, err = models.GetCourseMaxNum(course)
	if err != nil {
		holmes.Error("get course max num error: %v", err)
		io.WriteString(w, MSG_ERROR_SYSTEM)
		return
	}
	if !has {
		holmes.Error("get course max num has none")
		io.WriteString(w, MSG_ERROR_COURSE_NOT_FOUND)
		return
	}

	payMoney := course.Money
	if openid == "oaKrZwsAF6pRX6z3Qn_EhIZ3DG90" || openid == "oaKrZwotcenPmZyLKtMyoHZSTlaQ" {
		payMoney = 1
	}
	unifiedRsp, err := self.readingUnifiedOrder(
		user.ID,
		payMoney,
		user.OpenId,
		GetIPFromRequest(r),
		fmt.Sprintf("%s%s/%s", r.Host, ReadingPrefix, READING_URI_PAY_NOTIFY),
	)
	if err != nil {
		holmes.Error("reading unified order error: %v", err)
		if strings.Contains(err.Error(), "ORDERPAID") {
			// 订单已支付, 但未更新状态
			userCourse := &models.UserCourse{
				UserId:     user.ID,
				CourseId:   course.ID,
				CourseType: READING_COURSE_TYPE_GD,
				Money:      payMoney,
				Status:     READING_COURSE_STATUS_PAIED,
				PayTime:    time.Now().Unix(),
			}
			err = models.CreateUserCourse(userCourse)
			if err != nil {
				holmes.Error("create user course error: %v", err)
			}
			readingUserInfo := &ReadingEnrollUserInfo{
				NickName:   user.Name,
				AvatarUrl:  user.AvatarUrl,
				OpenId:     user.OpenId,
				CourseType: course.CourseType,
				CourseNum:  course.CourseNum,
				StartTime:  time.Unix(course.StartTime, 0).Format("2006.01.02"),
				EndTime:    time.Unix(course.EndTime, 0).Format("2006.01.02"),
			}
			renderView(w, "./views/reading_sign_success.html", readingUserInfo)
			return
		}
		return
	}

	readingUserInfo := &ReadingEnrollUserInfo{
		NickName:  user.Name,
		AvatarUrl: user.AvatarUrl,
		OpenId:    openid,
	}
	readingUserInfo.WxJsApiParams = WxJsApiParams{
		AppId:     self.l.cfg.ReadingOauth.ReadingWxAppId,
		TimeStamp: fmt.Sprintf("%d", time.Now().Unix()),
		NonceStr:  string(rand.NewHex()),
		Package:   fmt.Sprintf("prepay_id=%s", unifiedRsp.PrepayId),
		SignType:  "MD5",
	}
	readingUserInfo.WxJsApiParams.PaySign = mchcore.JsapiSign(
		readingUserInfo.WxJsApiParams.AppId,
		readingUserInfo.WxJsApiParams.TimeStamp,
		readingUserInfo.WxJsApiParams.NonceStr,
		readingUserInfo.WxJsApiParams.Package,
		readingUserInfo.WxJsApiParams.SignType,
		self.mchClient.ApiKey(),
	)
	renderView(w, "./views/reading_pay.html", readingUserInfo)
	// --- end ---

	//readingUser := &models.ReadingPay{
	//	OpenId: openid,
	//}
	//has, err := models.GetReadingPay(readingUser)
	//if err != nil {
	//	io.WriteString(w, "未找到你哦,请刷新重新登录")
	//	holmes.Error("get reading pay error: %v", err)
	//	return
	//}
	//if !has {
	//	io.WriteString(w, "未找到你哦,请刷新重新登录")
	//	return
	//}
	//
	//payMoney := READING_COURSE_GD_MONEY
	//if openid == "oaKrZwsAF6pRX6z3Qn_EhIZ3DG90" {
	//	payMoney = 1
	//}
	//unifiedRsp, err := self.readingUnifiedOrder(
	//	readingUser.ID,
	//	int64(payMoney),
	//	readingUser.OpenId,
	//	GetIPFromRequest(r),
	//	fmt.Sprintf("%s%s/%s", r.Host, ReadingPrefix, READING_URI_PAY_NOTIFY),
	//)
	//if err != nil {
	//	holmes.Error("reading unified order error: %v", err)
	//	if strings.Contains(err.Error(), "ORDERPAID") {
	//		// 订单已支付, 但未更新状态
	//		readingUser.Status = READING_COURSE_STATUS_PAIED
	//		readingUser.Number = self.l.cfg.NowCourseNumber
	//		err = models.UpdateReadingPayStatusFromOpenId(readingUser)
	//		if err != nil {
	//			holmes.Error("update reading pay status error: %v", err)
	//		}
	//		readingUserInfo := &ReadingEnrollUserInfo{
	//			NickName:  readingUser.Name,
	//			AvatarUrl: readingUser.AvatarUrl,
	//			OpenId:    readingUser.OpenId,
	//		}
	//		renderView(w, "./views/reading_sign_success.html", readingUserInfo)
	//		return
	//	}
	//	return
	//}
	//
	//readingUserInfo := &ReadingEnrollUserInfo{
	//	NickName:  readingUser.Name,
	//	AvatarUrl: readingUser.AvatarUrl,
	//	OpenId:    openid,
	//}
	//readingUserInfo.WxJsApiParams = WxJsApiParams{
	//	AppId:     self.l.cfg.ReadingOauth.ReadingWxAppId,
	//	TimeStamp: fmt.Sprintf("%d", time.Now().Unix()),
	//	NonceStr:  string(rand.NewHex()),
	//	Package:   fmt.Sprintf("prepay_id=%s", unifiedRsp.PrepayId),
	//	SignType:  "MD5",
	//}
	//readingUserInfo.WxJsApiParams.PaySign = mchcore.JsapiSign(
	//	readingUserInfo.WxJsApiParams.AppId,
	//	readingUserInfo.WxJsApiParams.TimeStamp,
	//	readingUserInfo.WxJsApiParams.NonceStr,
	//	readingUserInfo.WxJsApiParams.Package,
	//	readingUserInfo.WxJsApiParams.SignType,
	//	self.mchClient.ApiKey(),
	//)
	//renderView(w, "./views/reading_pay.html", readingUserInfo)
}

func (self *ReadingHandler) readingPayNotify(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	defer func() {
		io.WriteString(w, "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>")
	}()

	msg, err := util.DecodeXMLToMap(bytes.NewReader(rr.Val))
	if err != nil {
		holmes.Error("decode xml error: %v", err)
		return
	}
	holmes.Debug("reading pay notify msg: %v", msg)

	returnCode, ok := msg["return_code"]
	if returnCode == mchcore.ReturnCodeSuccess || !ok {
		haveAppId := msg["appid"]
		wantAppId := self.l.cfg.ReadingOauth.ReadingWxAppId
		if haveAppId != "" && wantAppId != "" && !security.SecureCompareString(haveAppId, wantAppId) {
			err = fmt.Errorf("appid mismatch, have: %s, want: %s", haveAppId, wantAppId)
			holmes.Error("%v", err)
			return
		}

		haveMchId := msg["mch_id"]
		wantMchId := self.l.cfg.ReadingOauth.ReadingMchId
		if haveMchId != "" && wantMchId != "" && !security.SecureCompareString(haveMchId, wantMchId) {
			err = fmt.Errorf("mch_id mismatch, have: %s, want: %s", haveMchId, wantMchId)
			holmes.Error("%v", err)
			return
		}

		// 认证签名
		haveSignature, ok := msg["sign"]
		if !ok {
			holmes.Error("msg sign not found")
			return
		}
		wantSignature := mchcore.Sign(msg, self.mchClient.ApiKey(), nil)
		if !security.SecureCompareString(haveSignature, wantSignature) {
			err = fmt.Errorf("sign mismatch,\nhave: %s,\nwant: %s", haveSignature, wantSignature)
			holmes.Error("%v", err)
			return
		}
	}

	openId, ok := msg["openid"]
	if !ok {
		holmes.Error("msg openid not found.")
		return
	}
	totalFee, ok := msg["total_fee"]
	if !ok {
		holmes.Error("msg total_fee not found.")
		return
	}
	money, err := strconv.Atoi(totalFee)
	if err != nil {
		holmes.Error("msg total_fee[%s] strconv error: %v", totalFee, err)
		return
	}

	// --- new logic ---
	course := &models.Course{
		CourseType: READING_COURSE_TYPE_GD,
	}
	has, err := models.GetCourseMaxNum(course)
	if err != nil {
		holmes.Error("get course max num error: %v", err)
		return
	}
	if !has {
		holmes.Error("get course max num has none")
		return
	}
	user := &models.User{
		OpenId: openId,
	}
	has, err = models.GetUserFromOpenid(user)
	if err != nil {
		holmes.Error("[pay notify] get user error: %v", err)
		return
	}
	if !has {
		holmes.Error("openid[%s] not found, create it", openId)
		user.AppId = self.l.cfg.ReadingOauth.ReadingWxAppId
		err = models.CreateUser(user)
		if err != nil {
			holmes.Error("create user error: %v", err)
			return
		}
	}
	userCourse := &models.UserCourse{
		UserId:     user.ID,
		CourseId:   course.ID,
		CourseType: READING_COURSE_TYPE_GD,
		Money:      int64(money),
		Status:     READING_COURSE_STATUS_PAIED,
		PayTime:    time.Now().Unix(),
	}
	err = models.CreateUserCourse(userCourse)
	if err != nil {
		holmes.Error("create user course error: %v", err)
	}
	// --- end ---

	//readingUser := &models.ReadingPay{
	//	OpenId: openId,
	//}
	//has, err := models.GetReadingPay(readingUser)
	//if err != nil {
	//	holmes.Error("get reading pay error: %v", err)
	//	return
	//}
	//if !has {
	//	holmes.Error("openid[%s] not found, create it", openId)
	//	readingUser = &models.ReadingPay{
	//		AppId:  self.l.cfg.ReadingOauth.ReadingWxAppId,
	//		OpenId: openId,
	//		Money:  int64(money),
	//		Status: READING_COURSE_STATUS_PAIED,
	//	}
	//	err = models.CreateReadingPay(readingUser)
	//	if err != nil {
	//		holmes.Error("create reading pay error: %v", err)
	//	}
	//	return
	//}
	//readingUser.Money = int64(money)
	//readingUser.Status = READING_COURSE_STATUS_PAIED
	//readingUser.Number = self.l.cfg.NowCourseNumber
	//err = models.UpdateReadingPayStatus(readingUser)
	//if err != nil {
	//	holmes.Error("update reading pay status error: %v", err)
	//}

	// send sms notify
	//err = self.smsExt.SMSNotify(user.Phone, user.RealName)
	//if err != nil {
	//	holmes.Error("sms notify error: %v", err)
	//}
	params := fmt.Sprintf("#RealName#=%s&#StartTime#=%s", user.RealName, time.Unix(course.StartTime, 0).Format("2006-01-02"))
	err = self.smsExt.SMSNotifyNormal(user.Phone, self.l.cfg.SMSNotify.TemplateId, params)
	if err != nil {
		holmes.Error("sms notify error: %v", err)
	}
}

func (self *ReadingHandler) readingSuccess(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		io.WriteString(w, err.Error())
		holmes.Error("url parse query error: %v", err)
		return
	}
	openid := queryValues.Get("openid")

	// --- new logic ---
	userCourseList, err := models.GetUserCourse(openid)
	if err != nil {
		holmes.Error("get reading pay error: %v", err)
		io.WriteString(w, MSG_ERROR_SYSTEM)
		return
	}
	for _, v := range userCourseList {
		if v.Course.CourseType == READING_COURSE_TYPE_GD {
			readingUserInfo := &ReadingEnrollUserInfo{
				NickName:   v.User.Name,
				AvatarUrl:  v.User.AvatarUrl,
				OpenId:     openid,
				CourseType: v.Course.CourseType,
				CourseNum:  v.Course.CourseNum,
				StartTime:  time.Unix(v.Course.StartTime, 0).Format("2006.01.02"),
				EndTime:    time.Unix(v.Course.EndTime, 0).Format("2006.01.02"),
			}
			if v.Course.StartTime <= time.Now().Unix() {
				readingUserInfo.IfCourseStart = 1
			}
			renderView(w, "./views/reading_sign_success.html", readingUserInfo)
			return
		}
	}
	//io.WriteString(w, MSG_ERROR_COURSE_NOT_FOUND)
	course := &models.Course{
		CourseType: READING_COURSE_TYPE_GD,
	}
	has, err := models.GetCourseMaxNum(course)
	if err != nil {
		holmes.Error("get course max num error: %v", err)
		io.WriteString(w, MSG_ERROR_SYSTEM)
		return
	}
	if !has {
		holmes.Error("get course max num has none")
		io.WriteString(w, MSG_ERROR_COURSE_NOT_FOUND)
		return
	}
	user := &models.User{
		OpenId: openid,
	}
	has, err = models.GetUserFromOpenid(user)
	if err != nil {
		holmes.Error("get user error: %v", err)
		io.WriteString(w, MSG_ERROR_SYSTEM)
		return
	}
	if !has {
		holmes.Error("get user has none")
		io.WriteString(w, MSG_ERROR_USER_NOT_FOUND)
		return
	}
	readingUserInfo := &ReadingEnrollUserInfo{
		NickName:   user.Name,
		AvatarUrl:  user.AvatarUrl,
		OpenId:     openid,
		CourseType: course.CourseType,
		CourseNum:  course.CourseNum,
		StartTime:  time.Unix(course.StartTime, 0).Format("2006.01.02"),
		EndTime:    time.Unix(course.EndTime, 0).Format("2006.01.02"),
	}
	if course.StartTime <= time.Now().Unix() {
		readingUserInfo.IfCourseStart = 1
	}
	renderView(w, "./views/reading_sign_success.html", readingUserInfo)
	// --- end ---

	//readingUser := &models.ReadingPay{
	//	OpenId: openid,
	//}
	//has, err := models.GetReadingPay(readingUser)
	//if err != nil {
	//	io.WriteString(w, "未找到你哦,请刷新重新登录")
	//	holmes.Error("get reading pay error: %v", err)
	//	return
	//}
	//if !has {
	//	io.WriteString(w, "未找到你哦,请刷新重新登录")
	//	return
	//}
	//
	//readingUserInfo := &ReadingEnrollUserInfo{
	//	NickName:  readingUser.Name,
	//	AvatarUrl: readingUser.AvatarUrl,
	//	OpenId:    openid,
	//}
	//renderView(w, "./views/reading_sign_success.html", readingUserInfo)
}

func (self *ReadingHandler) getOauthUserInfo(w http.ResponseWriter, r *http.Request) (bool, *mpoauth2.UserInfo, error) {
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		io.WriteString(w, err.Error())
		holmes.Error("url parse query error: %v", err)
		return false, nil, err
	}

	code := queryValues.Get("code")
	if code == "" {
		state := string(rand.NewHex())
		redirectUrl := fmt.Sprintf("http://%s%s", r.Host, r.URL.String())
		AuthCodeURL := mpoauth2.AuthCodeURL(self.l.cfg.ReadingOauth.ReadingWxAppId,
			redirectUrl,
			self.l.cfg.ReadingOauth.ReadingOauth2ScopeUser, state)
		http.Redirect(w, r, AuthCodeURL, http.StatusFound)
		return true, nil, nil
	}

	token, err := self.oauth2Client.ExchangeToken(code)
	if err != nil {
		//holmes.Error("exchange token error: %v", err)
		http.Redirect(w, r, fmt.Sprintf("http://%s%s", r.Host, r.URL.Path), http.StatusFound)
		return false, nil, err
	}
	holmes.Debug("token: %+v", token)

	userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		io.WriteString(w, err.Error())
		holmes.Error("get user info error: %v", err)
		return false, nil, err
	}
	holmes.Debug("user info: %+v", userinfo)

	return false, userinfo, nil
}

func (self *ReadingHandler) readingUnifiedOrder(payId, payMoney int64, openId, userIp, notifyUrl string) (*mchpay.UnifiedOrderResponse, error) {
	uor := &mchpay.UnifiedOrderRequest{
		DeviceInfo:     "WEB",
		Body:           "共读计划",
		Detail:         "「小鹿微课」共读计划",
		Attach:         "attach",
		OutTradeNo:     fmt.Sprintf("%s-%d", time.Now().Format("2006-01-02_15-04-05"), payId),
		TotalFee:       payMoney,
		SpbillCreateIP: userIp,
		NotifyURL:      notifyUrl,
		TradeType:      "JSAPI",
		OpenId:         openId,
	}

	rsp, err := mchpay.UnifiedOrder2(self.mchClient, uor)
	if err != nil {
		holmes.Error("mch pay unified order error: %v", err)
		return nil, err
	}

	return rsp, nil
}

func parseRequest(r *http.Request) (*HandlerRequest, error) {
	req := &HandlerRequest{}
	req.Path = r.URL.Path[len(ReadingPrefix)+1:]

	result, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return req, errors.New("parse request read error")
	}
	r.Body.Close()

	req.Method = r.Method
	req.Val = result

	return req, nil
}

func writeRsp(w http.ResponseWriter, rsp *proto.Response) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")

	if rsp != nil {
		WriteJSON(w, http.StatusOK, rsp)
	}
}

func renderView(w http.ResponseWriter, tpl string, data interface{}) {
	t := template.New(filepath.Base(tpl))
	t = t.Funcs(template.FuncMap{"unescaped": unescaped})
	t, err := t.ParseFiles(tpl)
	if err != nil {
		holmes.Error("parse file error: %v", err)
		return
	}
	err = t.ExecuteTemplate(w, filepath.Base(tpl), data)
	if err != nil {
		holmes.Error("execute tmp error: %v", err)
		return
	}
}

func unescaped(x string) interface{} { return template.HTML(x) }

func GetIPFromRequest(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		holmes.Error("userip: %q is not IP:port", r.RemoteAddr)
		return ""
	}
	userIP := net.ParseIP(ip)
	if userIP == nil {
		holmes.Error("userip: %q is not IP:port", r.RemoteAddr)
		return ""
	}
	return userIP.String()
}
