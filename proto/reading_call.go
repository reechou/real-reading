package proto

// about request
type ReadingEnrollReq struct {
	Mobile   string `json:"mobile"`
	Realname string `json:"realname"`
	Weixin   string `json:"weixin"`
	OpenId   string `json:"openid"`
	CourseId int64  `json:"courseid"`
}

type ReadingRefundReq struct {
	UserCourseId int64 `json:"userCourseId"`
}

type ReadingRefundFromTransactionReq struct {
	UserCourseId  int64  `json:"userCourseId"`
	TransactionId string `json:"transactionId"`
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

type GetGraduationReq struct {
	UserId   int64  `json:"userId"`
	UserName string `json:"userName"`
	CourseId int64  `json:"courseId"`
}

// about course manager
type BookReq struct {
	BookId int64 `json:"bookId"`
}

type BookCatalogReq struct {
	BookCatalogId int64 `json:"bookCatalogId"`
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

// about data statistics
type GetCourseChannelReq struct {
	CourseType int64 `json:"courseType"`
}

type GetCourseDataStatisticsReq struct {
	CourseType int64 `json:"courseType"`
	Source     int64 `json:"source"`
	StartTime  int64 `json:"startTime"`
	EndTime    int64 `json:"endTime"`
}

// coupon
type CreateCouponReq struct {
	Name   string
	Desc   string
	Amount int64
	Num    int
}

type GetCouponListReq struct {
	Offset int
	Num    int
}

// about response
type ReadingPayToday struct {
	OrderNum int64 `json:"orderNum"`
	AllMoney int64 `json:"allMoney"`
}
