package controller

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"encoding/json"

	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/session"
	mchpay "github.com/chanxuehong/wechat.v2/mch/pay"
	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/chanxuehong/wechat.v2/oauth2"
	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/proto"
)

const (
	ReadingPrefix = "/reading"
)

const (
	READING_URI_ENROLL    = "enroll"
	READING_URI_GO_ENROLL = "goenroll"
	READING_URI_PAY       = "pay"
	READING_URI_SUCCESS   = "success"
)

type ShareTpl struct {
	Title string
	Img   string
}

type HandlerRequest struct {
	Method string
	Path   string
	Val    []byte
}

type ReadingHandler struct {
	l *Logic

	lefitSessionStorage *session.Storage
	lefitOauth2Endpoint oauth2.Endpoint
	oauth2Client        *oauth2.Client
}

func NewReadingHandler(l *Logic) *ReadingHandler {
	lh := &ReadingHandler{l: l}

	lh.lefitSessionStorage = session.New(20*60, 60*60)
	lh.lefitOauth2Endpoint = mpoauth2.NewEndpoint(
		lh.l.cfg.ReadingOauth.ReadingWxAppId,
		lh.l.cfg.ReadingOauth.ReadingWxAppSecret)
	lh.oauth2Client = &oauth2.Client{
		Endpoint: lh.lefitOauth2Endpoint,
	}

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

	switch rr.Path {
	case READING_URI_ENROLL:
		self.readingEnroll(rr, w, r)
	case READING_URI_GO_ENROLL:
		self.readingGoEnroll(rr, w, r)
	case READING_URI_PAY:
		self.readingPay(rr, w, r)
	case READING_URI_SUCCESS:
		self.readingSuccess(rr, w, r)
	default:
		http.ServeFile(w, r, self.l.cfg.ReadingOauth.MpVerifyDir+rr.Path)
	}
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
		io.WriteString(w, err.Error())
		holmes.Error("get user info error: %v", err)
		return
	}
	holmes.Debug("user info: %+v", userinfo)

	readingUserInfo := &ReadingEnrollUserInfo{
		NickName:  userinfo.Nickname,
		AvatarUrl: userinfo.HeadImageURL,
		OpenId:    token.OpenId,
	}
	renderView(w, "./views/reading_enroll.html", readingUserInfo)
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
	holmes.Debug("reading go enroll: %+v %s", req, r.URL.String())
}

func (self *ReadingHandler) readingPay(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		io.WriteString(w, err.Error())
		holmes.Error("url parse query error: %v", err)
		return
	}
	openid := queryValues.Get("openid")
	
	readingUserInfo := &ReadingEnrollUserInfo{
		NickName:  "xxxxxx",
		AvatarUrl: "http://wx.qlogo.cn/mmopen/ibmyOaFEgYk09HCYrBXA7PHZSuFjHINfuNxBlIOyvPibrU0hD87gTrGI2YuBTtGibHrxdTyzFAMFvWIPO5ekuhibzQ/0",
		OpenId:    openid,
	}
	renderView(w, "./views/reading_pay.html", readingUserInfo)
}

func (self *ReadingHandler) readingSuccess(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		io.WriteString(w, err.Error())
		holmes.Error("url parse query error: %v", err)
		return
	}
	openid := queryValues.Get("openid")
	
	readingUserInfo := &ReadingEnrollUserInfo{
		NickName:  "xxxxxx",
		AvatarUrl: "http://wx.qlogo.cn/mmopen/ibmyOaFEgYk09HCYrBXA7PHZSuFjHINfuNxBlIOyvPibrU0hD87gTrGI2YuBTtGibHrxdTyzFAMFvWIPO5ekuhibzQ/0",
		OpenId:    openid,
	}
	renderView(w, "./views/reading_sign_success.html", readingUserInfo)
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

func (self *ReadingHandler) readingUnifiedOrderRequest(payMoney int64, openId, userIp, notifyUrl string) *mchpay.UnifiedOrderRequest {
	uor := &mchpay.UnifiedOrderRequest{
		DeviceInfo:     "WEB",
		Body:           "reading",
		Detail:         "reading detail",
		Attach:         "attach",
		TotalFee:       payMoney,
		SpbillCreateIP: userIp,
		NotifyURL:      notifyUrl,
		TradeType:      "JSAPI",
		OpenId:         openId,
	}

	return uor
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
	w.Header().Set("Content-Type", "application/json")

	if rsp != nil {
		WriteJSON(w, http.StatusOK, rsp)
	}
}

func renderView(w http.ResponseWriter, tpl string, data interface{}) {
	t, err := template.ParseFiles(tpl)
	if err != nil {
		holmes.Error("parse file error: %v", err)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		holmes.Error("execute tmp error: %v", err)
		return
	}
}
