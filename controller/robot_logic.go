package controller

import (
	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/robot_proto"
)

func (self *Logic) HandleReceiveMsg(msg *robot_proto.ReceiveMsgInfo) {
	holmes.Debug("receive robot msg: %v", msg)
	switch msg.BaseInfo.ReceiveEvent {
	case robot_proto.RECEIVE_EVENT_MSG:
		self.handleMsg(msg)
	}
}

func (self *Logic) handleMsg(msg *robot_proto.ReceiveMsgInfo) {

}

type MsgInfo struct {
	Msg     string
	MsgType string
}

func (self *Logic) sendMsgBack(msgs []MsgInfo, msg *robot_proto.ReceiveMsgInfo) error {
	robotHost := self.cfg.DefaultRobotHost
	var sendReq robot_proto.SendMsgInfo
	for _, v := range msgs {
		sendReq.SendMsgs = append(sendReq.SendMsgs, robot_proto.SendBaseInfo{
			WechatNick: msg.BaseInfo.WechatNick,
			ChatType:   msg.BaseInfo.FromType,
			UserName:   msg.BaseInfo.FromUserName,
			NickName:   msg.BaseInfo.FromNickName,
			MsgType:    v.MsgType,
			Msg:        v.Msg,
		})
	}
	err := self.robotExt.SendMsgs(robotHost, &sendReq)
	if err != nil {
		holmes.Error("send msg[%v] back error: %v", sendReq, err)
		return err
	}
	return nil
}
