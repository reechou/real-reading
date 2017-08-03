package models

import (
	"time"
	
	"github.com/reechou/holmes"
)

const (
	COMMENT_STATUS_SHOW = 1
)

type CourseComment struct {
	ID                   int64  `xorm:"pk autoincr" json:"id"`
	UserId               int64  `xorm:"not null default 0 int index" json:"userId"`
	MonthCourseCatalogId int64  `xorm:"not null default 0 int index" json:"monthCourseCatalogId"`
	Comment              string `xorm:"not null default '' varchar(512)" json:"comment"`
	Status               int64  `xorm:"not null default 0 int index" json:"status"`
	CreatedAtStr         string `xorm:"not null default '' varchar(64)" json:"createdAtStr"`
	CreatedAt            int64  `xorm:"not null default 0 int" json:"createdAt"`
	ReplyUser            string `xorm:"not null default '' varchar(128)" json:"replyUser"`
	ReplyComment         string `xorm:"not null default '' varchar(512)" json:"replyComment"`
	ReplyAtStr           string `xorm:"not null default '' varchar(64)" json:"replyAtStr"`
	ReplyAt              int64  `xorm:"not null default 0 int" json:"replydAt"`
}

func CreateCourseComment(info *CourseComment) error {
	now := time.Now()
	info.CreatedAt = now.Unix()
	info.CreatedAtStr = now.Format("2006-01-02 15:04")
	
	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create course comment error: %v", err)
		return err
	}
	holmes.Info("create course comment[%v] success.", info)
	
	return nil
}

func UpdateCourseCommentStatus(info *CourseComment) error {
	_, err := x.ID(info.ID).Cols("status").Update(info)
	return err
}

func UpdateCourseCommentReply(info *CourseComment) error {
	now := time.Now()
	info.ReplyAt = now.Unix()
	info.ReplyAtStr = now.Format("2006-01-02 15:04")
	_, err := x.ID(info.ID).Cols("reply_user", "reply_comment", "reply_at_str", "reply_at").Update(info)
	return err
}
