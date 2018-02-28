package models

import (
	"fmt"
	"time"

	"github.com/reechou/holmes"
)

func GetMonthCourseFromCourse(courseId int64) ([]MonthCourse, error) {
	var mcs []MonthCourse
	err := x.Where("course_id = ?", courseId).Find(&mcs)
	if err != nil {
		return nil, err
	}
	return mcs, nil
}

func GetMonthCourseBookFromMonthCourse(monthCourseId int64) ([]MonthCourseBook, error) {
	var mcbs []MonthCourseBook
	err := x.Where("month_course_id = ?", monthCourseId).Find(&mcbs)
	if err != nil {
		return nil, err
	}
	return mcbs, nil
}

func GetMonthCourseCatalogFromCourse(courseId, monthCourseId, bookId int64) ([]MonthCourseCatalog, error) {
	var mccs []MonthCourseCatalog
	err := x.Where("course_id = ?", courseId).And("month_course_id = ?", monthCourseId).And("book_id = ?", bookId).Find(&mccs)
	if err != nil {
		return nil, err
	}
	return mccs, nil
}

func GetMonthCourseCatalogChapterFromCatalog(monthCourseCatalogId int64) ([]MonthCourseCatalogChapter, error) {
	var mcccs []MonthCourseCatalogChapter
	err := x.Where("month_course_catalog_id = ?", monthCourseCatalogId).Find(&mcccs)
	if err != nil {
		return nil, err
	}
	return mcccs, nil
}

func GetMonthCourseCatalogAudioFromCatalog(monthCourseCatalogId int64) ([]MonthCourseCatalogAudio, error) {
	var mccas []MonthCourseCatalogAudio
	err := x.Where("month_course_catalog_id = ?", monthCourseCatalogId).Find(&mccas)
	if err != nil {
		return nil, err
	}
	return mccas, nil
}
