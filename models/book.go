package models

import (
	"time"
	"fmt"
	
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
	Remark    string `xorm:"not null default '' varchar(128)" json:"remark"`
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

func DelBook(info *Book) error {
	if info.ID == 0 {
		return fmt.Errorf("del id cannot be nil.")
	}
	_, err := x.ID(info.ID).Delete(info)
	if err != nil {
		return err
	}
	return nil
}

func UpdateBook(info *Book) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("book_name", "author", "abstract", "cover", "updated_at").Update(info)
	return err
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

func GetBookList() ([]Book, error) {
	var books []Book
	err := x.Find(&books)
	if err != nil {
		return nil, err
	}
	return books, nil
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

func DelChapter(info *Chapter) error {
	if info.ID == 0 {
		return fmt.Errorf("del id cannot be nil.")
	}
	_, err := x.ID(info.ID).Delete(info)
	if err != nil {
		return err
	}
	return nil
}

func UpdateChapter(info *Chapter) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols("index_id", "title", "cover", "content", "remark", "updated_at").Update(info)
	return err
}

func UpdateChapterRemark(info *Chapter) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.ID(info.ID).Cols( "remark", "updated_at").Update(info)
	return err
}

func GetChapter(info *Chapter) (bool, error) {
	has, err := x.Id(info.ID).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		return false, nil
	}
	return true, nil
}

func GetChapterList(bookId int64) ([]Chapter, error) {
	var chapters []Chapter
	err := x.Where("book_id = ?", bookId).Find(&chapters)
	if err != nil {
		return nil, err
	}
	return chapters, nil
}
