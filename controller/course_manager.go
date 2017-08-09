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
	READING_COURSE_MANAGER_URI_PREFIX             = "coursemanager"
	READING_COURSE_MANAGER_URI_GET_BOOKS          = "getbooklist"
	READING_COURSE_MANAGER_URI_CREATE_BOOK        = "createbook"
	READING_COURSE_MANAGER_URI_DELETE_BOOK        = "deletebook"
	READING_COURSE_MANAGER_URI_UPDATE_BOOK        = "updatebook"
	READING_COURSE_MANAGER_URI_GET_CHAPTER        = "getchapter"
	READING_COURSE_MANAGER_URI_GET_CHAPTERS       = "getchapterlist"
	READING_COURSE_MANAGER_URI_CREATE_CHAPTER     = "createchapter"
	READING_COURSE_MANAGER_URI_DELETE_CHAPTER     = "deletechapter"
	READING_COURSE_MANAGER_URI_UPDATE_CHAPTER     = "updatechapter"
	READING_COURSE_MANAGER_URI_UPDATE_CHAPTER_RE  = "updatechapterre"
	READING_COURSE_MANAGER_URI_GET_BCATALOGS      = "getbookcatalogs"
	READING_COURSE_MANAGER_URI_CREATE_BCATALOG    = "createbookcatalog"
	READING_COURSE_MANAGER_URI_DELETE_BCATALOG    = "deletebookcatalog"
	READING_COURSE_MANAGER_URI_UPDATE_BCATALOG    = "updatebookcatalog"
	READING_COURSE_MANAGER_URI_GET_BCCHAPTERS     = "getbcchapters"
	READING_COURSE_MANAGER_URI_CREATE_BCCHAPTERS  = "createbcchapters"
	READING_COURSE_MANAGER_URI_DELETE_BCCHAPTER   = "deletebcchapter"
	READING_COURSE_MANAGER_URI_GET_COURSES        = "getcourselist"
	READING_COURSE_MANAGER_URI_CREATE_COURSE      = "createcourse"
	READING_COURSE_MANAGER_URI_DELETE_COURSE      = "deletecourse"
	READING_COURSE_MANAGER_URI_UPDATE_COURSE      = "updatecourse"
	READING_COURSE_MANAGER_URI_GET_MCOURSES       = "getmcourselist"
	READING_COURSE_MANAGER_URI_CREATE_MCOURSE     = "createmcourse"
	READING_COURSE_MANAGER_URI_DELETE_MCOURSE     = "deletemcourse"
	READING_COURSE_MANAGER_URI_UPDATE_MCOURSE     = "updatemcourse"
	READING_COURSE_MANAGER_URI_GET_MCBOOKS        = "getmcbooklist"
	READING_COURSE_MANAGER_URI_CREATE_MCBOOK      = "createmcbook"
	READING_COURSE_MANAGER_URI_DELETE_MCBOOK      = "deletemcbook"
	READING_COURSE_MANAGER_URI_UPDATE_MCBOOK      = "updatemcbook"
	READING_COURSE_MANAGER_URI_GET_MCBCATALOGS    = "getmcbcataloglist"
	READING_COURSE_MANAGER_URI_CREATE_MCBCATALOG  = "createmcbcatalog"
	READING_COURSE_MANAGER_URI_DELETE_MCBCATALOG  = "deletemcbcatalog"
	READING_COURSE_MANAGER_URI_UPDATE_MCBCATALOG  = "updatemcbcatalog"
	READING_COURSE_MANAGER_URI_GET_MCBCCHAPTERS   = "getmcbcchapterlist"
	READING_COURSE_MANAGER_URI_CREATE_MCBCCHAPTER = "createmcbcchapter"
	READING_COURSE_MANAGER_URI_DELETE_MCBCCHAPTER = "deletemcbcchapter"
)

var (
	MonthConst map[int64]string = map[int64]string{
		1:  "Jan.",
		2:  "Feb.",
		3:  "Mar.",
		4:  "Apr.",
		5:  "May",
		6:  "Jun.",
		7:  "Jul.",
		8:  "Aug.",
		9:  "Sep.",
		10: "Oct.",
		11: "Nov.",
		12: "Dec.",
	}
)

func (self *ReadingHandler) courseManagerHandle(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	subPath := rr.Path[len(READING_COURSE_MANAGER_URI_PREFIX)+1:]
	rr.Params = strings.Split(subPath, "/")
	if len(rr.Params) == 0 {
		holmes.Error("path error: %s", rr.Path)
		return
	}
	switch rr.Params[0] {
	// about book manager
	case READING_COURSE_MANAGER_URI_GET_BOOKS:
		self.courseManagerGetBookList(rr, w, r)
	case READING_COURSE_MANAGER_URI_CREATE_BOOK:
		self.courseManagerCreateBook(rr, w, r)
	case READING_COURSE_MANAGER_URI_DELETE_BOOK:
		self.courseManagerDeleteBook(rr, w, r)
	case READING_COURSE_MANAGER_URI_UPDATE_BOOK:
		self.courseManagerUpdateBook(rr, w, r)
	case READING_COURSE_MANAGER_URI_GET_CHAPTER:
		self.courseManagerGetChapter(rr, w, r)
	case READING_COURSE_MANAGER_URI_GET_CHAPTERS:
		self.courseManagerGetChapterList(rr, w, r)
	case READING_COURSE_MANAGER_URI_CREATE_CHAPTER:
		self.courseManagerCreateChapter(rr, w, r)
	case READING_COURSE_MANAGER_URI_DELETE_CHAPTER:
		self.courseManagerDeleteChapter(rr, w, r)
	case READING_COURSE_MANAGER_URI_UPDATE_CHAPTER:
		self.courseManagerUpdateChapter(rr, w, r)
	// about book catalog
	case READING_COURSE_MANAGER_URI_GET_BCATALOGS:
		self.courseManagerGetBookCatalogList(rr, w, r)
	case READING_COURSE_MANAGER_URI_CREATE_BCATALOG:
		self.courseManagerCreateBookCatalog(rr, w, r)
	case READING_COURSE_MANAGER_URI_DELETE_BCATALOG:
		self.courseManagerDeleteBookCatalog(rr, w, r)
	case READING_COURSE_MANAGER_URI_UPDATE_BCATALOG:
		self.courseManagerUpdateBookCatalog(rr, w, r)
	case READING_COURSE_MANAGER_URI_GET_BCCHAPTERS:
		self.courseManagerGetBookCatalogChapterList(rr, w, r)
	case READING_COURSE_MANAGER_URI_CREATE_BCCHAPTERS:
		self.courseManagerCreateBookCatalogChapter(rr, w, r)
	case READING_COURSE_MANAGER_URI_DELETE_BCCHAPTER:
		self.courseManagerDeleteBookCatalogChapter(rr, w, r)
	// about course manager
	case READING_COURSE_MANAGER_URI_GET_COURSES:
		self.courseManagerGetCourseList(rr, w, r)
	case READING_COURSE_MANAGER_URI_CREATE_COURSE:
		self.courseManagerCreateCourse(rr, w, r)
	case READING_COURSE_MANAGER_URI_DELETE_COURSE:
		self.courseManagerDeleteCourse(rr, w, r)
	case READING_COURSE_MANAGER_URI_UPDATE_COURSE:
		self.courseManagerUpdateCourse(rr, w, r)
	case READING_COURSE_MANAGER_URI_GET_MCOURSES:
		self.courseManagerGetMonthCourseList(rr, w, r)
	case READING_COURSE_MANAGER_URI_CREATE_MCOURSE:
		self.courseManagerCreateMonthCourse(rr, w, r)
	case READING_COURSE_MANAGER_URI_DELETE_MCOURSE:
		self.courseManagerDeleteMonthCourse(rr, w, r)
	case READING_COURSE_MANAGER_URI_UPDATE_MCOURSE:
		self.courseManagerUpdateMonthCourse(rr, w, r)
	case READING_COURSE_MANAGER_URI_GET_MCBOOKS:
		self.courseManagerGetMonthCourseBookList(rr, w, r)
	case READING_COURSE_MANAGER_URI_CREATE_MCBOOK:
		self.courseManagerCreateMonthCourseBook(rr, w, r)
	case READING_COURSE_MANAGER_URI_DELETE_MCBOOK:
		self.courseManagerDeleteMonthCourseBook(rr, w, r)
	case READING_COURSE_MANAGER_URI_UPDATE_MCBOOK:
		self.courseManagerUpdateMonthCourseBook(rr, w, r)
	case READING_COURSE_MANAGER_URI_GET_MCBCATALOGS:
		self.courseManagerGetMonthCourseCatalogList(rr, w, r)
	case READING_COURSE_MANAGER_URI_CREATE_MCBCATALOG:
		self.courseManagerCreateMonthCourseCatalog(rr, w, r)
	case READING_COURSE_MANAGER_URI_DELETE_MCBCATALOG:
		self.courseManagerDeleteMonthCourseCatalog(rr, w, r)
	case READING_COURSE_MANAGER_URI_UPDATE_MCBCATALOG:
		self.courseManagerUpdateMonthCourseCatalog(rr, w, r)
	case READING_COURSE_MANAGER_URI_GET_MCBCCHAPTERS:
		self.courseManagerGetMonthCourseCatalogChapterList(rr, w, r)
	case READING_COURSE_MANAGER_URI_CREATE_MCBCCHAPTER:
		self.courseManagerCreateMonthCourseCatalogChapterList(rr, w, r)
	case READING_COURSE_MANAGER_URI_DELETE_MCBCCHAPTER:
		self.courseManagerDeleteMonthCourseCatalogChapter(rr, w, r)
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

func (self *ReadingHandler) courseManagerDeleteCourse(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
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

func (self *ReadingHandler) courseManagerUpdateCourse(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
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

func (self *ReadingHandler) courseManagerGetMonthCourseList(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &proto.CourseReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	list, err := models.GetMonthCourseList(req.CourseId)
	if err != nil {
		holmes.Error("get month course list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = list
}

func (self *ReadingHandler) courseManagerCreateMonthCourse(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.MonthCourse{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	if req.Month > 12 {
		holmes.Debug("month[%d] cannot be > 12", req.Month)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	req.MonthEn = MonthConst[req.Month]

	err = models.CreateMonthCourse(req)
	if err != nil {
		holmes.Error("create month course error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerDeleteMonthCourse(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.MonthCourse{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.DelMonthCourse(req)
	if err != nil {
		holmes.Error("delete month course error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerUpdateMonthCourse(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.MonthCourse{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	if req.CourseId == 0 {
		holmes.Error("course id cannot be 0")
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.UpdateMonthCourse(req)
	if err != nil {
		holmes.Error("update month course error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerGetMonthCourseBookList(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &proto.MonthCourseReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	list, err := models.GetMonthCourseBookList(req.CourseId, req.MonthCourseId)
	if err != nil {
		holmes.Error("get month course book list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = list
}

func (self *ReadingHandler) courseManagerCreateMonthCourseBook(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.MonthCourseBook{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.CreateMonthCourseBook(req)
	if err != nil {
		holmes.Error("create month course book error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerDeleteMonthCourseBook(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.MonthCourseBook{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.DelMonthCourseBook(req)
	if err != nil {
		holmes.Error("delete month course book error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerUpdateMonthCourseBook(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.MonthCourseBook{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	if req.CourseId == 0 {
		holmes.Error("course id cannot be 0")
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.UpdateMonthCourseBook(req)
	if err != nil {
		holmes.Error("update month course book error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerGetMonthCourseCatalogList(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &proto.MonthCourseBookReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	list, err := models.GetMonthCourseCatalogList(req.CourseId, req.MonthCourseId, req.BookId)
	if err != nil {
		holmes.Error("get month course catalog list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = list
}

func (self *ReadingHandler) courseManagerCreateMonthCourseCatalog(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.MonthCourseCatalog{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.CreateMonthCourseCatalog(req)
	if err != nil {
		holmes.Error("create month course catalog error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerDeleteMonthCourseCatalog(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.MonthCourseCatalog{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.DelMonthCourseCatalog(req)
	if err != nil {
		holmes.Error("delete month course catalog error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerUpdateMonthCourseCatalog(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.MonthCourseCatalog{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	if req.CourseId == 0 {
		holmes.Error("course id cannot be 0")
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.UpdateMonthCourseCatalog(req)
	if err != nil {
		holmes.Error("update month course catalog error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerGetMonthCourseCatalogChapterList(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &proto.MonthCourseCatalogReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	list, err := models.GetMonthCourseCatalogChapterList(req.MonthCourseCatalogId)
	if err != nil {
		holmes.Error("get month course catalog chapter list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = list
}

func (self *ReadingHandler) courseManagerCreateMonthCourseCatalogChapterList(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	var req []models.MonthCourseCatalogChapter
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.CreateMonthCourseCatalogChapterList(req)
	if err != nil {
		holmes.Error("create month course catalog chapter list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerDeleteMonthCourseCatalogChapter(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.MonthCourseCatalogChapter{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.DelMonthCourseCatalogChapter(req)
	if err != nil {
		holmes.Error("delete month course catalog chapter error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}
