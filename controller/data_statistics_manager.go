package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/models"
	"github.com/reechou/real-reading/proto"
)

const (
	DATA_STATISTICS_URI_PREFIX = "statistics"
)

const (
	DATA_STATISTICS_CREATE_COURSE_TYPE      = "createcoursetype"
	DATA_STATISTICS_GET_COURSE_TYPE_LIST    = "getcoursetypelist"
	DATA_STATISTICS_CREATE_COURSE_CHANNEL   = "createcoursechannel"
	DATA_STATISTICS_GET_COURSE_CHANNEL_LIST = "getcoursechannellist"
	DATA_STATISTICS_SET_USER_COURSE_REFUND  = "setusercourserefund"
	DATA_STATISTICS_GET_COURSE_STATISTICS   = "getcoursedatastatistics"
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
	case DATA_STATISTICS_GET_COURSE_CHANNEL_LIST:
		self.getCourseChannelList(rr, w, r)
	case DATA_STATISTICS_SET_USER_COURSE_REFUND:
		self.setUserCourseRefund(rr, w, r)
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
	
	req.Status = READING_COURSE_STATUS_REFUND
	err = models.UpdateUserCourseStatus(req)
	if err != nil {
		holmes.Error("update user course status of refund error: %v", err)
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
