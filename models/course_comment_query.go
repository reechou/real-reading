package models

type UserCourseComment struct {
	CourseComment `xorm:"extends" json:"courseComment"`
	User          `xorm:"extends" json:"user"`
}

func (UserCourseComment) TableName() string {
	return "course_comment"
}

func GetUserCourseComment(userId, monthCourseCatalogId int64) ([]UserCourseComment, error) {
	comments := make([]UserCourseComment, 0)
	err := x.Join("LEFT", "user", "course_comment.user_id = user.id").
		Where("course_comment.month_course_catalog_id = ?", monthCourseCatalogId).
		And("course_comment.status = 1").
		Limit(100).
		Find(&comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

type UserCourseCommentDetail struct {
	CourseComment      `xorm:"extends" json:"courseComment"`
	User               `xorm:"extends" json:"user"`
	MonthCourseCatalog `xorm:"extends" json:"monthCourseCatalog"`
}

func (UserCourseCommentDetail) TableName() string {
	return "course_comment"
}

func GetCourseCommentCount() (int64, error) {
	count, err := x.Join("LEFT", "user", "course_comment.user_id = user.id").
		Join("LEFT", "month_course_catalog", "course_comment.month_course_catalog_id = month_course_catalog.id").
		And("course_comment.status = 0").
		Count(&UserCourseCommentDetail{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetCourseCommentList(status, offset, num int64) ([]UserCourseCommentDetail, error) {
	comments := make([]UserCourseCommentDetail, 0)
	err := x.Join("LEFT", "user", "course_comment.user_id = user.id").
		Join("LEFT", "month_course_catalog", "course_comment.month_course_catalog_id = month_course_catalog.id").
		And("course_comment.status = ?", status).
		Limit(int(num), int(offset)).
		Find(&comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
