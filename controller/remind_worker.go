package controller

import (
	"fmt"
	"sync"
	"time"

	"github.com/chanxuehong/wechat.v2/mp/message/template"
	"github.com/jinzhu/now"
	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/config"
	"github.com/reechou/real-reading/models"
	"github.com/robfig/cron"
)

const (
	DEFAULT_REMIND_WORKER_NUM = 100
)

type RemindInfo struct {
	Course  models.Course
	Catalog models.CourseBookCatalogTime
}

type RemindWorker struct {
	wg  sync.WaitGroup
	cfg *config.Config
	wtw *WechatTplWorker

	remindChan chan *RemindInfo
	WorkerNum  int

	stop chan struct{}
}

func NewRemindWorker(cfg *config.Config, wc *WechatController) *RemindWorker {
	rw := &RemindWorker{
		cfg:        cfg,
		remindChan: make(chan *RemindInfo, 256),
		stop:       make(chan struct{}),
	}
	if cfg.RemindWorkerNum == 0 {
		rw.WorkerNum = DEFAULT_REMIND_WORKER_NUM
	} else {
		rw.WorkerNum = cfg.RemindWorkerNum
	}
	if rw.cfg.RemindCronTime == "" {
		holmes.Error("course remind cron time is nil, maybe remind not used.")
		return rw
	}
	rw.wtw = NewWechatTplWorker(cfg, wc)

	go rw.do()

	for i := 0; i < rw.WorkerNum; i++ {
		rw.wg.Add(1)
		go rw.runWorker()
	}

	holmes.Debug("remind worker start..")

	return rw
}

func (self *RemindWorker) Stop() {
	close(self.stop)
	self.wg.Wait()
}

func (self *RemindWorker) do() {
	c := cron.New()
	c.AddFunc(self.cfg.RemindCronTime, self.runRemind)
	c.Start()

	select {
	case <-self.stop:
		c.Stop()
		return
	}
}

func (self *RemindWorker) runRemind() {
	nowTime := time.Now()
	if nowTime.Hour() < 8 {
		return
	}
	courseList, err := models.GetCourseListActive(nowTime.Unix())
	if err != nil {
		holmes.Error("get course list active error: %v", err)
		return
	}
	holmes.Debug("run remind active course list: %v", courseList)
	zeroNow := now.BeginningOfDay().Unix()
	for _, v := range courseList {
		catalogList, err := models.GetCourseBookFromTime(v.ID, zeroNow)
		if err != nil {
			holmes.Error("get course catalog list error: %v", err)
			continue
		}
		for _, v2 := range catalogList {
			ri := &RemindInfo{
				Course:  v,
				Catalog: v2,
			}
			self.courseRemind(ri)
		}
	}
}

func (self *RemindWorker) courseRemind(ri *RemindInfo) {
	select {
	case self.remindChan <- ri:
	case <-self.stop:
		return
	}
}

func (self *RemindWorker) runWorker() {
	for {
		select {
		case ri := <-self.remindChan:
			self.doRemind(ri)
		case <-self.stop:
			self.wg.Done()
			return
		}
	}
}

func (self *RemindWorker) doRemind(ri *RemindInfo) {
	courseRemind := &models.CourseRemind{
		CourseId: ri.Course.ID,
	}
	has, err := models.GetCourseRemind(courseRemind)
	if err != nil {
		holmes.Error("get course remind error: %v", err)
		return
	}
	if !has {
		err = models.CreateCourseRemind(courseRemind)
		if err != nil {
			holmes.Error("create course remind error: %v", err)
			return
		}
	}
	zeroToday := now.BeginningOfDay().Unix()
	if courseRemind.EndTime > zeroToday {
		holmes.Debug("courseid[%d] today has remind", ri.Course.ID)
		return
	}
	holmes.Info("course[%d] user remind start at: %s", ri.Course.ID, time.Now().Format("2006-01-02 15:04:05"))
	var offset int64
	if zeroToday < courseRemind.UpdatedAt {
		offset = courseRemind.RemindUserNum
	} else {
		courseRemind.RemindUserNum = 0
	}
	userList, err := models.GetCourseUserList(ri.Course.ID, offset)
	if err != nil {
		holmes.Error("get course user list error: %v", err)
		return
	}
	holmes.Debug("course[%d] get course user list from offset[%d] get user num: %d", ri.Course.ID, offset, len(userList))
	courseName := fmt.Sprintf("%s 第%d期", ri.Course.Name, ri.Course.CourseNum)
	catalogName := fmt.Sprintf("%s - %s", ri.Catalog.Book.BookName, ri.Catalog.Title)
	readingDate := time.Unix(ri.Catalog.TaskTime, 0).Format("2006.01.02")
	var sendUserNum int64
	for _, v := range userList {
		homeworkMsg := &HomeworkRemindTplMsg{
			First:    &template.DataItem{Value: "您今天的阅读任务已更新", Color: "#f76e6c"},
			Keyword1: &template.DataItem{Value: courseName},
			Keyword2: &template.DataItem{Value: catalogName},
			Keyword3: &template.DataItem{Value: readingDate},
			Remark:   &template.DataItem{Value: ">>>点击这里前往阅读<<<", Color: "#f76e6c"},
		}
		msg := &TplMsg{
			ToUser: v.OpenId,
			TplId:  TPL_ID_HOMEWORK_REMIND,
			Url: fmt.Sprintf("%s/reading/course/bookdetail/%d/%d/%d/%d", self.cfg.ReadingHost,
				v.User.ID,
				ri.Catalog.MonthCourseCatalog.CourseId,
				ri.Catalog.Book.ID,
				ri.Catalog.MonthCourseCatalog.ID),
			Data: homeworkMsg,
		}
		self.wtw.Send(msg)
		// save send user num
		sendUserNum++
		if sendUserNum >= 10 {
			courseRemind.RemindUserNum += sendUserNum
			err = models.UpdateCourseRemindUserNum(courseRemind)
			if err != nil {
				holmes.Error("update course remind usernum error: %v", err)
			}
			sendUserNum = 0
		}
	}
	// update remind end time
	if sendUserNum != 0 {
		courseRemind.RemindUserNum += sendUserNum
	}
	nowTime := time.Now()
	courseRemind.EndTime = nowTime.Unix()
	err = models.UpdateCourseRemindEndTime(courseRemind)
	if err != nil {
		holmes.Error("update course remind endtime error: %v", err)
	}
	holmes.Info("course[%d] user remind end at: %s", ri.Course.ID, nowTime.Format("2006-01-02 15:04:05"))
}
