package models

import (
	"fmt"
	"testing"

	"github.com/jinzhu/now"
	"github.com/reechou/real-reading/config"
)

func initRealeaseDb() {
	InitDB(&config.Config{
		DBInfo: config.DBInfo{
			User:   "fenxiao",
			Pass:   "c^ljPgOGafAlo%pd",
			Host:   "shanzhuan.mysql.rds.aliyuncs.com:3306",
			DBName: "real-reading",
		},
	})
}

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

func TestTransOldData(t *testing.T) {
	initRealeaseDb()
	
	oldList, err := GetReadingPayFromTime(0)
	if err != nil {
		fmt.Printf("get reading pay from time error: %v\n", err)
		return
	}
	for _, v := range oldList {
		user := &User{
			OpenId: v.OpenId,
		}
		has, err := GetUserFromOpenid(user)
		if err != nil {
			fmt.Printf("get user from openid error: %v\n", err)
			return
		}
		if has {
			if user.Phone == "" {
				user.AppId = v.AppId
				user.Name = v.Name
				user.AvatarUrl = v.AvatarUrl
				user.RealName = v.RealName
				user.Phone = v.Phone
				user.Wechat = v.Wechat
				err = UpdateUserAll(user)
				if err != nil {
					fmt.Printf("update user error: %v\n", err)
					return
				}
			}
		}
		//if !has {
		//	err = CreateUser(user)
		//	if err != nil {
		//		fmt.Printf("create user error: %v\n", err)
		//		return
		//	}
		//}
		//userCourse := &UserCourse{
		//	UserId: user.ID,
		//}
		//has, err = GetUserCourseFromUser(userCourse)
		//if err != nil {
		//	fmt.Printf("get user course from user error: %v\n", err)
		//	return
		//}
		//if !has {
		//	userCourse.CourseId = 1
		//	userCourse.Money = 19900
		//	userCourse.Status = 1
		//	userCourse.PayTime = v.CreatedAt
		//	err = CreateUserCourse(userCourse)
		//	if err != nil {
		//		fmt.Printf("create user course error: %v\n", err)
		//		return
		//	}
		//}
	}
}

func TestInit(t *testing.T) {
	initRealeaseDb()
	course := &Course{
		CourseType:   1,
		CourseNum:    15,
		Name:         "共读计划",
		Introduction: "一起共读计划",
		StartTime:    1499616000,
		EndTime:      1515513600,
		Money:        19900,
	}
	CreateCourse(course)
	course = &Course{
		CourseType:   1,
		CourseNum:    16,
		Name:         "共读计划",
		Introduction: "一起共读计划",
		StartTime:    1501516800,
		EndTime:      1517414400,
		Money:        19900,
	}
	CreateCourse(course)
}

func TestQuery(t *testing.T) {
	initDb()

	course := &Course{
		CourseType: 1,
	}
	GetCourseMaxNum(course)
	fmt.Println(course)
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

func TestCreateUserCourse(t *testing.T) {
	initDb()

	user := &User{
		AppId:     "wxb86df8e64a5ed4cd",
		OpenId:    "oaKrZwsAF6pRX6z3Qn_EhIZ3DG90",
		Name:      "Mr.REE",
		AvatarUrl: "http://wx.qlogo.cn/mmopen/ibmyOaFEgYk09HCYrBXA7PHZSuFjHINfuNxBlIOyvPibrU0hD87gTrGI2YuBTtGibHrxdTyzFAMFvWIPO5ekuhibzQ/0",
		RealName:  "周林栋",
		Phone:     "15994798218",
	}
	CreateUser(user)
	CreateUserCourse(&UserCourse{
		UserId:   user.ID,
		CourseId: 2,
		Money:    19900,
		Status:   1,
	})
}

func TestCreateCourseDetail(t *testing.T) {
	initDb()

	var courseId int64 = 2

	monthCourse := &MonthCourse{
		CourseId:     courseId,
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
		CourseId:      courseId,
		MonthCourseId: monthCourse.ID,
		BookId:        3,
	}
	err = CreateMonthCourseBook(monthCourseBook)
	if err != nil {
		fmt.Printf("create month course book error: %v\n", err)
		return
	}

	monthCourseCatalog := &MonthCourseCatalog{
		CourseId:      courseId,
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
		CourseId:     course.ID,
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
		CourseId:      course.ID,
		MonthCourseId: monthCourse.ID,
		BookId:        2,
	}
	err = CreateMonthCourseBook(monthCourseBook)
	if err != nil {
		fmt.Printf("create month course book error: %v\n", err)
		return
	}

	monthCourseCatalog := &MonthCourseCatalog{
		CourseId:      course.ID,
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

	courses, err := GetMonthCourseList(2)
	if err != nil {
		fmt.Printf("get course books error: %v\n", err)
		return
	}
	fmt.Println(courses)

	unlockBooks, _ := GetMonthCourseBookUnlock(2)
	fmt.Println(unlockBooks)

	monthBookList, err := GetCourseBooks(2)
	if err != nil {
		fmt.Printf("get course books error: %v\n", err)
		return
	}
	fmt.Println(monthBookList)

	bookDetail, err := GetCourseBookDetail(2)
	if err != nil {
		fmt.Printf("get course books detail error: %v\n", err)
		return
	}
	fmt.Println(bookDetail)

	todayCourse, err := GetCourseBookFromTime(2, now.BeginningOfDay().Unix())
	if err != nil {
		fmt.Printf("get course book from time: %v\n", err)
		return
	}
	fmt.Println(todayCourse)

	userCourse, err := GetUserCourse("oaKrZwsAF6pRX6z3Qn_EhIZ3DG90")
	fmt.Println(userCourse)
}
