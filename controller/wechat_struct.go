package controller

import (
	"github.com/chanxuehong/wechat.v2/mp/message/template"
)

const (
	TPL_ID_HOMEWORK_REMIND = "VmbspTYELpYT2wLYG6kBpkZU5DOjrfMqtnUxr1I60fY"
)

type TplMsg struct {
	ToUser string
	TplId  string
	Url    string
	Data   interface{}
}

type HomeworkRemindTplMsg struct {
	First    *template.DataItem `json:"first"`
	Keyword1 *template.DataItem `json:"keyword1"`
	Keyword2 *template.DataItem `json:"keyword2"`
	Keyword3 *template.DataItem `json:"keyword3"`
	Remark   *template.DataItem `json:"remark"`
}

type JssdkInfo struct {
	NonceStr  string
	Timestamp string
	Sign      string
	Url       string
}
