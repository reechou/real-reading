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

type GetCommentListReq struct {
	UserId               int64 `json:"userId"`
	MonthCourseCatalogId int64 `json:"monthCourseCatalogId"`
}

type GetAllCommentListReq struct {
	Status int64 `json:"status"`
	Offset int64 `json:"offset"`
	Num    int64 `json:"num"`
}

// about course manager
type BookReq struct {
	BookId int64 `json:"bookId"`
}

type ChapterReq struct {
	ChapterId int64 `json:"chapterId"`
}

type CourseReq struct {
	CourseId int64 `json:"courseId"`
}

type MonthCourseReq struct {
	CourseId      int64 `json:"courseId"`
	MonthCourseId int64 `json:"monthCourseId"`
}

type MonthCourseBookReq struct {
	CourseId      int64 `json:"courseId"`
	MonthCourseId int64 `json:"monthCourseId"`
	BookId        int64 `json:"bookId"`
}

type MonthCourseCatalogReq struct {
	MonthCourseCatalogId int64 `json:"monthCourseCatalogId"`
}

// about response
type ReadingPayToday struct {
	OrderNum int64 `json:"orderNum"`
	AllMoney int64 `json:"allMoney"`
}
