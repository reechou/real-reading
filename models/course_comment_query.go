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
		Find(&comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
