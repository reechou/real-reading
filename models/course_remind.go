package models

import (
	"time"

	"github.com/reechou/holmes"
)

type CourseRemind struct {
	ID            int64 `xorm:"pk autoincr" json:"id"`
	CourseId      int64 `xorm:"not null default 0 int index" json:"courseId"`
	RemindUserNum int64 `xorm:"not null default 0 int" json:"remindUserNum"`
	EndTime       int64 `xorm:"not null default 0 int" json:"endTime"`
	CreatedAt     int64 `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt     int64 `xorm:"not null default 0 int" json:"-"`
}

func CreateCourseRemind(info *CourseRemind) error {
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

func GetCourseRemind(info *CourseRemind) (bool, error) {
	has, err := x.Where("course_id = ?", info.CourseId).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find course remind from courseid[%d]", info.CourseId)
		return false, nil
	}
	return true, nil
}

func UpdateCourseRemindUserNum(info *CourseRemind) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("remind_user_num", "updated_at").Update(info)
	return err
}

func UpdateCourseRemindEndTime(info *CourseRemind) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("remind_user_num", "end_time", "updated_at").Update(info)
	return err
}
