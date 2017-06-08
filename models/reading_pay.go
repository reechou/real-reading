package models

import (
	"fmt"
	"time"
	
	"github.com/reechou/holmes"
)

type ReadingPay struct {
	ID        int64  `xorm:"pk autoincr" json:"id"`
	OpenId    string `xorm:"not null default '' varchar(128) unique" json:"openId"`
	AppId     string `xorm:"not null default '' varchar(128)" json:"appId"`
	Name      string `xorm:"not null default '' varchar(256)" json:"name"`
	AvatarUrl string `xorm:"not null default '' varchar(256)" json:"avatarUrl"`
	Phone     string `xorm:"not null default '' varchar(64)" json:"phone"`
	Wechat    string `xorm:"not null default '' varchar(64)" json:"tempUri"`
	Course    int64  `xorm:"not null default 0 int" json:"course"`
	Money     int64  `xorm:"not null default 0 int" json:"money"`
	Status    int64  `xorm:"not null default 0 int" json:"status"`
	CreatedAt int64  `xorm:"not null default 0 int" json:"createAt"`
	UpdatedAt int64  `xorm:"not null default 0 int" json:"-"`
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
