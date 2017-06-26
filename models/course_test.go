package models

import (
	"fmt"
	"testing"

	"github.com/reechou/real-reading/config"
	"github.com/jinzhu/now"
)

func initDb() {
	InitDB(&config.Config{
		DBInfo: config.DBInfo{
			User:   "youzan",
			Pass:   "111111",
			Host:   "127.0.0.1:3306",
			DBName: "real_reading",
		},
	})
}

func TestCreateBook(t *testing.T) {
	initDb()

	book := &Book{
		BookName: "高效能人士的七个习惯",
		Author:   "(美)史蒂芬柯维",
		Abstract: "这是要不永恒的畅销书,里程碑式的著作,总销量超过2500万册,被评为“有史以来最具影响力的10大管理类图书之一”。",
		Cover:    "https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1498041630158&di=cd70acb1ae5353c31d9b975603e5aa76&imgtype=0&src=http%3A%2F%2Fpic.baike.soso.com%2Fp%2F20140421%2Fbki-20140421170747-225819815.jpg",
	}
	err := CreateBook(book)
	if err != nil {
		fmt.Printf("create book error: %v\n", err)
		return
	}

	err = CreateChapter(&Chapter{
		BookId:  book.ID,
		IndexId: 1,
		Title:   "前言",
		Cover:   "https://ss3.bdstatic.com/70cFv8Sh_Q1YnxGkpoWK1HF6hhy/it/u=1323751765,1242219830&fm=26&gp=0.jpg",
		Content: `<p>据不完全统计，我发现大家对工具类的内容（例如<a href="http://www.jianshu.com/p/6ca85cbf8210" target="_blank">高效软件</a>、<a href="http://www.jianshu.com/p/2ec2c101ac56" target="_blank">PPT制作</a>、成为表格女王、摄影技巧等）都非常感兴趣；可如果东西罗列太多，也只是一键收藏而已。</p><p>在很多场景下，掌握一些科学的方法，比使用再多的工具要好使得多。今天我们要讲的，就是大名鼎鼎的「番茄工作法」，这个方法的科学性和实用性不容小觑。</p>`,
	})
	if err != nil {
		fmt.Printf("create chapter error: %v\n", err)
		return
	}
}

func TestCreateCourseDetail(t *testing.T) {
	initDb()
	
	var courseNum int64 = 15
	
	monthCourse := &MonthCourse{
		CourseNum:    courseNum,
		Year:         2017,
		Month:        8,
		MonthEn:      "Aug",
		Title:        "职场规划",
		Introduction: "职场里的是是非非",
	}
	err := CreateMonthCourse(monthCourse)
	if err != nil {
		fmt.Printf("create month course error: %v\n", err)
		return
	}
	
	monthCourseBook := &MonthCourseBook{
		CourseNum:     courseNum,
		MonthCourseId: monthCourse.ID,
		BookId:        3,
	}
	err = CreateMonthCourseBook(monthCourseBook)
	if err != nil {
		fmt.Printf("create month course book error: %v\n", err)
		return
	}
	
	monthCourseCatalog := &MonthCourseCatalog{
		CourseNum:     courseNum,
		MonthCourseId: monthCourse.ID,
		BookId:        3,
		Title:         "职场的陷阱",
		TaskTime:      1498147200,
	}
	err = CreateMonthCourseCatalog(monthCourseCatalog)
	if err != nil {
		fmt.Printf("create month course catalog error: %v\n", err)
		return
	}
	
	monthCourseCatalogChapter := &MonthCourseCatalogChapter{
		MonthCourseCatalogId: monthCourseCatalog.ID,
		BookId:               3,
		ChapterId:            2,
	}
	err = CreateMonthCourseCatalogChapter(monthCourseCatalogChapter)
	if err != nil {
		fmt.Printf("create month course catalog chapter error: %v\n", err)
		return
	}
}

func TestCreateCourse(t *testing.T) {
	initDb()

	course := &Course{
		CourseNum: 15,
		Name:      "共读计划第15期",
	}
	err := CreateCourse(course)
	if err != nil {
		fmt.Printf("cretea course: %v\n", err)
		return
	}

	monthCourse := &MonthCourse{
		CourseNum:    course.CourseNum,
		Year:         2017,
		Month:        7,
		MonthEn:      "Jul",
		Title:        "时间管理",
		Introduction: "在正确的时间做该做的事",
	}
	err = CreateMonthCourse(monthCourse)
	if err != nil {
		fmt.Printf("create month course error: %v\n", err)
		return
	}

	monthCourseBook := &MonthCourseBook{
		CourseNum:     course.CourseNum,
		MonthCourseId: monthCourse.ID,
		BookId:        2,
	}
	err = CreateMonthCourseBook(monthCourseBook)
	if err != nil {
		fmt.Printf("create month course book error: %v\n", err)
		return
	}

	monthCourseCatalog := &MonthCourseCatalog{
		CourseNum:     course.CourseNum,
		MonthCourseId: monthCourse.ID,
		BookId:        2,
		Title:         "一次只做一件事",
		TaskTime:      now.BeginningOfDay().Unix(),
	}
	err = CreateMonthCourseCatalog(monthCourseCatalog)
	if err != nil {
		fmt.Printf("create month course catalog error: %v\n", err)
		return
	}

	monthCourseCatalogChapter := &MonthCourseCatalogChapter{
		MonthCourseCatalogId: monthCourseCatalog.ID,
		BookId:               2,
		ChapterId:            1,
	}
	err = CreateMonthCourseCatalogChapter(monthCourseCatalogChapter)
	if err != nil {
		fmt.Printf("create month course catalog chapter error: %v\n", err)
		return
	}
}

func TestQueryCourse(t *testing.T) {
	initDb()
	
	courses, err := GetMonthCourseList(15)
	if err != nil {
		fmt.Printf("get course books error: %v\n", err)
		return
	}
	fmt.Println(courses)
	
	unlockBooks, _ := GetMonthCourseBookUnlock(15)
	fmt.Println(unlockBooks)
	
	monthBookList, err := GetCourseBooks(15)
	if err != nil {
		fmt.Printf("get course books error: %v\n", err)
		return
	}
	fmt.Println(monthBookList)
	
	bookDetail, err := GetCourseBookDetail(15)
	if err != nil {
		fmt.Printf("get course books detail error: %v\n", err)
		return
	}
	fmt.Println(bookDetail)
	
	todayCourse, err := GetCourseBookFromTime(15, now.BeginningOfDay().Unix())
	if err != nil {
		fmt.Printf("get course book from time: %v\n", err)
		return
	}
	fmt.Println(todayCourse)
}
