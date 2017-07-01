package models

import (
	"time"

	"github.com/reechou/holmes"
)

type User struct {
	ID          int64  `xorm:"pk autoincr" json:"id"`
	AppId       string `xorm:"not null default '' varchar(128)" json:"appId"`
	OpenId      string `xorm:"not null default '' varchar(128) unique" json:"openId"`
	Name        string `xorm:"not null default '' varchar(256)" json:"name"`
	AvatarUrl   string `xorm:"not null default '' varchar(256)" json:"avatarUrl"`
	RealName    string `xorm:"not null default '' varchar(256)" json:"realName"`
	Phone       string `xorm:"not null default '' varchar(64)" json:"phone"`
	Wechat      string `xorm:"not null default '' varchar(64)" json:"tempUri"`
	IfNotRemind int64  `xorm:"not null default 0 int" json:"ifNotRemind"`
	RemindTime  int64  `xorm:"not null default 0 int" json:"remindTime"`
	CreatedAt   int64  `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt   int64  `xorm:"not null default 0 int index" json:"-"`
}

func CreateUser(info *User) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create user error: %v", err)
		return err
	}
	holmes.Info("create user[%v] success.", info)

	return nil
}

func UpdateUserInfo(info *User) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("real_name", "phone", "wechat", "updated_at").Update(info)
	return err
}

func UpdateUserWxInfo(info *User) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("name", "avatar_url", "updated_at").Update(info)
	return err
}

func GetUserFromOpenid(info *User) (bool, error) {
	has, err := x.Where("open_id = ?", info.OpenId).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find user from openid[%s]", info.OpenId)
		return false, nil
	}
	return true, nil
}
