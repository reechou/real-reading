package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/models"
	"github.com/reechou/real-reading/proto"
	mchpay "gopkg.in/chanxuehong/wechat.v2/mch/pay"
)

const (
	DATA_STATISTICS_URI_PREFIX = "statistics"
)

const (
	REFUND_WAY_AUTO             = 1
	REFUND_WAY_FROM_TRANSACTION = 2
)

const (
	DATA_STATISTICS_CREATE_COURSE_TYPE                  = "createcoursetype"
	DATA_STATISTICS_GET_COURSE_TYPE_LIST                = "getcoursetypelist"
	DATA_STATISTICS_CREATE_COURSE_CHANNEL               = "createcoursechannel"
	DATA_STATISTICS_DELETE_COURSE_CHANNEL               = "deletecoursechannel"
	DATA_STATISTICS_GET_COURSE_CHANNEL_LIST             = "getcoursechannellist"
	DATA_STATISTICS_SET_USER_COURSE_REFUND              = "setusercourserefund"
	DATA_STATISTICS_USER_COURSE_MANUAL_REFUND           = "usercoursemanualrefund"
	DATA_STATISTICS_USER_COURSE_REFUND_FROM_TRANSACTION = "usercourserefundfromtransaction"
	DATA_STATISTICS_GET_COURSE_STATISTICS               = "getcoursedatastatistics"
)

func (self *ReadingHandler) dataStatisticsHandle(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	subPath := rr.Path[len(DATA_STATISTICS_URI_PREFIX)+1:]
	rr.Params = strings.Split(subPath, "/")
	if len(rr.Params) == 0 {
		holmes.Error("path error: %s", rr.Path)
		return
	}
	switch rr.Params[0] {
	case DATA_STATISTICS_CREATE_COURSE_TYPE:
		self.createCourseType(rr, w, r)
	case DATA_STATISTICS_GET_COURSE_TYPE_LIST:
		self.getCourseTypeList(rr, w, r)
	case DATA_STATISTICS_CREATE_COURSE_CHANNEL:
		self.createCourseChannel(rr, w, r)
	case DATA_STATISTICS_DELETE_COURSE_CHANNEL:
		self.deleteCourseChannel(rr, w, r)
	case DATA_STATISTICS_GET_COURSE_CHANNEL_LIST:
		self.getCourseChannelList(rr, w, r)
	case DATA_STATISTICS_SET_USER_COURSE_REFUND:
		self.setUserCourseRefund(rr, w, r)
	case DATA_STATISTICS_USER_COURSE_MANUAL_REFUND:
		self.setUserCourseManualRefund(rr, w, r)
	case DATA_STATISTICS_USER_COURSE_REFUND_FROM_TRANSACTION:
		self.setUserCourseRefundFromTransaction(rr, w, r)
	case DATA_STATISTICS_GET_COURSE_STATISTICS:
		self.getCourseDataStatistics(rr, w, r)
	}
}

func (self *ReadingHandler) createCourseType(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.CourseType{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.CreateCourseType(req)
	if err != nil {
		holmes.Error("create course type error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) getCourseTypeList(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	list, err := models.GetCourseTypeList()
	if err != nil {
		holmes.Error("get course type list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = list
}

func (self *ReadingHandler) createCourseChannel(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.CourseChannel{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.CreateCourseChannel(req)
	if err != nil {
		holmes.Error("create course channel error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) deleteCourseChannel(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.CourseChannel{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.DelCourseChannel(req)
	if err != nil {
		holmes.Error("delete course channel error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) getCourseChannelList(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &proto.GetCourseChannelReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	list, err := models.GetCourseChannelList(req.CourseType)
	if err != nil {
		holmes.Error("get course channel list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = list
}

func (self *ReadingHandler) setUserCourseRefund(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.UserCourse{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	has, err := models.GetUserCourseFromId(req)
	if err != nil {
		holmes.Error("get user course from id error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	if !has {
		holmes.Error("user course[%d] not found", req.ID)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	if req.OutTradeNo == "" || req.TransactionId == "" {
		holmes.Debug("this user course[%d] cannot refund auto.", req.ID)
		rsp.Code = proto.RESPONSE_REFUND_NOT_AUTO
		rsp.Msg = "不能自动退款，请填写用户微信支付订单号"

		req.Status = READING_COURSE_STATUS_REFUND_NOT_AUTO
		err = models.UpdateUserCourseStatus(req)
		if err != nil {
			holmes.Error("update user course status of refund error: %v", err)
			rsp.Code = proto.RESPONSE_ERR
		}
		return
	}

	if req.RefundId != "" {
		rsp.Code = proto.RESPONSE_REFUND_HAS_EXECED
		rsp.Msg = "已退款，请勿重复退款"
		return
	}

	// auto refund
	refundReq := &mchpay.RefundRequest{
		TransactionId: req.TransactionId,
		OutTradeNo:    req.OutTradeNo,
		TotalFee:      req.Money,
		RefundFee:     req.Money,
		OutRefundNo:   fmt.Sprintf("%d", req.ID),
	}
	refundRsp, err := mchpay.Refund2(self.mchClient, refundReq)
	if err != nil {
		holmes.Error("mch pay refund error: %v", err)
		rsp.Code = proto.RESPONSE_ERR_EXT
		return
	}
	req.OutRefundNo = refundRsp.OutRefundNo
	req.RefundId = refundRsp.RefundId
	req.RefundFee = refundRsp.RefundFee
	req.Status = READING_COURSE_STATUS_REFUND
	req.RefundWay = REFUND_WAY_AUTO
	err = models.UpdateUserCourseRefundInfo(req)
	if err != nil {
		holmes.Error("update user course refund info error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) setUserCourseManualRefund(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.UserCourse{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	req.Status = READING_COURSE_STATUS_REFUND_MANUAL
	err = models.UpdateUserCourseStatus(req)
	if err != nil {
		holmes.Error("update user course status of refund error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) setUserCourseRefundFromTransaction(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &proto.ReadingRefundFromTransactionReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	userCourse := &models.UserCourse{
		ID: req.UserCourseId,
	}
	has, err := models.GetUserCourseFromId(userCourse)
	if err != nil {
		holmes.Error("get user course from id error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	if !has {
		holmes.Error("user course[%d] not found", userCourse.ID)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	user := &models.User{
		ID: userCourse.UserId,
	}
	has, err = models.GetUser(user)
	if err != nil {
		holmes.Error("get user error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	if !has {
		holmes.Error("cannot found this user[%d]", userCourse.UserId)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	orderQueryRsp, err := self.getWechatPayOrder(req.TransactionId)
	if err != nil {
		holmes.Error("get wechat pay order error: %v", err)
		rsp.Code = proto.RESPONSE_ERR_EXT
		return
	}
	if orderQueryRsp.TotalFee != userCourse.Money {
		holmes.Error("order query total fee[%d] != user course money[%d]", orderQueryRsp.TotalFee, userCourse.Money)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	if orderQueryRsp.OpenId != user.OpenId {
		holmes.Error("order query openid[%s] != user openid[%s]", orderQueryRsp.OpenId, user.OpenId)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	refundRsp, err := self.refundWechatPay(orderQueryRsp.TransactionId, orderQueryRsp.OutTradeNo, fmt.Sprintf("%d", userCourse.ID), userCourse.Money)
	if err != nil {
		holmes.Error("wechat pay refund error: %v", err)
		rsp.Code = proto.RESPONSE_ERR_EXT
		return
	}
	userCourse.TransactionId = refundRsp.TransactionId
	userCourse.OutTradeNo = refundRsp.OutTradeNo
	userCourse.OutRefundNo = refundRsp.OutRefundNo
	userCourse.RefundId = refundRsp.RefundId
	userCourse.RefundFee = refundRsp.RefundFee
	userCourse.Status = READING_COURSE_STATUS_REFUND
	userCourse.RefundWay = REFUND_WAY_FROM_TRANSACTION
	err = models.UpdateUserCourseTransactionRefundInfo(userCourse)
	if err != nil {
		holmes.Error("update user course refund info error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) getCourseDataStatistics(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &proto.GetCourseDataStatisticsReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	list, err := models.GetUserCourseDataStatistics(req.CourseType, req.Source, req.StartTime, req.EndTime)
	if err != nil {
		holmes.Error("get course data statistics list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = list
}

// wechat pay refund
func (self *ReadingHandler) getWechatPayOrder(transactionId string) (*mchpay.OrderQueryResponse, error) {
	orderQueryReq := &mchpay.OrderQueryRequest{
		TransactionId: transactionId,
	}
	orderQueryRsp, err := mchpay.OrderQuery2(self.mchClient, orderQueryReq)
	if err != nil {
		holmes.Error("order query error: %v", err)
		return nil, err
	}
	holmes.Debug("get wechat pay order success: %v", orderQueryRsp)
	return orderQueryRsp, nil
}

func (self *ReadingHandler) refundWechatPay(transactionId, outTradeNo, outRefundNo string, money int64) (*mchpay.RefundResponse, error) {
	refundReq := &mchpay.RefundRequest{
		TransactionId: transactionId,
		OutTradeNo:    outTradeNo,
		TotalFee:      money,
		RefundFee:     money,
		OutRefundNo:   outRefundNo,
	}
	refundRsp, err := mchpay.Refund2(self.mchClient, refundReq)
	if err != nil {
		holmes.Error("mch pay refund error: %v", err)
		return nil, err
	}
	return refundRsp, nil
}
