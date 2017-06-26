package controller

import (
	"net/http"
	"strings"
	"strconv"
	"time"
	
	"github.com/reechou/real-reading/models"
	"github.com/reechou/real-reading/proto"
	"github.com/jinzhu/now"
	"github.com/reechou/holmes"
)

const (
	READING_COURSE_URI_PREFIX       = "course"
	READING_COURSE_URI_CATA_AUDIOS  = "catalogaudios"
	READING_COURSE_URI_INDEX        = "index"
	READING_COURSE_URI_BOOK_CATALOG = "bookcatalog"
	READING_COURSE_URI_BOOK_DETAIL  = "bookdetail"
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
	case READING_COURSE_URI_INDEX:
		self.readingCourseIndex(rr, w, r)
	case READING_COURSE_URI_BOOK_CATALOG:
		self.readingCourseCatalog(rr, w, r)
	case READING_COURSE_URI_BOOK_DETAIL:
		self.readingCourseChapterDetail(rr, w, r)
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

func (self *ReadingHandler) readingCourseIndex(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	var courseNum int64 = 15
	
	var courseDetail CourseDetail
	var err error
	courseDetail.TodayCatalogs, err = models.GetCourseBookFromTime(courseNum, now.BeginningOfDay().Unix())
	if err != nil {
		holmes.Error("get course book from time error: %v", err)
		return
	}
	
	monthCourses, err := models.GetMonthCourseList(courseNum)
	if err != nil {
		holmes.Error("get month course list error: %v", err)
		return
	}
	courseBookDetails, err := models.GetCourseBookDetail(courseNum)
	if err != nil {
		holmes.Error("get course book detail list error: %v", err)
		return
	}
	unlockBooks, err := models.GetMonthCourseBookUnlock(courseNum)
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
	holmes.Debug("result: %+v %+v", courseDetail, unlockBooks)
	renderView(w, "./views/course/course.html", &courseDetail)
}

func (self *ReadingHandler) readingCourseCatalog(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	if len(rr.Params) < 3 {
		holmes.Error("params error: %v", rr.Params)
		return
	}
	
	var err error
	courseNum, err := strconv.ParseInt(rr.Params[1], 10, 0)
	if err != nil {
		holmes.Error("params[1][%s] strconv error: %v", rr.Params[1], err)
		return
	}
	bookId, err := strconv.ParseInt(rr.Params[2], 10, 0)
	if err != nil {
		holmes.Error("params[2][%s] strconv error: %v", rr.Params[2], err)
		return
	}
	
	catalogList := new(CourseCatalogList)
	catalogList.Book.ID = bookId
	has, err := models.GetBook(&catalogList.Book)
	if err != nil {
		holmes.Error("get book error: %v", err)
		return
	}
	if !has {
		holmes.Error("can not found this book[%d]", bookId)
		return
	}
	courseCatalogList, err := models.GetCourseCatalogList(courseNum, bookId)
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
	if len(rr.Params) < 4 {
		holmes.Error("params error: %v", rr.Params)
		return
	}
	
	var err error
	//courseNum, err := strconv.ParseInt(rr.Params[1], 10, 0)
	//if err != nil {
	//	holmes.Error("params[1][%s] strconv error: %v", rr.Params[1], err)
	//	return
	//}
	//bookId, err := strconv.ParseInt(rr.Params[2], 10, 0)
	//if err != nil {
	//	holmes.Error("params[2][%s] strconv error: %v", rr.Params[2], err)
	//	return
	//}
	catalogId, err := strconv.ParseInt(rr.Params[3], 10, 0)
	if err != nil {
		holmes.Error("params[3][%s] strconv error: %v", rr.Params[3], err)
		return
	}
	
	catalogDetailList := new(CourseCatalogDetailList)
	catalogDetailList.MonthCourseCatalog.ID = catalogId
	has, err := models.GetMonthCourseCatalog(&catalogDetailList.MonthCourseCatalog)
	if err != nil {
		holmes.Error("get month course catalog error: %v", err)
		return
	}
	if !has {
		holmes.Error("cannot found month course catalog from[%d]", catalogId)
		return
	}
	catalogDetailList.TaskTime = time.Unix(catalogDetailList.MonthCourseCatalog.TaskTime, 0).Format("2006-01-02")
	catalogDetailList.ChapterDetailList, err = models.GetCourseBookCatalogDetailList(catalogId)
	if err != nil {
		holmes.Error("get course book catalog detail list error: %v", err)
		return
	}
	holmes.Debug("%+v", catalogDetailList)
	
	renderView(w, "./views/course/course_detail.html", catalogDetailList)
}
