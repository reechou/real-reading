package models

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/config"
)

var x *xorm.Engine

func InitDB(cfg *config.Config) {
	var err error
	x, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4",
		cfg.DBInfo.User,
		cfg.DBInfo.Pass,
		cfg.DBInfo.Host,
		cfg.DBInfo.DBName))
	if err != nil {
		holmes.Fatal("Fail to init new engine: %v", err)
	}
	//x.SetLogger(nil)
	x.SetMapper(core.GonicMapper{})
	x.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	// if need show raw sql in log
	//x.ShowSQL(true)
	if cfg.IfShowSql {
		x.ShowSQL(true)
	}

	// sync tables
	if err = x.Sync2(new(ReadingPay),
		new(Book),
		new(Chapter),
		new(Course),
		new(MonthCourse),
		new(MonthCourseBook),
		new(MonthCourseCatalog),
		new(MonthCourseCatalogChapter),
		new(MonthCourseCatalogAudio),
		new(User),
		new(UserCourse),
		new(UserCourseCheckin),
		new(CourseRemind)); err != nil {
		holmes.Fatal("Fail to sync database: %v", err)
	}
}
