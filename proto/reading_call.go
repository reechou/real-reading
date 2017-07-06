package proto

// about request
type ReadingEnrollReq struct {
	Mobile   string `json:"mobile"`
	Realname string `json:"realname"`
	Weixin   string `json:"weixin"`
	OpenId   string `json:"openid"`
}

type ReadingCourseSignIn struct {
	OpenId    string `json:"openId"`
	UserId    int64  `json:"userId"`
	CourseId  int64  `json:"courseId"`
	CatalogId int64  `json:"catalogId"`
}

// about course manager
type BookReq struct {
	BookId int64 `json:"bookId"`
}

// about response
type ReadingPayToday struct {
	OrderNum int64 `json:"orderNum"`
	AllMoney int64 `json:"allMoney"`
}
