package models

import (
	"fmt"
	"time"

	"github.com/reechou/holmes"
)

type ReadingPay struct {
	ID        int64  `xorm:"pk autoincr" json:"id"`
	AppId     string `xorm:"not null default '' varchar(128)" json:"appId"`
	OpenId    string `xorm:"not null default '' varchar(128) unique" json:"openId"`
	Name      string `xorm:"not null default '' varchar(256)" json:"name"`
	AvatarUrl string `xorm:"not null default '' varchar(256)" json:"avatarUrl"`
	RealName  string `xorm:"not null default '' varchar(256)" json:"realName"`
	Phone     string `xorm:"not null default '' varchar(64)" json:"phone"`
	Wechat    string `xorm:"not null default '' varchar(64)" json:"tempUri"`
	Course    int64  `xorm:"not null default 0 int" json:"course"`
	Number    int64  `xorm:"not null default 0 int index" json:"number"`
	Money     int64  `xorm:"not null default 0 int" json:"money"`
	Status    int64  `xorm:"not null default 0 int index" json:"status"`
	CreatedAt int64  `xorm:"not null default 0 int" json:"createAt"`
	UpdatedAt int64  `xorm:"not null default 0 int index" json:"-"`
}

func CreateReadingPay(info *ReadingPay) error {
	if info.OpenId == "" {
		return fmt.Errorf("reading pay openid[%s] cannot be nil.", info.OpenId)
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create reading pay error: %v", err)
		return err
	}
	holmes.Info("create reading pay[%v] success.", info)

	return nil
}

func GetReadingPay(info *ReadingPay) (bool, error) {
	has, err := x.Where("open_id = ?", info.OpenId).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find reading pay from openid[%s]", info.OpenId)
		return false, nil
	}
	return true, nil
}

func UpdateReadingPayUserInfo(info *ReadingPay) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("real_name", "phone", "wechat", "updated_at").Update(info)
	return err
}

func UpdateReadingPayWxInfo(info *ReadingPay) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("name", "avatar_url", "updated_at").Update(info)
	return err
}

func UpdateReadingPayStatus(info *ReadingPay) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("money", "number", "status", "updated_at").Update(info)
	return err
}

func UpdateReadingPayStatusFromOpenId(info *ReadingPay) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Cols("money", "number", "status", "updated_at").Where("open_id = ?", info.OpenId).Update(info)
	return err
}

func GetReadingPayFromTime(fromTime int64) ([]ReadingPay, error) {
	var list []ReadingPay
	err := x.Where("status = 1").And("updated_at >= ?", fromTime).Find(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
