package controller

import (
	"github.com/reechou/real-reading/models"
)

const (
	COURSE_BOOK_LOCK = iota
	COURSE_BOOK_UNLOCK
)

const (
	COURSE_BOOK_CATALOG_LOCK = iota
	COURSE_BOOK_CATALOG_UNLOCK
)

// user course
type UserCourseDetail struct {
	TodayCatalogs  []models.CourseBookCatalogTime
	UserCourseList []models.UserCourseDetail
	UserId         int64
}

// course
type BookDetail struct {
	models.CourseBookDetail
	Status int
}

type MonthCourseDetail struct {
	models.MonthCourse
	Books []BookDetail
}

type CourseDetail struct {
	TodayCatalogs []models.CourseBookCatalogTime
	MonthCourses  []MonthCourseDetail
	UserId        int64
	CourseId      int64
}

// catalog
type CatalogDetail struct {
	models.MonthCourseCatalog
	Status   int
	TaskTime string
}

type CourseCatalogList struct {
	models.Book
	CatalogList []CatalogDetail
	UserId      int64
	CourseId    int64
	BookId      int64
}

// book detail
type CourseCatalogDetailList struct {
	models.Book
	models.MonthCourseCatalog
	models.MonthCourseCatalogAudio
	AudioTime         string
	PrevCatalogId     int64
	NextCatalogId     int64
	TaskTime          string
	ChapterDetailList []models.CourseBookCatalogDetail
	OpenId            string
	UserId            int64
	CourseId          int64
	BookId            int64
	CatalogId         int64
}

// course attendance
type UserCourseAttendance struct {
	AttendanceList []models.UserCourseList
	UserId         int64
}

type UserCourseAttendanceDetail struct {
	Course               models.Course
	AttendanceDetailList []models.UserCourseAttendance
	UserId               int64
	CourseId             int64

	StartYear        int
	StartMonth       int
	StartDay         int
	EndYear          int
	EndMonth         int
	EndDay           int
	NowYear          int
	NowMonth         int
	NowDay           int
	BeforeRenderAttr string
}
