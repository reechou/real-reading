package controller

import (
	"sync"

	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/config"
)

const (
	DEFAULT_TPL_WORKER_NUM = 1024
)

type WechatTplWorker struct {
	wg  sync.WaitGroup
	cfg *config.Config
	wc  *WechatController

	msgChan   chan *TplMsg
	WorkerNum int

	stop chan struct{}
}

func NewWechatTplWorker(cfg *config.Config, wc *WechatController) *WechatTplWorker {
	wtw := &WechatTplWorker{
		cfg:     cfg,
		wc:      wc,
		msgChan: make(chan *TplMsg, 10240),
		stop:    make(chan struct{}),
	}
	if cfg.TplWorkerNum == 0 {
		wtw.WorkerNum = DEFAULT_TPL_WORKER_NUM
	} else {
		wtw.WorkerNum = cfg.TplWorkerNum
	}

	for i := 0; i < wtw.WorkerNum; i++ {
		wtw.wg.Add(1)
		go wtw.runWorker()
	}

	holmes.Debug("wechat tpl worker start..")

	return wtw
}

func (self *WechatTplWorker) Stop() {
	close(self.stop)
	self.wg.Wait()
}

func (self *WechatTplWorker) Send(msg *TplMsg) {
	select {
	case self.msgChan <- msg:
	case <-self.stop:
		return
	}
}

func (self *WechatTplWorker) runWorker() {
	for {
		select {
		case msg := <-self.msgChan:
			self.sendTplMsg(msg)
		case <-self.stop:
			self.wg.Done()
			return
		}
	}
}

func (self *WechatTplWorker) sendTplMsg(msg *TplMsg) {
	err := self.wc.SendTplMsg(msg)
	if err != nil {
		holmes.Error("send tpl msg error: %v", err)
	}
}
