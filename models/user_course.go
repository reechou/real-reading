package models

import (
	"time"

	"github.com/reechou/holmes"
)

type UserCourse struct {
	ID            int64  `xorm:"pk autoincr" json:"id"`
	UserId        int64  `xorm:"not null default 0 int index" json:"userId"`
	CourseId      int64  `xorm:"not null default 0 int index" json:"courseId"`
	CourseType    int64  `xorm:"not null default 0 int index" json:"courseType"`
	Money         int64  `xorm:"not null default 0 int" json:"money"`
	Status        int64  `xorm:"not null default 0 int index" json:"status"`
	PayTime       int64  `xorm:"not null default 0 int index" json:"payTime"`
	Source        int64  `xorm:"not null default 0 int index" json:"source"`
	OutTradeNo    string `xorm:"not null default '' varchar(128)" json:"outTradeNo"`
	TransactionId string `xorm:"not null default '' varchar(128)" json:"transactionId"`
	OutRefundNo   string `xorm:"not null default '' varchar(128)" json:"outRefundNo"`
	RefundId      string `xorm:"not null default '' varchar(128)" json:"refundId"`
	RefundFee     int64  `xorm:"not null default 0 int" json:"refundFee"`
	RefundWay     int64  `xorm:"not null default 0 int" json:"refundWay"`
	CreatedAt     int64  `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt     int64  `xorm:"not null default 0 int" json:"-"`
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

func GetUserCourseFromId(info *UserCourse) (bool, error) {
	has, err := x.Id(info.ID).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find user course from id[%d]", info.ID)
		return false, nil
	}
	return true, nil
}

func GetUserCourseFromUser(info *UserCourse) (bool, error) {
	has, err := x.Where("user_id = ?", info.UserId).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}
	return true, nil
}

func GetUserCourseFromTime(fromTime int64) ([]UserCourse, error) {
	var list []UserCourse
	err := x.Where("pay_time >= ?", fromTime).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func UpdateUserCourseStatus(info *UserCourse) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("status", "updated_at").Update(info)
	return err
}

func UpdateUserCourseRefundInfo(info *UserCourse) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("out_refund_no", "refund_id", "refund_fee", "status", "refund_way", "updated_at").Update(info)
	return err
}

func UpdateUserCourseTransactionRefundInfo(info *UserCourse) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("out_trade_no", "transaction_id", "out_refund_no", "refund_id", "refund_fee", "status", "refund_way", "updated_at").Update(info)
	return err
}
