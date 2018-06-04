package controller

import (
	"github.com/reechou/real-reading/models"
)

type RegisterInfo struct {
	NickName     string
	AvatarUrl    string
	OpenId       string
	EnrollName   string
	EnrollMobile string
	EnrollWechat string
	Source       int
	UserId       int64

	Course models.Course
	Coupon models.Coupon

	RealPay string

	IfCourseStart  int64
	StartTime      string
	EndTime        string
	IndexStartTime string

	WxJsApiParams
}
