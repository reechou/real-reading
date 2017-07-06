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
	READING_COURSE_MANAGER_URI_PREFIX         = "coursemanager"
	READING_COURSE_MANAGER_URI_GET_BOOKS      = "getbooklist"
	READING_COURSE_MANAGER_URI_CREATE_BOOK    = "createbook"
	READING_COURSE_MANAGER_URI_DELETE_BOOK    = "deletebook"
	READING_COURSE_MANAGER_URI_UPDATE_BOOK    = "updatebook"
	READING_COURSE_MANAGER_URI_GET_CHAPTERS   = "getchapterlist"
	READING_COURSE_MANAGER_URI_CREATE_CHAPTER = "createchapter"
	READING_COURSE_MANAGER_URI_DELETE_CHAPTER = "deletechapter"
	READING_COURSE_MANAGER_URI_UPDATE_CHAPTER = "updatechapter"
	READING_COURSE_MANAGER_URI_GET_COURSES    = "getcourselist"
	READING_COURSE_MANAGER_URI_CREATE_COURSE  = "createcourse"
	READING_COURSE_MANAGER_URI_DELETE_COURSE  = "deletecourse"
	READING_COURSE_MANAGER_URI_UPDATE_COURSE  = "updatecourse"
)

func (self *ReadingHandler) courseManagerHandle(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	subPath := rr.Path[len(READING_COURSE_MANAGER_URI_PREFIX)+1:]
	rr.Params = strings.Split(subPath, "/")
	if len(rr.Params) == 0 {
		holmes.Error("path error: %s", rr.Path)
		return
	}
	switch rr.Params[0] {
	case READING_COURSE_MANAGER_URI_GET_BOOKS:
		self.courseManagerGetBookList(rr, w, r)
	case READING_COURSE_MANAGER_URI_CREATE_BOOK:
		self.courseManagerCreateBook(rr, w, r)
	case READING_COURSE_MANAGER_URI_DELETE_BOOK:
		self.courseManagerDeleteBook(rr, w, r)
	case READING_COURSE_MANAGER_URI_UPDATE_BOOK:
		self.courseManagerUpdateBook(rr, w, r)
	case READING_COURSE_MANAGER_URI_GET_CHAPTERS:
		self.courseManagerGetChapterList(rr, w, r)
	case READING_COURSE_MANAGER_URI_CREATE_CHAPTER:
		self.courseManagerCreateChapter(rr, w, r)
	case READING_COURSE_MANAGER_URI_DELETE_CHAPTER:
		self.courseManagerDeleteChapter(rr, w, r)
	case READING_COURSE_MANAGER_URI_UPDATE_CHAPTER:
		self.courseManagerUpdateChapter(rr, w, r)
	case READING_COURSE_MANAGER_URI_GET_COURSES:
		self.courseManagerGetCourseList(rr, w, r)
	case READING_COURSE_MANAGER_URI_CREATE_COURSE:
		self.courseManagerCreateCourse(rr, w, r)
	case READING_COURSE_MANAGER_URI_DELETE_COURSE:
		self.courseManagerDeleteCourse(rr, w, r)
	case READING_COURSE_MANAGER_URI_UPDATE_COURSE:
		self.courseManagerUpdateCourse(rr, w, r)
	}
}

func (self *ReadingHandler) courseManagerGetCourseList(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	bookList, err := models.GetCourseList()
	if err != nil {
		holmes.Error("get course list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = bookList
}

func (self *ReadingHandler) courseManagerCreateCourse(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.Course{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.CreateCourse(req)
	if err != nil {
		holmes.Error("create course error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerDeleteCourse(rr *HandlerRequest, w http.ResponseWriter, r *http.Request)  {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()
	
	req := &models.Course{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	err = models.DelCourse(req)
	if err != nil {
		holmes.Error("delete course error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerUpdateCourse(rr *HandlerRequest, w http.ResponseWriter, r *http.Request)  {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()
	
	req := &models.Course{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	if req.CourseType == 0 || req.CourseNum == 0 {
		holmes.Error("course type or num cannot be 0")
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	err = models.UpdateCourse(req)
	if err != nil {
		holmes.Error("update course error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}
