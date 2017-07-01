package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/now"
	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/models"
	"github.com/reechou/real-reading/proto"
)

const (
	TEST_OPEN_ID = "oaKrZwsAF6pRX6z3Qn_EhIZ3DG90"
)

const (
	READING_COURSE_URI_PREFIX            = "course"
	READING_COURSE_URI_CATA_AUDIOS       = "catalogaudios"
	READING_COURSE_URI_CATA_SIGNIN       = "catalogsignin"
	READING_COURSE_URI_LIST              = "usercourselist"
	READING_COURSE_URI_INDEX             = "index"
	READING_COURSE_URI_BOOK_CATALOG      = "bookcatalog"
	READING_COURSE_URI_BOOK_DETAIL       = "bookdetail"
	READING_COURSE_URI_MYSELF            = "myself"
	READING_COURSE_URI_ATTENDANCE        = "attendance"
	READING_COURSE_URI_ATTENDANCE_DETAIL = "attendancedetail"
	READING_COURSE_URI_REMIND            = "remind"
)

func (self *ReadingHandler) courseHandle(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	subPath := rr.Path[len(READING_COURSE_URI_PREFIX)+1:]
	rr.Params = strings.Split(subPath, "/")
	if len(rr.Params) == 0 {
		holmes.Error("path error: %s", rr.Path)
		return
	}
	switch rr.Params[0] {
	case READING_COURSE_URI_CATA_AUDIOS:
		self.readingCourseCatalogAudios(rr, w, r)
	case READING_COURSE_URI_CATA_SIGNIN:
		self.readingCourseSignIn(rr, w, r)
	case READING_COURSE_URI_LIST:
		self.readingCourseList(rr, w, r)
	case READING_COURSE_URI_INDEX:
		self.readingCourseIndex(rr, w, r)
	case READING_COURSE_URI_BOOK_CATALOG:
		self.readingCourseCatalog(rr, w, r)
	case READING_COURSE_URI_MYSELF:
		self.readingCourseMyself(rr, w, r)
	case READING_COURSE_URI_BOOK_DETAIL:
		self.readingCourseChapterDetail(rr, w, r)
	case READING_COURSE_URI_ATTENDANCE:
		self.readingCourseAttendance(rr, w, r)
	case READING_COURSE_URI_ATTENDANCE_DETAIL:
		self.readingCourseAttendanceDetail(rr, w, r)
	case READING_COURSE_URI_REMIND:
		self.readingCourseRemind(rr, w, r)
	}
}

func (self *ReadingHandler) readingCourseCatalogAudios(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	if len(rr.Params) < 2 {
		holmes.Error("params error: %v", rr.Params)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	catalogId, err := strconv.ParseInt(rr.Params[1], 10, 0)
	if err != nil {
		holmes.Error("params[1][%s] strconv error: %v", rr.Params[1], err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	audios, err := models.GetCourseCatalogAudioList(catalogId)
	if err != nil {
		holmes.Error("get course catalog audio list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = audios
}

func (self *ReadingHandler) readingCourseSignIn(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &proto.ReadingCourseSignIn{}
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
		holmes.Error("get user error: %v", err)
		return
	}
	if !has {
		holmes.Error("cannot found this openid")
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	if user.ID != req.UserId {
		holmes.Error("openid userid[%d] != req.userid[%d]", user.ID, req.UserId)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	courseCheckin := &models.UserCourseCheckin{
		UserId:               user.ID,
		CourseId:             req.CourseId,
		MonthCourseCatalogId: req.CatalogId,
	}
	err = models.CreateUserCourseCheckin(courseCheckin)
	if err != nil {
		holmes.Error("create course checkin error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) readingCourseList(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	openId := TEST_OPEN_ID

	var err error
	userCourseList := new(UserCourseDetail)
	userCourseList.UserCourseList, err = models.GetUserCourse(openId)
	if err != nil {
		holmes.Error("get user course error: %v", err)
		return
	}
	if len(userCourseList.UserCourseList) == 0 {
		// TODO: redirect to course recommend
		return
	}
	userCourseList.UserId = userCourseList.UserCourseList[0].User.ID

	courseIds := make([]int64, 0)
	for _, v := range userCourseList.UserCourseList {
		courseIds = append(courseIds, v.Course.ID)
	}
	userCourseList.TodayCatalogs, err = models.GetCourseBookCatalogListFromTime(courseIds, now.BeginningOfDay().Unix())
	if err != nil {
		holmes.Error("get course book catalog list error: %v", err)
		return
	}

	renderView(w, "./views/course/course_list.html", userCourseList)
}

func (self *ReadingHandler) readingCourseIndex(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	if len(rr.Params) < 3 {
		holmes.Error("params error: %v", rr.Params)
		return
	}

	var err error
	courseDetail := new(CourseDetail)
	courseDetail.UserId, err = strconv.ParseInt(rr.Params[1], 10, 0)
	if err != nil {
		holmes.Error("params[1][%s] strconv error: %v", rr.Params[1], err)
		return
	}
	courseDetail.CourseId, err = strconv.ParseInt(rr.Params[2], 10, 0)
	if err != nil {
		holmes.Error("params[2][%s] strconv error: %v", rr.Params[2], err)
		return
	}

	courseDetail.TodayCatalogs, err = models.GetCourseBookFromTime(courseDetail.CourseId, now.BeginningOfDay().Unix())
	if err != nil {
		holmes.Error("get course book from time error: %v", err)
		return
	}

	monthCourses, err := models.GetMonthCourseList(courseDetail.CourseId)
	if err != nil {
		holmes.Error("get month course list error: %v", err)
		return
	}
	courseBookDetails, err := models.GetCourseBookDetail(courseDetail.CourseId)
	if err != nil {
		holmes.Error("get course book detail list error: %v", err)
		return
	}
	unlockBooks, err := models.GetMonthCourseBookUnlock(courseDetail.CourseId)
	if err != nil {
		holmes.Error("get month course book unlock error: %v", err)
		return
	}
	for _, v := range monthCourses {
		courseDetail.MonthCourses = append(courseDetail.MonthCourses, MonthCourseDetail{
			MonthCourse: v,
		})
	}
	for _, v := range courseBookDetails {
		for i := 0; i < len(courseDetail.MonthCourses); i++ {
			if courseDetail.MonthCourses[i].MonthCourse.ID == v.MonthCourseBook.MonthCourseId {
				bd := BookDetail{
					CourseBookDetail: v,
				}
				_, ok := unlockBooks[v.MonthCourseBook.BookId]
				if ok {
					bd.Status = COURSE_BOOK_UNLOCK
				}
				courseDetail.MonthCourses[i].Books = append(courseDetail.MonthCourses[i].Books, bd)
				break
			}
		}
	}

	renderView(w, "./views/course/course.html", courseDetail)
}

func (self *ReadingHandler) readingCourseCatalog(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	if len(rr.Params) < 4 {
		holmes.Error("params error: %v", rr.Params)
		return
	}

	var err error
	catalogList := new(CourseCatalogList)
	catalogList.UserId, err = strconv.ParseInt(rr.Params[1], 10, 0)
	if err != nil {
		holmes.Error("params[1][%s] strconv error: %v", rr.Params[1], err)
		return
	}
	catalogList.CourseId, err = strconv.ParseInt(rr.Params[2], 10, 0)
	if err != nil {
		holmes.Error("params[2][%s] strconv error: %v", rr.Params[2], err)
		return
	}
	catalogList.BookId, err = strconv.ParseInt(rr.Params[3], 10, 0)
	if err != nil {
		holmes.Error("params[3][%s] strconv error: %v", rr.Params[3], err)
		return
	}

	catalogList.Book.ID = catalogList.BookId
	has, err := models.GetBook(&catalogList.Book)
	if err != nil {
		holmes.Error("get book error: %v", err)
		return
	}
	if !has {
		holmes.Error("can not found this book[%d]", catalogList.BookId)
		return
	}
	courseCatalogList, err := models.GetCourseCatalogList(catalogList.CourseId, catalogList.BookId)
	if err != nil {
		holmes.Error("get course catalog list error: %v", err)
		return
	}
	now := time.Now().Unix()
	for _, v := range courseCatalogList {
		cd := CatalogDetail{
			MonthCourseCatalog: v,
		}
		if v.TaskTime <= now {
			cd.Status = COURSE_BOOK_CATALOG_UNLOCK
		}
		cd.TaskTime = time.Unix(v.TaskTime, 0).Format("2006-01-02")
		catalogList.CatalogList = append(catalogList.CatalogList, cd)
	}

	renderView(w, "./views/course/course_catalog.html", catalogList)
}

func (self *ReadingHandler) readingCourseChapterDetail(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	if len(rr.Params) < 5 {
		holmes.Error("params error: %v", rr.Params)
		return
	}

	var err error
	catalogDetailList := new(CourseCatalogDetailList)
	catalogDetailList.OpenId = TEST_OPEN_ID
	catalogDetailList.UserId, err = strconv.ParseInt(rr.Params[1], 10, 0)
	if err != nil {
		holmes.Error("params[1][%s] strconv error: %v", rr.Params[1], err)
		return
	}
	catalogDetailList.CourseId, err = strconv.ParseInt(rr.Params[2], 10, 0)
	if err != nil {
		holmes.Error("params[2][%s] strconv error: %v", rr.Params[2], err)
		return
	}
	catalogDetailList.BookId, err = strconv.ParseInt(rr.Params[3], 10, 0)
	if err != nil {
		holmes.Error("params[3][%s] strconv error: %v", rr.Params[3], err)
		return
	}
	catalogDetailList.CatalogId, err = strconv.ParseInt(rr.Params[4], 10, 0)
	if err != nil {
		holmes.Error("params[4][%s] strconv error: %v", rr.Params[4], err)
		return
	}

	catalogDetailList.MonthCourseCatalog.ID = catalogDetailList.CatalogId
	has, err := models.GetMonthCourseCatalog(&catalogDetailList.MonthCourseCatalog)
	if err != nil {
		holmes.Error("get month course catalog error: %v", err)
		return
	}
	if !has {
		holmes.Error("cannot found month course catalog from[%d]", catalogDetailList.CatalogId)
		return
	}
	catalogDetailList.TaskTime = time.Unix(catalogDetailList.MonthCourseCatalog.TaskTime, 0).Format("2006-01-02")
	catalogDetailList.ChapterDetailList, err = models.GetCourseBookCatalogDetailList(catalogDetailList.CatalogId)
	if err != nil {
		holmes.Error("get course book catalog detail list error: %v", err)
		return
	}
	holmes.Debug("%+v", catalogDetailList)

	renderView(w, "./views/course/course_detail.html", catalogDetailList)
}

func (self *ReadingHandler) readingCourseMyself(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	openId := TEST_OPEN_ID

	user := &models.User{
		OpenId: openId,
	}
	has, err := models.GetUserFromOpenid(user)
	if err != nil {
		holmes.Error("get user from openid[%s] error: %v", openId, err)
		return
	}
	if !has {
		holmes.Error("cannot found this user[%s]", openId)
		return
	}

	renderView(w, "./views/course/course_myself.html", user)
}

func (self *ReadingHandler) readingCourseAttendance(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	if len(rr.Params) < 2 {
		holmes.Error("params error: %v", rr.Params)
		return
	}

	var err error
	userCourseAttendanceList := new(UserCourseAttendance)
	userCourseAttendanceList.UserId, err = strconv.ParseInt(rr.Params[1], 10, 0)
	if err != nil {
		holmes.Error("params[1][%s] strconv error: %v", rr.Params[1], err)
		return
	}

	userCourseAttendanceList.AttendanceList, err = models.GetUserCourseList(userCourseAttendanceList.UserId)
	if err != nil {
		holmes.Error("get user course list error: %v", err)
		return
	}

	renderView(w, "./views/course/course_check_list.html", userCourseAttendanceList)
}

func (self *ReadingHandler) readingCourseAttendanceDetail(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	if len(rr.Params) < 3 {
		holmes.Error("params error: %v", rr.Params)
		return
	}

	var err error
	userCourseAttendanceDetailList := new(UserCourseAttendanceDetail)
	userCourseAttendanceDetailList.UserId, err = strconv.ParseInt(rr.Params[1], 10, 0)
	if err != nil {
		holmes.Error("params[1][%s] strconv error: %v", rr.Params[1], err)
		return
	}
	userCourseAttendanceDetailList.CourseId, err = strconv.ParseInt(rr.Params[2], 10, 0)
	if err != nil {
		holmes.Error("params[2][%s] strconv error: %v", rr.Params[2], err)
		return
	}

	userCourseAttendanceDetailList.Course.ID = userCourseAttendanceDetailList.CourseId
	has, err := models.GetCourse(&userCourseAttendanceDetailList.Course)
	if err != nil {
		holmes.Error("get course error: %v", err)
		return
	}
	if !has {
		holmes.Error("cannot found this course[%d]", userCourseAttendanceDetailList.CourseId)
		return
	}
	startTime := time.Unix(userCourseAttendanceDetailList.Course.StartTime, 0)
	endTime := time.Unix(userCourseAttendanceDetailList.Course.EndTime, 0)
	nowTime := time.Now()
	year, month, day := startTime.Date()
	userCourseAttendanceDetailList.StartYear = year
	userCourseAttendanceDetailList.StartMonth = int(month)
	userCourseAttendanceDetailList.StartDay = day
	year, month, day = endTime.Date()
	userCourseAttendanceDetailList.EndYear = year
	userCourseAttendanceDetailList.EndMonth = int(month)
	userCourseAttendanceDetailList.EndDay = day
	year, month, day = nowTime.Date()
	userCourseAttendanceDetailList.NowYear = year
	userCourseAttendanceDetailList.NowMonth = int(month)
	userCourseAttendanceDetailList.NowDay = day

	userCourseAttendanceDetailList.AttendanceDetailList, err = models.GetUserCourseCheckList(userCourseAttendanceDetailList.UserId, userCourseAttendanceDetailList.CourseId)
	if err != nil {
		holmes.Error("get user course check list error: %v", err)
		return
	}
	for i := 0; i < len(userCourseAttendanceDetailList.AttendanceDetailList); i++ {
		userCourseAttendanceDetailList.AttendanceDetailList[i].TaskTime *= 1000
	}

	renderView(w, "./views/course/course_check_detail.html", userCourseAttendanceDetailList)
}

func (self *ReadingHandler) readingCourseRemind(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {

	renderView(w, "./views/course/course_remind.html", nil)
}
