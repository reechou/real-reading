package models

type ReadingPay struct {
	ID        int64  `xorm:"pk autoincr" json:"id"`
	OpenId    string `xorm:"not null default '' varchar(128) unique" json:"openId"`
	AppId     string `xorm:"not null default '' varchar(128)" json:"appId"`
	Name      string `xorm:"not null default '' varchar(256)" json:"name"`
	Phone     string `xorm:"not null default '' varchar(64)" json:"phone"`
	Wechat    string `xorm:"not null default '' varchar(64)" json:"tempUri"`
	Course    int64  `xorm:"not null default 0 int" json:"course"`
	Money     int64  `xorm:"not null default 0 int" json:"money"`
	Status    int64  `xorm:"not null default 0 int" json:"status"`
	CreatedAt int64  `xorm:"not null default 0 int" json:"createAt"`
	UpdatedAt int64  `xorm:"not null default 0 int" json:"-"`
}
