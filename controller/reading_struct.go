package controller

const (
	DEFAULT_ENROLL_NAME   = "请输入姓名"
	DEFAULT_ENROLL_MOBILE = "请输入手机号"
	DEFAULT_ENROLL_WECHAT = "微信号（选填）"
)

// about html tml struct
type ReadingEnrollUserInfo struct {
	NickName     string
	AvatarUrl    string
	OpenId       string
	EnrollName   string
	EnrollMobile string
	EnrollWechat string

	CourseType     int64
	CourseNum      int64
	IfCourseStart  int64
	StartTime      string
	EndTime        string
	IndexStartTime string

	WxJsApiParams
}

type WxJsApiParams struct {
	AppId     string `json:"appid"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}
