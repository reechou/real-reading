package proto

// about request
type ReadingEnrollReq struct {
	Mobile   string `json:"mobile"`
	Realname string `json:"realname"`
	Weixin   string `json:"weixin"`
	OpenId   string `json:"openid"`
}

// about response
type ReadingPayToday struct {
	OrderNum int64 `json:"orderNum"`
	AllMoney int64 `json:"allMoney"`
}
