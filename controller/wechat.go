package controller

import (
	"github.com/chanxuehong/wechat.v2/mp/core"
	"github.com/chanxuehong/wechat.v2/mp/message/template"
	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/config"
)

type WechatController struct {
	cfg               *config.Config
	accessTokenServer core.AccessTokenServer
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
	holmes.Debug("template send msg success, msgid: %s", msgId)
	return nil
}
