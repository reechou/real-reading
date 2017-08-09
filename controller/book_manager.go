package controller

import (
	"encoding/json"
	"net/http"

	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/models"
	"github.com/reechou/real-reading/proto"
)

func (self *ReadingHandler) courseManagerGetBookList(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	bookList, err := models.GetBookList()
	if err != nil {
		holmes.Error("get book list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = bookList
}

func (self *ReadingHandler) courseManagerCreateBook(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.Book{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.CreateBook(req)
	if err != nil {
		holmes.Error("create book error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerDeleteBook(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.Book{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.DelBook(req)
	if err != nil {
		holmes.Error("delete book error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerUpdateBook(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.Book{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.UpdateBook(req)
	if err != nil {
		holmes.Error("update book error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerGetChapter(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &proto.ChapterReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	chapter := &models.Chapter{
		ID: req.ChapterId,
	}
	has, err := models.GetChapter(chapter)
	if err != nil {
		holmes.Error("get chapter list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	if !has {
		holmes.Error("this chapter[%d] cannot found", req.ChapterId)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = chapter
}

func (self *ReadingHandler) courseManagerGetChapterList(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &proto.BookReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	list, err := models.GetChapterList(req.BookId)
	if err != nil {
		holmes.Error("get chapter list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = list
}

func (self *ReadingHandler) courseManagerCreateChapter(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.Chapter{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.CreateChapter(req)
	if err != nil {
		holmes.Error("create chapter error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerDeleteChapter(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.Chapter{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.DelChapter(req)
	if err != nil {
		holmes.Error("delete chapter error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerUpdateChapter(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.Chapter{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.UpdateChapter(req)
	if err != nil {
		holmes.Error("update chapter error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerUpdateChapterRemark(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	req := &models.Chapter{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}

	err = models.UpdateChapterRemark(req)
	if err != nil {
		holmes.Error("update chapter remark error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerCreateBookCatalog(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()
	
	req := &models.BookCatalog{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	err = models.CreateBookCatalog(req)
	if err != nil {
		holmes.Error("create book catalog error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerDeleteBookCatalog(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()
	
	req := &models.BookCatalog{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	err = models.DelBookCatalog(req)
	if err != nil {
		holmes.Error("delete book catalog error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerUpdateBookCatalog(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()
	
	req := &models.BookCatalog{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	err = models.UpdateBookCatalog(req)
	if err != nil {
		holmes.Error("update book catalog error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerGetBookCatalogList(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()
	
	req := &proto.BookReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	list, err := models.GetBookCatalogList(req.BookId)
	if err != nil {
		holmes.Error("get book catalog list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = list
}

func (self *ReadingHandler) courseManagerCreateBookCatalogChapter(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()
	
	var req []models.BookCatalogChapter
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	err = models.CreateBookCatalogChapterList(req)
	if err != nil {
		holmes.Error("create book catalog chapter list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerDeleteBookCatalogChapter(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()
	
	req := &models.BookCatalogChapter{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	err = models.DelBookCatalogChapter(req)
	if err != nil {
		holmes.Error("delete book catalog chapter error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
}

func (self *ReadingHandler) courseManagerGetBookCatalogChapterList(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()
	
	req := &proto.BookCatalogReq{}
	err := json.Unmarshal(rr.Val, &req)
	if err != nil {
		holmes.Error("json unmarshal error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	
	list, err := models.GetBookCatalogChapterList(req.BookCatalogId)
	if err != nil {
		holmes.Error("get book catalog chapter list error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rsp.Data = list
}

