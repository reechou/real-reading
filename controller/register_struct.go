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
	
	Course models.Course
	
	IfCourseStart  int64
	StartTime      string
	EndTime        string
	IndexStartTime string
	
	WxJsApiParams
}
