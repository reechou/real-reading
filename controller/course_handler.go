package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/chanxuehong/rand"
	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/jinzhu/now"
	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/models"
	"github.com/reechou/real-reading/proto"
)

const (
	TEST_OPEN_ID = "oaKrZwsAF6pRX6z3Qn_EhIZ3DG90"
)

const (
	READING_COURSE_URI_PREFIX          = "course"
	READING_COURSE_URI_CATA_AUDIOS     = "catalogaudios"
	READING_COURSE_URI_CATA_SIGNIN     = "catalogsignin"
	READING_COURSE_URI_CREATE_COMMENT  = "createcomment"
	READING_COURSE_URI_GET_COMMENT     = "getcommentlist"
	READING_COURSE_URI_UPDATE_COMMENT  = "updatecomment"
	READING_COURSE_URI_GET_ALL_COMMENT = "getallcomment"
	READING_COURSE_URI_REPLY_COMMENT   = "replycomment"
	// tmp html
	READING_COURSE_URI_LIST              = "usercourselist"
	READING_COURSE_URI_INDEX             = "index"
	READING_COURSE_URI_BOOK_CATALOG      = "bookcatalog"
	READING_COURSE_URI_BOOK_DETAIL       = "bookdetail"
	READING_COURSE_URI_MYSELF            = "myself"
	READING_COURSE_URI_ATTENDANCE        = "attendance"
	READING_COURSE_URI_ATTENDANCE_DETAIL = "attendancedetail"
	READING_COURSE_URI_REMIND            = "remind"
	READING_COURSE_URI_SHARE             = "share"
	READING_COURSE_URI_SHARE_TO          = "shareto"
)

func (self *ReadingHandler) courseHandle(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	subPath := rr.Path[len(READING_COURSE_URI_PREFIX)+1:]
	rr.Params = strings.Split(subPath, "/")
	if len(rr.Params) == 0 {
		holmes.Error("path error: %s", rr.Path)
		return
	}
	switch rr.Params[0] {
	// api
	case READING_COURSE_URI_CATA_AUDIOS:
		self.readingCourseCatalogAudios(rr, w, r)
	case READING_COURSE_URI_CATA_SIGNIN:
		self.readingCourseSignIn(rr, w, r)
	case READING_COURSE_URI_CREATE_COMMENT:
		self.readingCreateComment(rr, w, r)
	case READING_COURSE_URI_GET_COMMENT:
		self.readingGetCommentList(rr, w, r)
	case READING_COURSE_URI_UPDATE_COMMENT:
		self.readingUpdateComment(rr, w, r)
	case READING_COURSE_URI_GET_ALL_COMMENT:
		self.readingGetAllComment(rr, w, r)
	case READING_COURSE_URI_REPLY_COMMENT:
		self.readingReplyComment(rr, w, r)
	// tmp html
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
	case READING_COURSE_URI_SHARE:
		self.readingCourseShare(rr, w, r)
	case READING_COURSE_URI_SHARE_TO:
		self.readingCourseShareTo(rr, w, r)
	}
}

// api
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

func (self *ReadingHandler) readingCreateComment(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.CourseComment{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.CreateCourseComment(req)
	if err != nil {
		holmes.Error("create course comment error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	user := &models.User{ID: req.UserId}
	has, err := models.GetUser(user)
	if err != nil {
		holmes.Error("get user error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	if !has {
		holmes.Error("cannot found this user[%d]", req.UserId)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = user
}

func (self *ReadingHandler) readingGetCommentList(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &proto.GetCommentListReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	list, err := models.GetUserCourseComment(req.UserId, req.MonthCourseCatalogId)
	if err != nil {
		holmes.Error("get course comment error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = list
}

func (self *ReadingHandler) readingUpdateComment(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.CourseComment{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.UpdateCourseCommentStatus(req)
	if err != nil {
		holmes.Error("update course comment error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) readingGetAllComment(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &proto.GetAllCommentListReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	type CourseCommentList struct {
		Count int64                            `json:"count"`
		List  []models.UserCourseCommentDetail `json:"list"`
	}
	result := new(CourseCommentList)
	result.Count, err = models.GetCourseCommentCount()
	if err != nil {
		holmes.Error("get course comment count error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	result.List, err = models.GetCourseCommentList(req.Status, req.Offset, req.Num)
	if err != nil {
		holmes.Error("get course comment list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = result
}

func (self *ReadingHandler) readingReplyComment(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.CourseComment{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.UpdateCourseCommentReply(req)
	if err != nil {
		holmes.Error("update course comment reply error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

// check user
type UserInfo struct {
	OpenId    string
	Name      string
	AvatarUrl string
	Source    int
}

func (self *ReadingHandler) checkUserBase(w http.ResponseWriter, r *http.Request) (ui *UserInfo, ifRedirect bool) {
	ui = &UserInfo{}

	defer func() {
		if ui.OpenId != "" {
			// set cookie
			http.SetCookie(w, &http.Cookie{
				Name:    "user",
				Value:   ui.OpenId,
				Path:    "/",
				Expires: time.Now().Add(time.Hour),
			})
		}
	}()

	// get cookie
	cookie, err := r.Cookie("user")
	if err == nil {
		ui.OpenId = cookie.Value
		holmes.Debug("get user[%s] from cookie", ui.OpenId)
		return
	}

	//holmes.Debug("start to redirect oauth,")
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		holmes.Error("url parse query error: %v", err)
		return
	}

	code := queryValues.Get("code")
	if code == "" {
		state := string(rand.NewHex())
		redirectUrl := fmt.Sprintf("http://%s%s", r.Host, r.URL.String())
		AuthCodeURL := mpoauth2.AuthCodeURL(self.l.cfg.ReadingOauth.ReadingWxAppId,
			redirectUrl,
			self.l.cfg.ReadingOauth.ReadingOauth2ScopeBase, state)
		ifRedirect = true
		http.Redirect(w, r, AuthCodeURL, http.StatusFound)
		return
	}

	token, err := self.oauth2Client.ExchangeToken(code)
	if err != nil {
		ifRedirect = true
		http.Redirect(w, r, fmt.Sprintf("http://%s%s", r.Host, r.URL.Path), http.StatusFound)
		return
	}
	ui.OpenId = token.OpenId
	return
}

func (self *ReadingHandler) checkUser(w http.ResponseWriter, r *http.Request, ifNeedInfo bool) (ui *UserInfo, ifRedirect bool) {
	ui = &UserInfo{}

	defer func() {
		if ui.OpenId != "" {
			// set cookie
			http.SetCookie(w, &http.Cookie{
				Name:    "user",
				Value:   ui.OpenId,
				Path:    "/",
				Expires: time.Now().Add(time.Hour),
			})
		}
	}()

	// get cookie
	cookie, err := r.Cookie("user")
	if err == nil {
		ui.OpenId = cookie.Value
		holmes.Debug("get user[%s] from cookie", ui.OpenId)
		if ifNeedInfo {
			user := &models.User{
				OpenId: ui.OpenId,
			}
			has, err := models.GetUserFromOpenid(user)
			if err != nil {
				holmes.Error("get user from openid error: %v", err)
				goto NEED_OAUTH
			}
			if !has {
				goto NEED_OAUTH
			}
			ui.Name = user.Name
			ui.AvatarUrl = user.AvatarUrl
		}
		return
	}
NEED_OAUTH:
	//holmes.Debug("start to redirect oauth,")
	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
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
		ifRedirect = true
		http.Redirect(w, r, AuthCodeURL, http.StatusFound)
		return
	}
	src := queryValues.Get("src")
	if src != "" {
		ui.Source, err = strconv.Atoi(src)
		if err != nil {
			holmes.Error("strconv src[%s] error: %v", src, err)
		}
	}

	token, err := self.oauth2Client.ExchangeToken(code)
	if err != nil {
		ifRedirect = true
		http.Redirect(w, r, fmt.Sprintf("http://%s%s", r.Host, r.URL.Path), http.StatusFound)
		return
	}

	userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		holmes.Error("get user info error: %v", err)
		return
	}
	holmes.Debug("user info: %+v", userinfo)
	ui.OpenId = userinfo.OpenId
	ui.Name = userinfo.Nickname
	ui.AvatarUrl = userinfo.HeadImageURL

	user := &models.User{
		OpenId: ui.OpenId,
	}
	has, err := models.GetUserFromOpenid(user)
	if err != nil {
		holmes.Error("get user from openid error: %v", err)
		return
	}
	if !has {
		user.AppId = self.l.cfg.ReadingOauth.ReadingWxAppId
		user.Name = ui.Name
		user.AvatarUrl = ui.AvatarUrl
		user.Source = int64(ui.Source)
		err = models.CreateUser(user)
		if err != nil {
			holmes.Error("create user error: %v", err)
		}
	} else {
		if user.Name != ui.Name || user.AvatarUrl != ui.AvatarUrl {
			user.Name = ui.Name
			user.AvatarUrl = ui.AvatarUrl
			err = models.UpdateUserWxInfo(user)
			if err != nil {
				holmes.Error("update user wx info error: %v", err)
			}
		}
	}
	return
}

func (self *ReadingHandler) checkUserCourse(openId string, userId, courseId int64) bool {
	courseList, err := models.GetUserCourseFromOpenId(openId)
	if err != nil {
		holmes.Error("get user course error: %v", err)
		return false
	}
	//holmes.Debug("%+v user: %d course: %d", courseList, userId, courseId)
	if len(courseList) == 0 {
		return false
	}
	var userCheck, courseCheck bool
	for _, v := range courseList {
		if v.User.ID == userId {
			userCheck = true
		}
		if v.UserCourse.CourseId == courseId {
			courseCheck = true
		}
	}
	if userCheck && courseCheck {
		return true
	}
	return false
}

func (self *ReadingHandler) readingCourseError(w http.ResponseWriter, redirectUrl string) {
	courseError := &CourseError{
		RedirectUrl: redirectUrl,
	}
	renderView(w, "./views/course/course_error.html", courseError)
}

// template view
func (self *ReadingHandler) readingCourseList(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	//openId := TEST_OPEN_ID

	userinfo, ifRedirect := self.checkUser(w, r, false)
	if ifRedirect {
		return
	}

	var err error
	userCourseList := new(UserCourseDetail)
	userCourseList.UserCourseList, err = models.GetUserCourse(userinfo.OpenId)
	if err != nil {
		holmes.Error("get user course error: %v", err)
		return
	}
	if len(userCourseList.UserCourseList) == 0 {
		// TODO: redirect to course recommend
		self.readingCourseError(w, "/reading/signup")
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
	for i := 0; i < len(userCourseList.TodayCatalogs); i++ {
		ucc := &models.UserCourseCheckin{
			UserId:               userCourseList.UserId,
			CourseId:             userCourseList.TodayCatalogs[i].MonthCourseCatalog.CourseId,
			MonthCourseCatalogId: userCourseList.TodayCatalogs[i].MonthCourseCatalog.ID,
		}
		has, err := models.GetUserCourseCheckFromUCM(ucc)
		if err != nil {
			holmes.Error("get user course checkin error: %v", err)
			continue
		}
		if has {
			userCourseList.TodayCatalogs[i].IfCheck = 1
		}
	}

	renderView(w, "./views/course/course_list.html", userCourseList)
}

func (self *ReadingHandler) readingCourseIndex(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	userinfo, ifRedirect := self.checkUser(w, r, false)
	if ifRedirect {
		return
	}

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
	if !self.checkUserCourse(userinfo.OpenId, courseDetail.UserId, courseDetail.CourseId) {
		holmes.Error("user[%s] cannot found this courseid[%d]", userinfo.OpenId, courseDetail.UserId)
		// todo: redirect to sign
		self.readingCourseError(w, "/reading/signup")
		return
	}

	courseDetail.TodayCatalogs, err = models.GetCourseBookFromTime(courseDetail.CourseId, now.BeginningOfDay().Unix())
	if err != nil {
		holmes.Error("get course book from time error: %v", err)
		return
	}
	for i := 0; i < len(courseDetail.TodayCatalogs); i++ {
		ucc := &models.UserCourseCheckin{
			UserId:               courseDetail.UserId,
			CourseId:             courseDetail.TodayCatalogs[i].MonthCourseCatalog.CourseId,
			MonthCourseCatalogId: courseDetail.TodayCatalogs[i].MonthCourseCatalog.ID,
		}
		has, err := models.GetUserCourseCheckFromUCM(ucc)
		if err != nil {
			holmes.Error("get user course checkin error: %v", err)
			continue
		}
		if has {
			courseDetail.TodayCatalogs[i].IfCheck = 1
		}
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
	userinfo, ifRedirect := self.checkUser(w, r, false)
	if ifRedirect {
		return
	}

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
	if !self.checkUserCourse(userinfo.OpenId, catalogList.UserId, catalogList.CourseId) {
		holmes.Error("user[%s] cannot found this courseid[%d]", userinfo.OpenId, catalogList.UserId)
		// todo: redirect to sign
		self.readingCourseError(w, "/reading/signup")
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
	userinfo, ifRedirect := self.checkUser(w, r, false)
	if ifRedirect {
		return
	}

	if len(rr.Params) < 5 {
		holmes.Error("params error: %v", rr.Params)
		return
	}

	var err error
	catalogDetailList := new(CourseCatalogDetailList)
	catalogDetailList.OpenId = userinfo.OpenId
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
	if !self.checkUserCourse(userinfo.OpenId, catalogDetailList.UserId, catalogDetailList.CourseId) {
		holmes.Error("user[%s] cannot found this courseid[%d]", userinfo.OpenId, catalogDetailList.CourseId)
		//io.WriteString(w, MSG_ERROR_USER_COURSE_NOT_JOIN)
		// todo: redirect to sign
		self.readingCourseError(w, "/reading/signup")
		return
	}

	catalogDetailList.Book.ID = catalogDetailList.BookId
	_, err = models.GetBook(&catalogDetailList.Book)
	if err != nil {
		holmes.Error("get book error: %v", err)
		return
	}

	catalogs, err := models.GetMonthCourseCatalogFromBook(catalogDetailList.BookId)
	if err != nil {
		holmes.Error("get month course catalog from book error: %v", err)
		return
	}
	var ifHas bool
	var prevCatalogId, nextCatalogId int64
	for i := 0; i < len(catalogs); i++ {
		if catalogs[i].ID == catalogDetailList.CatalogId {
			ifHas = true
			catalogDetailList.MonthCourseCatalog = catalogs[i]
			if i < (len(catalogs) - 1) {
				nextCatalogId = catalogs[i+1].ID
			}
			break
		}
		prevCatalogId = catalogs[i].ID
	}
	if ifHas {
		catalogDetailList.PrevCatalogId = prevCatalogId
		catalogDetailList.NextCatalogId = nextCatalogId
	} else {
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
	}

	catalogDetailList.MonthCourseCatalogAudio.MonthCourseCatalogId = catalogDetailList.CatalogId
	has, err := models.GetMonthCourseCatalogAudio(&catalogDetailList.MonthCourseCatalogAudio)
	if err != nil {
		holmes.Error("get month course catalog audio error: %v", err)
		return
	}
	if has {
		catalogDetailList.AudioTime = fmt.Sprintf("%2d:%2d", catalogDetailList.MonthCourseCatalogAudio.AudioTime/60, catalogDetailList.MonthCourseCatalogAudio.AudioTime%60)
	}

	userCourseCheckin := &models.UserCourseCheckin{
		UserId:               catalogDetailList.UserId,
		CourseId:             catalogDetailList.CourseId,
		MonthCourseCatalogId: catalogDetailList.MonthCourseCatalog.ID,
	}
	has, err = models.GetUserCourseCheckFromUCM(userCourseCheckin)
	if err != nil {
		holmes.Error("get user course check from ucm error: %v", err)
		return
	}
	if has {
		catalogDetailList.IfChecked = 1
	}

	catalogDetailList.TaskTime = time.Unix(catalogDetailList.MonthCourseCatalog.TaskTime, 0).Format("2006-01-02")
	catalogDetailList.ChapterDetailList, err = models.GetCourseBookCatalogDetailList(catalogDetailList.CatalogId)
	if err != nil {
		holmes.Error("get course book catalog detail list error: %v", err)
		return
	}
	//holmes.Debug("%+v", catalogDetailList)

	renderView(w, "./views/course/course_detail.html", catalogDetailList)
}

func (self *ReadingHandler) readingCourseMyself(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	//openId := TEST_OPEN_ID

	userinfo, ifRedirect := self.checkUser(w, r, true)
	if ifRedirect {
		return
	}

	user := &models.User{
		OpenId: userinfo.OpenId,
	}
	has, err := models.GetUserFromOpenid(user)
	if err != nil {
		holmes.Error("get user from openid[%s] error: %v", userinfo.OpenId, err)
		return
	}
	if !has {
		holmes.Error("cannot found this user[%s]", userinfo.OpenId)
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
	userinfo, ifRedirect := self.checkUser(w, r, false)
	if ifRedirect {
		return
	}

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
	if !self.checkUserCourse(userinfo.OpenId, userCourseAttendanceDetailList.UserId, userCourseAttendanceDetailList.CourseId) {
		holmes.Error("user[%s] cannot found this courseid[%d]", userinfo.OpenId, userCourseAttendanceDetailList.UserId)
		// todo: redirect to sign
		self.readingCourseError(w, "/reading/signup")
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

func (self *ReadingHandler) readingCourseShare(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	userinfo, ifRedirect := self.checkUser(w, r, false)
	if ifRedirect {
		return
	}

	if len(rr.Params) < 3 {
		holmes.Error("params error: %v", rr.Params)
		return
	}

	var err error
	courseShare := new(CourseShare)
	courseShare.UserId, err = strconv.ParseInt(rr.Params[1], 10, 0)
	if err != nil {
		holmes.Error("params[1][%s] strconv error: %v", rr.Params[1], err)
		return
	}
	courseShare.CourseId, err = strconv.ParseInt(rr.Params[2], 10, 0)
	if err != nil {
		holmes.Error("params[2][%s] strconv error: %v", rr.Params[2], err)
		return
	}
	if !self.checkUserCourse(userinfo.OpenId, courseShare.UserId, courseShare.CourseId) {
		holmes.Error("user[%s] cannot found this courseid[%d]", userinfo.OpenId, courseShare.UserId)
		// todo: redirect to sign
		return
	}

	course := &models.Course{
		ID: courseShare.CourseId,
	}
	has, err := models.GetCourse(course)
	if err != nil {
		holmes.Error("get course error: %v", err)
		return
	}
	if !has {
		holmes.Error("cannot found this course[%d]", courseShare.CourseId)
		return
	}
	// add course type
	courseShare.CourseType = course.CourseType

	courseShare.DayNum, err = models.GetUserCourseCheckinCount(courseShare.UserId, courseShare.CourseId)
	if err != nil {
		holmes.Error("get user course checkin count error: %v", err)
		return
	}

	courseShare.JssdkInfo.Url = fmt.Sprintf("http://%s%s", r.Host, r.URL.String())
	self.l.wc.JssdkSign(&courseShare.JssdkInfo)
	courseShare.OpenId = userinfo.OpenId
	courseShare.AppId = self.l.cfg.ReadingOauth.ReadingWxAppId

	renderView(w, "./views/course/course_share.html", courseShare)
}

func (self *ReadingHandler) readingCourseShareTo(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	userinfo, ifRedirect := self.checkUserBase(w, r)
	if ifRedirect {
		return
	}

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		holmes.Error("url parse query error: %v", err)
		return
	}
	openid := queryValues.Get("openid")

	courseShareInfo := new(CourseShareInfo)
	user := &models.User{
		OpenId: openid,
	}
	has, err := models.GetUserFromOpenid(user)
	if err != nil {
		holmes.Error("get user from openid error: %v", err)
		return
	}
	if !has {
		holmes.Error("cannot found this man[%s]", openid)
		return
	}
	courseShareInfo.OpenId = userinfo.OpenId
	courseShareInfo.NickName = user.Name
	courseShareInfo.AvatarUrl = user.AvatarUrl
	courseShareInfo.AppId = self.l.cfg.ReadingOauth.ReadingWxAppId
	courseShareInfo.JssdkInfo.Url = fmt.Sprintf("http://%s%s", r.Host, r.URL.String())
	self.l.wc.JssdkSign(&courseShareInfo.JssdkInfo)

	renderView(w, "./views/course/course_share_to.html", courseShareInfo)
}
