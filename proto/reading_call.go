package proto

// about request
type ReadingEnrollReq struct {
	Mobile   string `json:"mobile"`
	Realname string `json:"realname"`
	Weixin   string `json:"weixin"`
	OpenId   string `json:"openid"`
}
