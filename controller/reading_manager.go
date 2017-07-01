package controller

import (
	"net/http"

	"github.com/jinzhu/now"
	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/models"
	"github.com/reechou/real-reading/proto"
)

func (self *ReadingHandler) readingPayToday(rr *HandlerRequest, w http.ResponseWriter, r *http.Request) {
	rsp := &proto.Response{Code: proto.RESPONSE_OK}
	defer func() {
		writeRsp(w, rsp)
	}()

	//now := time.Now().Unix()
	//todayZero := now - (now % 86400) - 28800
	todayZero := now.BeginningOfDay().Unix()
	list, err := models.GetUserCourseFromTime(todayZero)
	if err != nil {
		holmes.Error("get reading pay from time error: %v", err)
		rsp.Code = proto.RESPONSE_ERR
		return
	}
	rpToday := &proto.ReadingPayToday{
		OrderNum: int64(len(list)),
	}
	for _, v := range list {
		rpToday.AllMoney = rpToday.AllMoney + (v.Money / 100)
	}
	rsp.Data = rpToday
}
