package controller

import (
	"sync"

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

func NewWechatTplWorker(cfg *config.Config) *WechatTplWorker {
	wtw := &WechatTplWorker{
		cfg:     cfg,
		msgChan: make(chan *TplMsg, 10240),
	}
	if cfg.TplWorkerNum == 0 {
		wtw.WorkerNum = DEFAULT_TPL_WORKER_NUM
	} else {
		wtw.WorkerNum = cfg.TplWorkerNum
	}
	wtw.wc = NewWechatController(cfg)

	return wtw
}

func (self *WechatTplWorker) runWorker() {

}
