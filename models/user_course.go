package models

import (
	"time"

	"github.com/reechou/holmes"
)

type UserCourse struct {
	ID        int64 `xorm:"pk autoincr" json:"id"`
	UserId    int64 `xorm:"not null default 0 int index" json:"userId"`
	CourseId  int64 `xorm:"not null default 0 int" json:"courseId"`
	Money     int64 `xorm:"not null default 0 int" json:"money"`
	Status    int64 `xorm:"not null default 0 int index" json:"status"`
	PayTime   int64 `xorm:"not null default 0 int index" json:"payTime"`
	CreatedAt int64 `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt int64 `xorm:"not null default 0 int" json:"-"`
}

func CreateUserCourse(info *UserCourse) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create user course error: %v", err)
		return err
	}
	holmes.Info("create user course[%v] success.", info)

	return nil
}

func GetUserCourseFromTime(fromTime int64) ([]UserCourse, error) {
	var list []UserCourse
	err := x.Where("pay_time >= ?", fromTime).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
