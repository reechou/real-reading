package models

import (
	"github.com/reechou/holmes"
	"time"
)

type Coupon struct {
	ID         int64  `xorm:"pk autoincr" json:"id"`
	Name       string `xorm:"not null default '' varchar(128)" json:"name"`
	Desc       string `xorm:"not null default '' varchar(256)" json:"desc"`
	CdKey      string `xorm:"not null default '' varchar(128) unique" json:"cdkey"`
	Amount     int64  `xorm:"not null default 0 int" json:"amount"`
	IfUse      int64  `xorm:"not null default 0 int" json:"ifUse"`
	LockTime   int64  `xorm:"not null default 0 int" json:"lockTime"`
	CourseType int64  `xorm:"not null default 0 int" json:"courseType"`
	UserId     int64  `xorm:"not null default 0 int" json:"userId"`
	CreatedAt  int64  `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt  int64  `xorm:"not null default 0 int" json:"-"`
}

func CreateCouponList(list []Coupon) error {
	if len(list) == 0 {
		return nil
	}
	_, err := x.Insert(&list)
	if err != nil {
		holmes.Error("create Coupon list error: %v", err)
		return err
	}
	return nil
}

func GetCoupon(info *Coupon) (bool, error) {
	has, err := x.Where("cd_key = ?", info.CdKey).Get(info)
	if err != nil {
		return false, err
	}
	return has, nil
}

func UpdateCoupon(info *Coupon) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Incr("if_use").Update(info)
	return err
}

func UpdateCouponLock(info *Coupon) error {
	info.UpdatedAt = time.Now().Unix()
	info.LockTime = info.UpdatedAt + 300
	_, err := x.ID(info.ID).Cols("course_type", "user_id", "lock_time", "updated_at").Update(info)
	return err
}
