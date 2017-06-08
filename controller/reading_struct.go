package controller

// about html tml struct
type ReadingEnrollUserInfo struct {
	NickName          string
	AvatarUrl         string
	OpenId            string
	WxJsApiParameters string
}

type WxJsApiParams struct {
	AppId     string `json:"appid"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}
