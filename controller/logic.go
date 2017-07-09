package controller

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/config"
	"github.com/reechou/real-reading/ext"
	"github.com/reechou/real-reading/models"
)

type Logic struct {
	sync.Mutex

	robotExt *ext.RobotExt
	wc       *WechatController

	cfg *config.Config
}

func NewLogic(cfg *config.Config) *Logic {
	l := &Logic{
		cfg: cfg,
	}
	l.robotExt = ext.NewRobotExt(cfg)
	l.init()

	models.InitDB(cfg)

	return l
}

func (self *Logic) init() {
	http.HandleFunc("/robot/receive_msg", self.RobotReceiveMsg)
}

func (self *Logic) Run() {
	defer holmes.Start(holmes.LogFilePath("./log"),
		holmes.EveryDay,
		holmes.AlsoStdout,
		holmes.DebugLevel).Stop()

	if self.cfg.Debug {
		EnableDebug()
	}

	self.wc = NewWechatController(self.cfg)
	NewRemindWorker(self.cfg, self.wc)

	mux := http.NewServeMux()
	mux.Handle(ReadingPrefix+"/", NewReadingHandler(self))

	holmes.Info("server starting on[%s]..", self.cfg.Host)
	holmes.Infoln(http.ListenAndServe(self.cfg.Host, mux))
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func WriteBytes(w http.ResponseWriter, code int, v []byte) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	w.Write(v)
}

func EnableDebug() {

}
