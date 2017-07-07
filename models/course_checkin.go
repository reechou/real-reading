package models

import (
	"time"

	"github.com/reechou/holmes"
)

type UserCourseCheckin struct {
	ID                   int64 `xorm:"pk autoincr" json:"id"`
	UserId               int64 `xorm:"not null default 0 int index" json:"userId"`
	CourseId             int64 `xorm:"not null default 0 int index" json:"courseId"`
	MonthCourseCatalogId int64 `xorm:"not null default 0 int index" json:"monthCourseCatalogId"`
	CreatedAt            int64 `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt            int64 `xorm:"not null default 0 int" json:"-"`
}

func CreateUserCourseCheckin(info *UserCourseCheckin) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create user course checkin error: %v", err)
		return err
	}
	holmes.Info("create user course checkin[%v] success.", info)

	return nil
}

func GetUserCourseCheckFromUCM(info *UserCourseCheckin) (bool, error) {
	has, err := x.Where("user_id = ?", info.UserId).And("course_id = ?", info.CourseId).And("month_course_catalog_id = ?", info.MonthCourseCatalogId).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}
	return true, nil
}

func GetUserCourseCheckinCount(userId, courseId int64) (int64, error) {
	count, err := x.Where("user_id = ?", userId).And("course_id = ?", courseId).Count(&UserCourseCheckin{})
	if err != nil {
		holmes.Error("get user course checkin list count error: %v", err)
		return 0, err
	}
	return count, nil
}
