package models

import (
	"time"

	"github.com/reechou/holmes"
)

type Book struct {
	ID        int64  `xorm:"pk autoincr" json:"id"`
	BookName  string `xorm:"not null default '' varchar(256)" json:"bookName"`
	Author    string `xorm:"not null default '' varchar(256)" json:"author"`
	Abstract  string `xorm:"not null default '' varchar(512)" json:"abstract"`
	Cover     string `xorm:"not null default '' varchar(256)" json:"cover"`
	CreatedAt int64  `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt int64  `xorm:"not null default 0 int" json:"-"`
}

type Chapter struct {
	ID        int64  `xorm:"pk autoincr" json:"id"`
	BookId    int64  `xorm:"not null default 0 int index" json:"bookId"`
	IndexId   int64  `xorm:"not null default 0 int index" json:"indexId"`
	Title     string `xorm:"not null default '' varchar(256)" json:"title"`
	Cover     string `xorm:"not null default '' varchar(256)" json:"cover"`
	Content   string `xorm:"text" json:"content"`
	CreatedAt int64  `xorm:"not null default 0 int" json:"createdAt"`
	UpdatedAt int64  `xorm:"not null default 0 int" json:"-"`
}

func CreateBook(info *Book) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create book error: %v", err)
		return err
	}
	holmes.Info("create book[%v] success.", info)

	return nil
}

func CreateChapter(info *Chapter) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		holmes.Error("create Chapter error: %v", err)
		return err
	}
	holmes.Info("create Chapter[%d %s] success.", info.BookId, info.Title)

	return nil
}

func GetBook(info *Book) (bool, error) {
	has, err := x.Id(info.ID).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		holmes.Debug("cannot find book from id[%d]", info.ID)
		return false, nil
	}
	return true, nil
}
