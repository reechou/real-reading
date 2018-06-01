package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/util"
	"github.com/chanxuehong/util/security"
	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/models"
	"github.com/reechou/real-reading/proto"
	mchcore "gopkg.in/chanxuehong/wechat.v2/mch/core"
	mchpay "gopkg.in/chanxuehong/wechat.v2/mch/pay"
)

const (
	REGISTER_URI_PREFIX = "register"
)

func (self *ReadingHandler) registerHandle(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	//holmes.Debug("register handle uri: %s", rr.Path)
	subPath := rr.Path[len(REGISTER_URI_PREFIX)+1:]
	rr.Params = strings.Split(subPath, "/")
	if len(rr.Params) == 0 {
		holmes.Error("path error: %s", rr.Path)
		return
	}
	switch rr.Params[0] {
	case READING_URI_SIGN_UP:
		self.registerSignup(rr, w, r)
	case READING_URI_ENROLL:
		self.registerEnroll(rr, w, r)
	case READING_URI_GO_ENROLL:
		self.registerGoEnroll(rr, w, r)
	case READING_URI_PAY:
		self.registerPay(rr, w, r)
	case READING_URI_PAY_NOTIFY:
		self.registerPayNotify(rr, w, r)
	case READING_URI_SUCCESS:
		self.registerSuccess(rr, w, r)
	case READING_URI_PROTO:
		self.registerProto(rr, w, r)
	}
}

func (self *ReadingHandler) registerSignup(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	if len(rr.Params) < 2 {
		holmes.Error("params error: %v", rr.Params)
		return
	}

	registerInfo := new(RegisterInfo)
	var err error
	registerInfo.Course.CourseType, err = strconv.ParseInt(rr.Params[1], 10, 0)
	if err != nil {
		holmes.Error("params[1][%s] strconv error: %v", rr.Params[1], err)
		return
	}
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		holmes.Error("url parse query error: %v", err)
		return
	}
	src := queryValues.Get("src")
	if src != "" {
		registerInfo.Source, err = strconv.Atoi(src)
		if err != nil {
			holmes.Error("strconv src[%s] error: %v", src, err)
		}
	}

	has, err := models.GetCourseMaxNum(&registerInfo.Course)
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

	registerInfo.Course.Money /= 100
	registerInfo.IndexStartTime = time.Unix(registerInfo.Course.StartTime, 0).Format("2006年01月02日")
	registerInfo.StartTime = time.Unix(registerInfo.Course.StartTime, 0).Format("2006.01")
	registerInfo.EndTime = time.Unix(registerInfo.Course.EndTime, 0).Format("2006.01")

	renderView(w, "./views/register/sign.html", registerInfo)
}

func (self *ReadingHandler) registerEnroll(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	userinfo, ifRedirect := self.checkUserBase(w, r)
	if ifRedirect {
		return
	}

	if len(rr.Params) < 2 {
		holmes.Error("params error: %v", rr.Params)
		return
	}

	registerInfo := new(RegisterInfo)
	registerInfo.OpenId = userinfo.OpenId
	var err error
	registerInfo.Course.CourseType, err = strconv.ParseInt(rr.Params[1], 10, 0)
	if err != nil {
		holmes.Error("params[1][%s] strconv error: %v", rr.Params[1], err)
		return
	}

	user := &models.User{
		OpenId: userinfo.OpenId,
	}
	has, err := models.GetUserFromOpenid(user)
	if err != nil {
		holmes.Error("get user from openid error: %v", err)
		io.WriteString(w, MSG_ERROR_SYSTEM)
		return
	}
	if !has {
		user.AppId = self.l.cfg.ReadingOauth.ReadingWxAppId
		user.Name = userinfo.Name
		user.AvatarUrl = userinfo.AvatarUrl
		user.Source = int64(registerInfo.Source)
		err = models.CreateUser(user)
		if err != nil {
			holmes.Error("create user error: %v", err)
			return
		}
		io.WriteString(w, MSG_ERROR_USER_COURSE_NOT_JOIN)
		return
	}

	registerInfo.NickName = user.Name
	registerInfo.AvatarUrl = user.AvatarUrl
	registerInfo.Source = int(user.Source)

	courseList, err := models.GetUserCourseList(user.ID)
	if err != nil {
		holmes.Error("get user course list error: %v", err)
		io.WriteString(w, MSG_ERROR_SYSTEM)
		return
	}
	var ifHasThisCourse bool
	for _, v := range courseList {
		if v.Course.CourseType == registerInfo.Course.CourseType {
			registerInfo.Course = v.Course
			registerInfo.StartTime = time.Unix(v.Course.StartTime, 0).Format("2006.01.02")
			registerInfo.EndTime = time.Unix(v.Course.EndTime, 0).Format("2006.01.02")
			if v.Course.StartTime <= time.Now().Unix() {
				registerInfo.IfCourseStart = 1
			}
			//renderView(w, "./views/register/sign_success.html", registerInfo)
			//return
			ifHasThisCourse = true
		}
	}
	if !ifHasThisCourse {
		holmes.Error("user[%d] not has this course[%d], but in this course's enroll.", user.ID, registerInfo.Course.CourseType)
		io.WriteString(w, MSG_ERROR_USER_COURSE_NOT_JOIN)
		return
	}

	//has, err = models.GetCourseMaxNum(&registerInfo.Course)
	//if err != nil {
	//	holmes.Error("get course max num error: %v", err)
	//	io.WriteString(w, MSG_ERROR_SYSTEM)
	//	return
	//}
	//if !has {
	//	holmes.Error("get course max num has none")
	//	io.WriteString(w, MSG_ERROR_COURSE_NOT_FOUND)
	//	return
	//}

	if user.RealName != "" {
		registerInfo.EnrollName = user.RealName
	}
	if user.Phone != "" {
		registerInfo.EnrollMobile = user.Phone
	}
	if user.Wechat != "" {
		registerInfo.EnrollWechat = user.Wechat
	}
	//registerInfo.StartTime = time.Unix(registerInfo.Course.StartTime, 0).Format("2006-01-02")
	//registerInfo.EndTime = time.Unix(registerInfo.Course.EndTime, 0).Format("2006-01-02")
	renderView(w, "./views/register/enroll.html", registerInfo)
}

func (self *ReadingHandler) registerGoEnroll(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
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

	course := &models.Course{
		ID: req.CourseId,
	}
	has, err = models.GetCourse(course)
	if err != nil {
		holmes.Error("get course error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	if !has {
		holmes.Error("cannot found this course id[%d]", req.CourseId)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	params := fmt.Sprintf("#RealName#=%s&#PlanName#=%s&#StartTime#=%s", user.RealName, course.Name, time.Unix(course.StartTime, 0).Format("2006-01-02"))
	err = self.smsExt.SMSNotifyNormal(user.Phone, self.l.cfg.SMSNotify.RegisterTid, params)
	if err != nil {
		holmes.Error("sms notify error: %v", err)
	}
}

func (self *ReadingHandler) registerPay(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	userinfo, ifRedirect := self.checkUser(w, r, true)
	if ifRedirect {
		return
	}

	if len(rr.Params) < 2 {
		holmes.Error("params error: %v", rr.Params)
		return
	}

	registerInfo := new(RegisterInfo)
	registerInfo.OpenId = userinfo.OpenId
	registerInfo.Source = userinfo.Source
	registerInfo.NickName = userinfo.Name
	registerInfo.AvatarUrl = userinfo.AvatarUrl
	var err error
	registerInfo.Course.CourseType, err = strconv.ParseInt(rr.Params[1], 10, 0)
	if err != nil {
		holmes.Error("params[1][%s] strconv error: %v", rr.Params[1], err)
		return
	}

	user := &models.User{
		OpenId: userinfo.OpenId,
	}
	has, err := models.GetUserFromOpenid(user)
	if err != nil {
		holmes.Error("get user from openid error: %v", err)
		io.WriteString(w, MSG_ERROR_SYSTEM)
		return
	}
	if !has {
		user.AppId = self.l.cfg.ReadingOauth.ReadingWxAppId
		user.Name = userinfo.Name
		user.AvatarUrl = userinfo.AvatarUrl
		user.Source = int64(registerInfo.Source)
		err = models.CreateUser(user)
		if err != nil {
			holmes.Error("create user error: %v", err)
			return
		}
	}

	// check if user has this course
	courseList, err := models.GetUserCourseList(user.ID)
	if err != nil {
		holmes.Error("get user course list error: %v", err)
		io.WriteString(w, MSG_ERROR_SYSTEM)
		return
	}
	for _, v := range courseList {
		if v.Course.CourseType == registerInfo.Course.CourseType {
			registerInfo.Course = v.Course
			registerInfo.StartTime = time.Unix(v.Course.StartTime, 0).Format("2006.01.02")
			registerInfo.EndTime = time.Unix(v.Course.EndTime, 0).Format("2006.01.02")
			if v.Course.StartTime <= time.Now().Unix() {
				registerInfo.IfCourseStart = 1
			}
			renderView(w, "./views/register/sign_success.html", registerInfo)
			return
		}
	}

	has, err = models.GetCourseMaxNum(&registerInfo.Course)
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

	payMoney := registerInfo.Course.Money
	//if userinfo.OpenId == "oaKrZwsAF6pRX6z3Qn_EhIZ3DG90" || userinfo.OpenId == "oaKrZwotcenPmZyLKtMyoHZSTlaQ" {
	//	payMoney = 1
	//}
	var useCoupon int64
	if userinfo.CouponCDKey != "" {
		registerInfo.Coupon.CdKey = userinfo.CouponCDKey
		if has, err := models.GetCoupon(&registerInfo.Coupon); err != nil {
			holmes.Error("get coupon error: %v", err)
		} else {
			if has {
				if registerInfo.Coupon.IfUse == 0 {
					if time.Now().Unix() > registerInfo.Coupon.LockTime ||
						(registerInfo.Coupon.CourseType == registerInfo.Coupon.CourseType && user.ID == registerInfo.Coupon.UserId) {
						registerInfo.Coupon.CourseType = registerInfo.Course.CourseType
						registerInfo.Coupon.UserId = user.ID
						if err = models.UpdateCouponLock(&registerInfo.Coupon); err != nil {
							holmes.Error("update coupon lock error: %v", err)
						} else {
							useCoupon = registerInfo.Coupon.ID
							if payMoney > registerInfo.Coupon.Amount {
								payMoney = payMoney - registerInfo.Coupon.Amount
							} else {
								payMoney = 1
							}
							holmes.Info("user[%d] has locked the coupon[%s]: %s and new pay: %d",
								user.ID, registerInfo.Coupon.Desc, registerInfo.Coupon.CdKey, payMoney)
						}
					}
				}
			}
		}
	}

	registerInfo.RealPay = strconv.FormatFloat(float64(payMoney)/100, 'f', -1, 64)

	outTradeNo := fmt.Sprintf("%s-%d-%d", time.Now().Format("2006-01-02_15-04-05"), user.ID, registerInfo.Course.ID)
	unifiedRsp, err := self.registerUnifiedOrder(
		outTradeNo,
		payMoney,
		registerInfo.Course.Name,
		fmt.Sprintf("%d_%d", registerInfo.Course.CourseType, useCoupon),
		user.OpenId,
		GetIPFromRequest(r),
		fmt.Sprintf("%s%s/%s/%s", r.Host, ReadingPrefix, REGISTER_URI_PREFIX, READING_URI_PAY_NOTIFY),
	)
	if err != nil {
		holmes.Error("reading unified order error: %v", err)
		if strings.Contains(err.Error(), "ORDERPAID") {
			// 订单已支付, 但未更新状态
			userCourse := &models.UserCourse{
				UserId:     user.ID,
				CourseId:   registerInfo.Course.ID,
				CourseType: registerInfo.Course.CourseType,
				Money:      payMoney,
				Status:     READING_COURSE_STATUS_PAIED,
				PayTime:    time.Now().Unix(),
				Source:     int64(registerInfo.Source),
			}
			err = models.CreateUserCourse(userCourse)
			if err != nil {
				holmes.Error("create user course error: %v", err)
			}
			registerInfo.StartTime = time.Unix(registerInfo.Course.StartTime, 0).Format("2006.01.02")
			registerInfo.EndTime = time.Unix(registerInfo.Course.EndTime, 0).Format("2006.01.02")
			renderView(w, "./views/register/sign_success.html", registerInfo)
			return
		}
		return
	}

	registerInfo.Course.Money /= 100

	registerInfo.WxJsApiParams = WxJsApiParams{
		AppId:     self.l.cfg.ReadingOauth.ReadingWxAppId,
		TimeStamp: fmt.Sprintf("%d", time.Now().Unix()),
		NonceStr:  string(rand.NewHex()),
		Package:   fmt.Sprintf("prepay_id=%s", unifiedRsp.PrepayId),
		SignType:  "MD5",
	}
	registerInfo.WxJsApiParams.PaySign = mchcore.JsapiSign(
		registerInfo.WxJsApiParams.AppId,
		registerInfo.WxJsApiParams.TimeStamp,
		registerInfo.WxJsApiParams.NonceStr,
		registerInfo.WxJsApiParams.Package,
		registerInfo.WxJsApiParams.SignType,
		self.mchClient.ApiKey(),
	)
	renderView(w, "./views/register/pay.html", registerInfo)
}

func (self *ReadingHandler) registerPayNotify(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	defer func() {
		io.WriteString(w, "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>")
	}()

	msg, err := util.DecodeXMLToMap(bytes.NewReader(rr.Val))
	if err != nil {
		holmes.Error("decode xml error: %v", err)
		return
	}
	holmes.Debug("pay notify msg: %v", msg)

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
	attach, ok := msg["attach"]
	if !ok {
		holmes.Error("msg attach not found.")
		return
	}
	attachs := strings.Split(attach, "_")
	if len(attach) != 2 {
		holmes.Error("attach error: %s", attach)
		return
	}
	courseType, err := strconv.Atoi(attachs[0])
	if err != nil {
		holmes.Error("msg attach-course_type[%s] strconv error: %v", attachs[0], err)
		return
	}
	couponId, err := strconv.ParseInt(attachs[1], 10, 0)
	if err != nil {
		holmes.Error("msg attach-coupon_id[%s] strconv error: %v", attachs[1], err)
		return
	}
	if couponId != 0 {
		if err = models.UpdateCoupon(&models.Coupon{ID: couponId}); err != nil {
			holmes.Error("update coupon error: %v", err)
		}
	}
	outTradeNo, ok := msg["out_trade_no"]
	if !ok {
		holmes.Error("msg out_trade_no not found.")
		return
	}
	transactionId, ok := msg["transaction_id"]
	if !ok {
		holmes.Error("msg transaction_id not found.")
		return
	}

	course := &models.Course{
		CourseType: int64(courseType),
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
		UserId:        user.ID,
		CourseId:      course.ID,
		CourseType:    course.CourseType,
		Money:         int64(money),
		Status:        READING_COURSE_STATUS_PAIED,
		PayTime:       time.Now().Unix(),
		Source:        user.Source,
		OutTradeNo:    outTradeNo,
		TransactionId: transactionId,
	}
	err = models.CreateUserCourse(userCourse)
	if err != nil {
		holmes.Error("create user course error: %v", err)
	}

	//params := fmt.Sprintf("#RealName#=%s&#PlanName#=%s&#StartTime#=%s", user.RealName, course.Name, time.Unix(course.StartTime, 0).Format("2006-01-02"))
	//err = self.smsExt.SMSNotifyNormal(user.Phone, self.l.cfg.SMSNotify.RegisterTid, params)
	//if err != nil {
	//	holmes.Error("sms notify error: %v", err)
	//}
}

func (self *ReadingHandler) registerSuccess(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	userinfo, ifRedirect := self.checkUserBase(w, r)
	if ifRedirect {
		return
	}

	if len(rr.Params) < 2 {
		holmes.Error("params error: %v", rr.Params)
		return
	}

	registerInfo := new(RegisterInfo)
	registerInfo.OpenId = userinfo.OpenId
	var err error
	registerInfo.Course.CourseType, err = strconv.ParseInt(rr.Params[1], 10, 0)
	if err != nil {
		holmes.Error("params[1][%s] strconv error: %v", rr.Params[1], err)
		return
	}

	userCourseList, err := models.GetUserCourse(registerInfo.OpenId)
	if err != nil {
		holmes.Error("get reading pay error: %v", err)
		io.WriteString(w, MSG_ERROR_SYSTEM)
		return
	}
	for _, v := range userCourseList {
		if v.Course.CourseType == registerInfo.Course.CourseType {
			registerInfo.NickName = v.User.Name
			registerInfo.AvatarUrl = v.User.AvatarUrl
			registerInfo.Course = v.Course
			registerInfo.StartTime = time.Unix(v.Course.StartTime, 0).Format("2006.01.02")
			registerInfo.EndTime = time.Unix(v.Course.EndTime, 0).Format("2006.01.02")
			if v.Course.StartTime <= time.Now().Unix() {
				registerInfo.IfCourseStart = 1
			}
			holmes.Debug("success register info: %v", registerInfo)
			renderView(w, "./views/register/sign_success.html", registerInfo)
			return
		}
	}
	has, err := models.GetCourseMaxNum(&registerInfo.Course)
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
		OpenId: registerInfo.OpenId,
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
	registerInfo.NickName = user.Name
	registerInfo.AvatarUrl = user.AvatarUrl
	registerInfo.StartTime = time.Unix(registerInfo.Course.StartTime, 0).Format("2006.01.02")
	registerInfo.EndTime = time.Unix(registerInfo.Course.EndTime, 0).Format("2006.01.02")
	if registerInfo.Course.StartTime <= time.Now().Unix() {
		registerInfo.IfCourseStart = 1
	}
	renderView(w, "./views/register/sign_success.html", registerInfo)
}

func (self *ReadingHandler) registerProto(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	renderView(w, "./views/register/proto.html", nil)
}

func (self *ReadingHandler) registerUnifiedOrder(outTradeNo string, payMoney int64, title, attach, openId, userIp, notifyUrl string) (*mchpay.UnifiedOrderResponse, error) {
	uor := &mchpay.UnifiedOrderRequest{
		DeviceInfo:     "WEB",
		Body:           title,
		Detail:         fmt.Sprintf("「小鹿微课」%s", title),
		Attach:         attach,
		OutTradeNo:     outTradeNo,
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
