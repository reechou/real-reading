package logic

import (
	"time"

	"github.com/jinzhu/now"
	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/config"
	"github.com/reechou/real-reading/models"
)

var (
	MonthConst map[int64]string = map[int64]string{
		1:  "Jan.",
		2:  "Feb.",
		3:  "Mar.",
		4:  "Apr.",
		5:  "May",
		6:  "Jun.",
		7:  "Jul.",
		8:  "Aug.",
		9:  "Sep.",
		10: "Oct.",
		11: "Nov.",
		12: "Dec.",
	}
)

type UpdateCourseLogic struct {
	cfg *config.Config
}

func NewUpdateCourseLogic(cfg *config.Config) *UpdateCourseLogic {
	ucl := &UpdateCourseLogic{cfg: cfg}
	ucl.init()

	return ucl
}

func (ucl *UpdateCourseLogic) init() {
	models.InitDB(ucl.cfg)
}

func (ucl *UpdateCourseLogic) Run() {
	defer holmes.Start(holmes.LogFilePath("./log"),
		holmes.EveryDay,
		holmes.AlsoStdout,
		holmes.DebugLevel).Stop()

	now.FirstDayMonday = true

	// get dst course info
	dstCourse := &models.Course{ID: ucl.cfg.CopyCourse.DstCourseId}
	has, err := models.GetCourse(dstCourse)
	if err != nil {
		holmes.Error("get course error: %v", err)
		return
	}
	if !has {
		holmes.Error("dst course[%d] not found", ucl.cfg.CopyCourse.DstCourseId)
		return
	}

	// get month course list
	mcs, err := models.GetMonthCourseFromCourse(ucl.cfg.CopyCourse.SrcCourseId)
	if err != nil {
		holmes.Error("get src month course from course error: %v", err)
		return
	}
	mcs2, err := models.GetMonthCourseFromCourse(ucl.cfg.CopyCourse.DstCourseId)
	if err != nil {
		holmes.Error("get dst month course from course error: %v", err)
		return
	}
	// handle month course
	if len(mcs2) > len(mcs) {
		holmes.Error("len(dst)[%d] > len(src)[%d]", len(mcs2), len(mcs))
		return
	} else if len(mcs2) < len(mcs) {
		// check
		var i = 0
		year := int64(time.Unix(dstCourse.StartTime, 0).Year())
		month := int64(time.Unix(dstCourse.StartTime, 0).Month())
		//holmes.Debug("year: %d month: %d start: %v", year, month, time.Unix(dstCourse.StartTime, 0))
		month--
		var idx int64
		for ; i < len(mcs2); i++ {
			if mcs2[i].IndexId != mcs[i].IndexId && mcs2[i].Title != mcs[i].Title {
				holmes.Error("check month course error, maybe not same course.")
				return
			}
			year = mcs2[i].Year
			month = mcs2[i].Month
			idx = mcs2[i].IndexId
		}
		// add month course
		for ; i < len(mcs); i++ {
			if month >= 12 {
				month = 1
				year++
			} else {
				month++
			}
			idx++
			mc := models.MonthCourse{
				CourseId:     ucl.cfg.CopyCourse.DstCourseId,
				IndexId:      idx,
				Year:         year,
				Month:        month,
				MonthEn:      MonthConst[month],
				Title:        mcs[i].Title,
				Introduction: mcs[i].Introduction,
			}
			if err = models.CreateMonthCourse(&mc); err != nil {
				holmes.Error("create month course error: %v", err)
				return
			}
			mcs2 = append(mcs2, mc)
		}
	}

	dstCourseTaskTimeOfWeekBegin := now.New(time.Unix(dstCourse.StartTime, 0).UTC()).BeginningOfWeek().Unix()
	holmes.Debug("start: %d begin: %d", dstCourse.StartTime, dstCourseTaskTimeOfWeekBegin)
	var taskWeekIdx = -1
	// get month course book list
	for i := 0; i < len(mcs); i++ {
		mcbs, err := models.GetMonthCourseBookFromMonthCourse(mcs[i].ID)
		if err != nil {
			holmes.Error("get month course book from src id[%d] error: %v", mcs[i].ID, err)
			return
		}
		mcbs2, err := models.GetMonthCourseBookFromMonthCourse(mcs2[i].ID)
		if err != nil {
			holmes.Error("get month course book from dst id[%d] error: %v", mcs2[i].ID, err)
			return
		}
		var books []int64
		if len(mcbs2) > len(mcbs) {
			holmes.Error("len(dst-mcb)[%d] > len(src-mcb)[%d]", len(mcbs2), len(mcbs))
			return
		} else if len(mcbs2) < len(mcbs) {
			// check
			var j = 0
			var idx int64
			for ; j < len(mcbs2); j++ {
				if mcbs2[j].IndexId != mcbs[j].IndexId && mcbs2[j].BookId != mcbs[j].BookId {
					holmes.Error("check month course book error, maybe not same course.")
					return
				}
				idx = mcbs[j].IndexId
				books = append(books, mcbs[j].BookId)
			}
			// add month course book
			for ; j < len(mcbs); j++ {
				idx++
				mcb := models.MonthCourseBook{
					CourseId:      ucl.cfg.CopyCourse.DstCourseId,
					MonthCourseId: mcs2[i].ID,
					BookId:        mcbs[j].BookId,
					IndexId:       idx,
				}
				if err = models.CreateMonthCourseBook(&mcb); err != nil {
					holmes.Error("create month course book error: %v", err)
					return
				}
				books = append(books, mcbs[j].BookId)
			}
		} else {
			for j := 0; j < len(mcbs2); j++ {
				if mcbs2[j].IndexId != mcbs[j].IndexId && mcbs2[j].BookId != mcbs[j].BookId {
					holmes.Error("check month course book error, maybe not same course.")
					return
				}
				books = append(books, mcbs[j].BookId)
			}
		}
		holmes.Debug("books: %v", books)

		for j := 0; j < len(books); j++ {
			mccs, err := models.GetMonthCourseCatalogFromCourse(ucl.cfg.CopyCourse.SrcCourseId, mcs[i].ID, books[j])
			if err != nil {
				holmes.Error("get month course catalog from src error: %v", err)
				return
			}
			mccs2, err := models.GetMonthCourseCatalogFromCourse(ucl.cfg.CopyCourse.DstCourseId, mcs2[i].ID, books[j])
			if err != nil {
				holmes.Error("get month course catalog from dst error: %v", err)
				return
			}
			if len(mccs2) > len(mccs) {
				holmes.Error("len(dst-mcc)[%d] > len(src-mcc)[%d]", len(mccs2), len(mccs))
				return
			} else if len(mccs2) <= len(mccs) {
				// check
				var k = 0
				var idx int64
				var taskTime int64
				for ; k < len(mccs2); k++ {
					if mccs2[k].IndexId != mccs[k].IndexId && mccs2[k].Title != mccs[k].Title {
						holmes.Error("check month course catalog error, maybe not same course.")
						return
					}
					idx = mccs2[k].IndexId
					taskTime = mccs2[k].TaskTime
				}
				if taskTime != 0 {
					dstCourseTaskTimeOfWeekBegin = now.New(time.Unix(taskTime, 0).UTC()).BeginningOfWeek().Unix()
					//holmes.Debug("tasktime: %d %d", taskTime, dstCourseTaskTimeOfWeekBegin)
					taskWeekDay := (taskTime-dstCourseTaskTimeOfWeekBegin)/86400 + 1
					var ifFindWeekIdx bool
					for ii := 0; ii < len(ucl.cfg.CopyCourse.TaskWeekDate); ii++ {
						if taskWeekDay == ucl.cfg.CopyCourse.TaskWeekDate[ii] {
							taskWeekIdx = ii
							ifFindWeekIdx = true
							break
						}
					}
					if !ifFindWeekIdx {
						holmes.Error("task week day[%d] cannot found in config: %v", taskWeekDay, ucl.cfg.CopyCourse.TaskWeekDate)
						return
					}
				}
				holmes.Debug("taskTime: %d dstCourseTaskTimeOfWeekBegin: %d", taskTime, dstCourseTaskTimeOfWeekBegin)
				for ; k < len(mccs); k++ {
					if taskWeekIdx == (len(ucl.cfg.CopyCourse.TaskWeekDate) - 1) {
						taskWeekIdx = 0
						dstCourseTaskTimeOfWeekBegin = dstCourseTaskTimeOfWeekBegin + 7*86400
					} else {
						taskWeekIdx++
					}
					idx++
					newTaskTime := dstCourseTaskTimeOfWeekBegin + (ucl.cfg.CopyCourse.TaskWeekDate[taskWeekIdx]-1)*86400
					mcc := models.MonthCourseCatalog{
						CourseId:      ucl.cfg.CopyCourse.DstCourseId,
						MonthCourseId: mcs2[i].ID,
						BookId:        books[j],
						IndexId:       idx,
						Title:         mccs[k].Title,
						TaskTime:      newTaskTime,
					}
					if err = models.CreateMonthCourseCatalog(&mcc); err != nil {
						holmes.Error("create month course catalog error: %v", err)
						return
					}
					mccs2 = append(mccs2, mcc)
				}
			}

			// get month course catalog chapter
			for k := 0; k < len(mccs); k++ {
				chapters2, err := models.GetMonthCourseCatalogChapterFromCatalog(mccs2[k].ID)
				if err != nil {
					holmes.Error("get month course catalog chapter error: %v", err)
					return
				}
				if len(chapters2) != 0 {
					continue
				}
				chapters, err := models.GetMonthCourseCatalogChapterFromCatalog(mccs[k].ID)
				if err != nil {
					holmes.Error("get month course catalog chapter error: %v", err)
					return
				}
				// add month course catalog chapter
				for m := 0; m < len(chapters); m++ {
					mccc := &models.MonthCourseCatalogChapter{
						MonthCourseCatalogId: mccs2[k].ID,
						BookId:               books[j],
						ChapterId:            chapters[m].ChapterId,
						IndexId:              chapters[m].IndexId,
					}
					if err = models.CreateMonthCourseCatalogChapter(mccc); err != nil {
						holmes.Error("create month course catalog chapter error: %v", err)
						return
					}
				}

				audios2, err := models.GetMonthCourseCatalogAudioFromCatalog(mccs2[k].ID)
				if err != nil {
					holmes.Error("get month course catalog audio error: %v", err)
					return
				}
				if len(audios2) != 0 {
					continue
				}
				audios, err := models.GetMonthCourseCatalogAudioFromCatalog(mccs[k].ID)
				if err != nil {
					holmes.Error("get month course catalog audio error: %v", err)
					return
				}
				// add month course catalog audio
				for m := 0; m < len(audios); m++ {
					mcca := &models.MonthCourseCatalogAudio{
						MonthCourseCatalogId: mccs2[k].ID,
						AudioTitle:           audios[m].AudioTitle,
						AudioUrl:             audios[m].AudioUrl,
						AudioTime:            audios[m].AudioTime,
					}
					if err = models.CreateMonthCourseCatalogAudio(mcca); err != nil {
						holmes.Error("create month course catalog audio error: %v", err)
						return
					}
				}
			}
		}
	}
}
