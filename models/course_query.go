package models

import (
	"time"
)

func GetMonthCourseList(courseId int64) ([]MonthCourse, error) {
	var courses []MonthCourse
	err := x.Where("course_id = ?", courseId).OrderBy("index_id").Find(&courses)
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func GetMonthCourseBookUnlock(courseId int64) (map[int64]int, error) {
	var books []MonthCourseCatalog
	err := x.Distinct("book_id").Where("course_id = ?", courseId).
		And("task_time <= ?", time.Now().Unix()).
		Find(&books)
	if err != nil {
		return nil, err
	}
	bookIds := make(map[int64]int)
	for _, v := range books {
		bookIds[v.BookId] = 1
	}
	return bookIds, nil
}

func GetCourseCatalogList(courseId, bookId int64) ([]MonthCourseCatalog, error) {
	var catalogs []MonthCourseCatalog
	err := x.Where("course_id = ?", courseId).
		And("book_id = ?", bookId).
		OrderBy("index_id").
		Find(&catalogs)
	if err != nil {
		return nil, err
	}
	return catalogs, nil
}

func GetCourseCatalogAudioList(catalogId int64) ([]MonthCourseCatalogAudio, error) {
	var audios []MonthCourseCatalogAudio
	err := x.Where("month_course_catalog_id = ?", catalogId).
		Find(&audios)
	if err != nil {
		return nil, err
	}
	return audios, nil
}

func GetCourse(info *Course) (bool, error) {
	has, err := x.Id(info.ID).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}
	return true, nil
}

func GetCourseMaxNum(info *Course) (bool, error) {
	has, err := x.Where("course_type = ?", info.CourseType).Desc("course_num").Limit(1).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}
	return true, nil
}

func GetCourseFromCourse(info *Course) (bool, error) {
	has, err := x.Where("course_type = ?", info.CourseType).And("course_num = ?", info.CourseNum).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}
	return true, nil
}

func GetMonthCourseCatalog(info *MonthCourseCatalog) (bool, error) {
	has, err := x.Id(info.ID).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}
	return true, nil
}

type MonthCourseBooks struct {
	MonthCourse     `xorm:"extends"`
	MonthCourseBook `xorm:"extends"`
}

func (MonthCourseBooks) TableName() string {
	return "month_course"
}

func GetCourseBooks(courseId int64) ([]MonthCourseBooks, error) {
	monthBooks := make([]MonthCourseBooks, 0)
	err := x.Join("LEFT", "month_course_book", "month_course_book.month_course_id = month_course.id").
		Where("month_course.course_id = ?", courseId).
		Find(&monthBooks)
	if err != nil {
		return nil, err
	}
	return monthBooks, nil
}

type CourseBookDetail struct {
	MonthCourseBook `xorm:"extends"`
	Book            `xorm:"extends"`
}

func (CourseBookDetail) TableName() string {
	return "month_course_book"
}

func GetCourseBookDetail(courseId int64) ([]CourseBookDetail, error) {
	bookDetail := make([]CourseBookDetail, 0)
	err := x.Join("LEFT", "book", "month_course_book.book_id = book.id").
		Where("month_course_book.course_id = ?", courseId).
		OrderBy("month_course_book.index_id").
		Find(&bookDetail)
	if err != nil {
		return nil, err
	}
	return bookDetail, nil
}

type CourseBookCatalogTime struct {
	MonthCourseCatalog `xorm:"extends"`
	Book               `xorm:"extends"`
}

func (CourseBookCatalogTime) TableName() string {
	return "month_course_catalog"
}

func GetCourseBookFromTime(courseId, day int64) ([]CourseBookCatalogTime, error) {
	bookCatalogBooks := make([]CourseBookCatalogTime, 0)
	err := x.Join("LEFT", "book", "month_course_catalog.book_id = book.id").
		Where("month_course_catalog.course_id = ?", courseId).
		And("task_time = ?", day).
		Find(&bookCatalogBooks)
	if err != nil {
		return nil, err
	}
	return bookCatalogBooks, nil
}

func GetCourseBookCatalogListFromTime(courseIds []int64, day int64) ([]CourseBookCatalogTime, error) {
	bookCatalogBooks := make([]CourseBookCatalogTime, 0)
	err := x.Join("LEFT", "book", "month_course_catalog.book_id = book.id").
		In("month_course_catalog.course_id", courseIds).
		And("task_time = ?", day).
		Find(&bookCatalogBooks)
	if err != nil {
		return nil, err
	}
	return bookCatalogBooks, nil
}

type CourseBookCatalogDetail struct {
	MonthCourseCatalogChapter `xorm:"extends"`
	Chapter                   `xorm:"extends"`
}

func (CourseBookCatalogDetail) TableName() string {
	return "month_course_catalog_chapter"
}

func GetCourseBookCatalogDetailList(monthCourseCatalogId int64) ([]CourseBookCatalogDetail, error) {
	catalogDetailList := make([]CourseBookCatalogDetail, 0)
	err := x.Join("LEFT", "chapter", "month_course_catalog_chapter.chapter_id = chapter.id").
		Where("month_course_catalog_chapter.month_course_catalog_id = ?", monthCourseCatalogId).
		OrderBy("month_course_catalog_chapter.index_id").
		Find(&catalogDetailList)
	if err != nil {
		return nil, err
	}
	return catalogDetailList, nil
}

type UserCourseList struct {
	UserCourse `xorm:"extends"`
	Course     `xorm:"extends"`
}

func (UserCourseList) TableName() string {
	return "user_course"
}

func GetUserCourseList(userId int64) ([]UserCourseList, error) {
	userCourseList := make([]UserCourseList, 0)
	err := x.Join("LEFT", "course", "user_course.course_id = course.id").
		Where("user_course.user_id = ?", userId).
		Find(&userCourseList)
	if err != nil {
		return nil, err
	}
	return userCourseList, nil
}

type UserCourseDetail struct {
	UserCourse `xorm:"extends"`
	User       `xorm:"extends"`
	Course     `xorm:"extends"`
}

func (UserCourseDetail) TableName() string {
	return "user_course"
}

func GetUserCourse(openId string) ([]UserCourseDetail, error) {
	userCourseList := make([]UserCourseDetail, 0)
	err := x.Join("LEFT", "user", "user_course.user_id = user.id").
		Join("LEFT", "course", "user_course.course_id = course.id").
		Where("user.open_id = ?", openId).
		Find(&userCourseList)
	if err != nil {
		return nil, err
	}
	return userCourseList, nil
}

type UserCourseAttendance struct {
	UserCourseCheckin  `xorm:"extends"`
	MonthCourseCatalog `xorm:"extends"`
}

func (UserCourseAttendance) TableName() string {
	return "user_course_checkin"
}

func GetUserCourseCheckList(userId, courseId int64) ([]UserCourseAttendance, error) {
	userCourseAttendance := make([]UserCourseAttendance, 0)
	err := x.Join("LEFT", "month_course_catalog", "user_course_checkin.month_course_catalog_id = month_course_catalog.id").
		Where("user_course_checkin.user_id = ?", userId).
		And("user_course_checkin.course_id = ?", courseId).
		Find(&userCourseAttendance)
	if err != nil {
		return nil, err
	}
	return userCourseAttendance, nil
}
