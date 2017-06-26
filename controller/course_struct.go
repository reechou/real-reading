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
}

// book detail
type CourseCatalogDetailList struct {
	models.Book
	models.MonthCourseCatalog
	TaskTime          string
	ChapterDetailList []models.CourseBookCatalogDetail
}
