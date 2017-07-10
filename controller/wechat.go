package controller

import (
	"fmt"
	"time"

	"github.com/chanxuehong/wechat.v2/mp/core"
	"github.com/chanxuehong/wechat.v2/mp/jssdk"
	"github.com/chanxuehong/wechat.v2/mp/message/template"
	"github.com/chanxuehong/wechat.v2/util"
	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/config"
)

type WechatController struct {
	cfg               *config.Config
	accessTokenServer core.AccessTokenServer
	ticketServer      jssdk.TicketServer
	wxClient          *core.Client
}

func NewWechatController(cfg *config.Config) *WechatController {
	wc := &WechatController{
		cfg: cfg,
	}
	wc.accessTokenServer = core.NewDefaultAccessTokenServer(
		cfg.ReadingOauth.ReadingWxAppId,
		cfg.ReadingOauth.ReadingWxAppSecret,
		nil)
	wc.wxClient = core.NewClient(wc.accessTokenServer, nil)
	wc.ticketServer = jssdk.NewDefaultTicketServer(wc.wxClient)

	return wc
}

func (self *WechatController) SendTplMsg(msg *TplMsg) error {
	tplMsg := &template.TemplateMessage2{
		ToUser:     msg.ToUser,
		TemplateId: msg.TplId,
		URL:        msg.Url,
		Data:       msg.Data,
	}
	msgId, err := template.Send(self.wxClient, tplMsg)
	if err != nil {
		holmes.Error("template send error: %v", err)
		return err
	}
	holmes.Debug("template send msg success, msgid: %d", msgId)
	return nil
}

func (self *WechatController) JssdkSign(info *JssdkInfo) {
	info.NonceStr = util.NonceStr()
	info.Timestamp = fmt.Sprintf("%d", time.Now().Unix())
	ticket, err := self.ticketServer.Ticket()
	if err != nil {
		holmes.Error("get jssdk ticket error: %v", err)
		return
	}
	info.Sign = jssdk.WXConfigSign(ticket,
		info.NonceStr,
		info.Timestamp,
		info.Url)
}
