package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/models"
	"github.com/reechou/real-reading/proto"
)

const (
	COUPON_URI_PREFIX      = "coupon"
	COUPON_URI_CREATE      = "create"
	COUPON_URI_EXCHANGE    = "exchange"
	COUPON_URI_GO_EXCHANGE = "go-exchange"
	COUPON_URI_LIST        = "list"
)

func (self *ReadingHandler) couponHandle(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	holmes.Debug("coupon handle uri: %s", rr.Path)
	subPath := rr.Path[len(COUPON_URI_PREFIX)+1:]
	rr.Params = strings.Split(subPath, "/")
	if len(rr.Params) == 0 {
		holmes.Error("path error: %s", rr.Path)
		return
	}
	switch rr.Params[0] {
	case COUPON_URI_CREATE:
		self.createCoupon(rr, w, r)
	case COUPON_URI_EXCHANGE:
		self.exchangeCoupon(rr, w, r)
	case COUPON_URI_GO_EXCHANGE:
		self.goexchangeCoupon(rr, w, r)
	case COUPON_URI_LIST:
		self.listCoupon(rr, w, r)
	}
}

func (self *ReadingHandler) createCoupon(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &proto.CreateCouponReq{}
	var err error
	if err = json.Unmarshal(rr.Val, &req); err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	now := time.Now().Unix()
	coupons := make([]models.Coupon, req.Num)
	for i := 0; i < req.Num; i++ {
		coupons[i].Name = req.Name
		coupons[i].Desc = req.Desc
		coupons[i].Amount = req.Amount
		coupons[i].CdKey = uniuri.New()
		coupons[i].CreatedAt = now
		coupons[i].UpdatedAt = now
	}
	if err = models.CreateCouponList(coupons); err != nil {
		holmes.Error("create coupon list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = coupons
}

func (self *ReadingHandler) exchangeCoupon(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
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

	renderView(w, "./views/register/coupon.html", registerInfo)
}

func (self *ReadingHandler) goexchangeCoupon(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.Coupon{}
	var err error
	if err = json.Unmarshal(rr.Val, &req); err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	has, err := models.GetCoupon(req)
	if err != nil {
		holmes.Error("get coupon error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	if !has {
		holmes.Error("coupon[%s] not found", req.CdKey)
		rsp.Code = proto.RESPONSE_ERR_COUPON_NOT_FOUND
		rsp.Msg = "未找到该优惠券码"
		return
	}
	if req.IfUse != 0 {
		rsp.Code = proto.RESPONSE_COUPON_HAS_USED
		rsp.Msg = "该优惠券码已使用"
		return
	}
	if time.Now().Unix() < req.LockTime {
		rsp.Code = proto.RESPONSE_COUPON_HAS_LOCKED
		lost := req.LockTime - time.Now().Unix()
		rsp.Msg = fmt.Sprintf("该优惠券码已锁定，请%d分%d秒后再试", lost/60, lost%60)
		return
	}
}

func (self *ReadingHandler) listCoupon(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &proto.GetCouponListReq{}
	var err error
	if err = json.Unmarshal(rr.Val, &req); err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	list, err := models.GetCouponList(req.Offset, req.Num)
	if err != nil {
		holmes.Error("get coupon list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = list
}
