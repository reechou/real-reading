package models

import (
	"fmt"
	"time"

	"github.com/reechou/holmes"
)

type Course struct {
	ID           int64  `xorm:"pk autoincr" json:"id"`
	CourseType   int64  `xorm:"not null default 0 int index" json:"courseType"`
	CourseNum    int64  `xorm:"not null default 0 int index" json:"courseNum"`
	Name         string `xorm:"not null default '' varchar(128)" json:"name"`
	Introduction string `xorm:"not null default '' varchar(512)" json:"introduction"`
	Cover        string `xorm:"not null default '' varchar(256)" json:"cover"`
	StartTime    int64  `xorm:"not null default 0 int index" json:"startTime"`
	EndTime      int64  `xorm:"not null default 0 int index" json:"endTime"`
	Money        int64  `xorm:"not null default 0 int" json:"money"`
	CreatedAt    int64  `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt    int64  `xorm:"not null default 0 int" json:"-"`
}

type MonthCourse struct {
	ID           int64  `xorm:"pk autoincr" json:"id"`
	CourseId     int64  `xorm:"not null default 0 int" json:"courseId"`
	IndexId      int64  `xorm:"not null default 0 int index" json:"indexId"`
	Year         int64  `xorm:"not null default 0 int" json:"year"`
	Month        int64  `xorm:"not null default 0 int" json:"month"`
	MonthEn      string `xorm:"not null default '' varchar(64)" json:"monthEn"`
	Title        string `xorm:"not null default '' varchar(256)" json:"title"`
	Introduction string `xorm:"not null default '' varchar(512)" json:"introduction"`
	CreatedAt    int64  `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt    int64  `xorm:"not null default 0 int" json:"-"`
}

type MonthCourseBook struct {
	ID            int64 `xorm:"pk autoincr" json:"id"`
	CourseId      int64 `xorm:"not null default 0 int" json:"courseId"`
	MonthCourseId int64 `xorm:"not null default 0 int index" json:"monthCourseId"`
	BookId        int64 `xorm:"not null default 0 int" json:"bookId"`
	IndexId       int64 `xorm:"not null default 0 int index" json:"indexId"`
	CreatedAt     int64 `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt     int64 `xorm:"not null default 0 int" json:"-"`
}

type MonthCourseCatalog struct {
	ID            int64  `xorm:"pk autoincr" json:"id"`
	CourseId      int64  `xorm:"not null default 0 int" json:"courseId"`
	MonthCourseId int64  `xorm:"not null default 0 int index" json:"monthCourseId"`
	BookId        int64  `xorm:"not null default 0 int index" json:"bookId"`
	IndexId       int64  `xorm:"not null default 0 int index" json:"indexId"`
	Title         string `xorm:"not null default '' varchar(256)" json:"title"`
	TaskTime      int64  `xorm:"not null default 0 int index" json:"taskTime"`
	CreatedAt     int64  `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt     int64  `xorm:"not null default 0 int" json:"-"`
}

type MonthCourseCatalogChapter struct {
	ID                   int64 `xorm:"pk autoincr" json:"id"`
	MonthCourseCatalogId int64 `xorm:"not null default 0 int index" json:"monthCourseCatalogId"`
	BookId               int64 `xorm:"not null default 0 int index" json:"bookId"`
	ChapterId            int64 `xorm:"not null default 0 int" json:"chapterId"`
	IndexId              int64 `xorm:"not null default 0 int index" json:"indexId"`
	CreatedAt            int64 `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt            int64 `xorm:"not null default 0 int" json:"-"`
}

type MonthCourseCatalogAudio struct {
	ID                   int64  `xorm:"pk autoincr" json:"id"`
	MonthCourseCatalogId int64  `xorm:"not null default 0 int index" json:"monthCourseCatalogId"`
	AudioTitle           string `xorm:"not null default '' varchar(128)" json:"AudioTitle"`
	AudioUrl             string `xorm:"not null default '' varchar(128)" json:"AudioUrl"`
	AudioTime            int64  `xorm:"not null default 0 int" json:"AudioTime"`
	CreatedAt            int64  `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt            int64  `xorm:"not null default 0 int" json:"-"`
}

func CreateCourse(info *Course) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create course error: %v", err)
		return err
	}
	holmes.Info("create course[%v] success.", info)

	return nil
}

func DelCourse(info *Course) error {
	if info.ID == 0 {
		return fmt.Errorf("del id cannot be nil.")
	}
	_, err := x.ID(info.ID).Delete(info)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCourse(info *Course) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("course_type", "course_num", "name", "introduction",
		"cover", "start_time", "end_time", "money", "updated_at").Update(info)
	return err
}

func CreateMonthCourse(info *MonthCourse) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create month course error: %v", err)
		return err
	}
	holmes.Info("create month course[%v] success.", info)

	return nil
}

func CreateMonthCourseBook(info *MonthCourseBook) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create month course book error: %v", err)
		return err
	}
	holmes.Info("create month course book[%v] success.", info)

	return nil
}

func CreateMonthCourseCatalog(info *MonthCourseCatalog) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create month course catalog error: %v", err)
		return err
	}
	holmes.Info("create month course catalog[%v] success.", info)

	return nil
}

func CreateMonthCourseCatalogChapter(info *MonthCourseCatalogChapter) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create month course catalog chapter error: %v", err)
		return err
	}
	holmes.Info("create month course catalog chapter[%v] success.", info)

	return nil
}

func CreateMonthCourseCatalogAudio(info *MonthCourseCatalogAudio) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create month course catalog audio error: %v", err)
		return err
	}
	holmes.Info("create month course catalog audio[%v] success.", info)

	return nil
}
