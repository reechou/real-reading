package robot_proto

const (
	FROM_TYPE_PEOPLE = "people"
	FROM_TYPE_GROUP  = "group"

	RECEIVE_EVENT_MSG                  = "receivemsg"
	RECEIVE_EVENT_MOD_GROUP_ADD        = "modgroupadd"
	RECEIVE_EVENT_MOD_GROUP_ADD_DETAIL = "modgroupadddetail"
	RECEIVE_EVENT_ADD_FRIEND           = "addfriend"
	RECEIVE_EVENT_ADD                  = "receiveadd"
	RECEIVE_EVENT_ADD_GROUP            = "addgroup"
)

const (
	RECEIVE_MSG_TYPE_TEXT       = "text"
	RECEIVE_MSG_TYPE_IMG        = "img"
	RECEIVE_MSG_TYPE_VOICE      = "voice"
	RECEIVE_MSG_TYPE_VIDEO      = "video"
	RECEIVE_MSG_TYPE_CARD       = "card"
	RECEIVE_MSG_TYPE_SHARE      = "shareurl"
	RECEIVE_MSG_TYPE_TRANSFER   = "transfer"   // 转账
	RECEIVE_MSG_TYPE_RED_PACKET = "red-packet" // 红包
)

type WxGroup struct {
	NickName       string `json:"nickname"`
	UserName       string `json:"username"`
	GroupMemberNum int    `json:"groupMemberNum"`
}

type GroupUserInfo struct {
	DisplayName string `json:"displayName"`
	NickName    string `json:"nickname"`
	UserName    string `json:"username"`
}

type UserFriend struct {
	Alias       string `json:"alias"`
	City        string `json:"city"`
	VerifyFlag  int    `json:"verifyFlag"`
	ContactFlag int    `json:"contactFlag"`
	NickName    string `json:"nickName"`
	RemarkName  string `json:"remarkName"`
	Sex         int    `json:"sex"`
	UserName    string `json:"userName"`
}

type BaseInfo struct {
	Uin                string `json:"uin"`
	UserName           string `json:"userName,omitempty"`   // 机器人username
	WechatNick         string `json:"wechatNick,omitempty"` // 微信昵称
	ReceiveEvent       string `json:"receiveEvent,omitempty"`
	FromType           string `json:"fromType,omitempty"`
	FromUserName       string `json:"fromUserName,omitempty"`
	FromMemberUserName string `json:"fromMemberUserName,omitempty"`
	FromNickName       string `json:"fromNickName,omitempty"`
	FromGroupName      string `json:"fromGroupName,omitempty"`
}

type BaseToUserInfo struct {
	ToUserName  string `json:"toUserName,omitempty"`
	ToNickName  string `json:"toNickName,omitempty"`
	ToGroupName string `json:"toGroupName,omitempty"`
}

type SendBaseInfo struct {
	WechatNick string `json:"wechatNick,omitempty"` // 微信昵称
	ChatType   string `json:"chatType,omitempty"`
	NickName   string `json:"nickName,omitempty"`
	UserName   string `json:"userName,omitempty"`
	MsgType    string `json:"msgType,omitempty"`
	Msg        string `json:"msg,omitempty"`
}

type RetResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

type AddFriend struct {
	SourceWechat string `json:"sourceWechat,omitempty"`
	SourceNick   string `json:"sourceNick,omitempty"`
	UserWxid     string `json:"userWxid,omitempty"`
	UserWechat   string `json:"userWechat,omitempty"`
	UserNick     string `json:"userNick,omitempty"`
	UserCity     string `json:"userCity,omitempty"`
	UserSex      int    `json:"userSex,omitempty"`
	Ticket       string `json:"-"` // for verify
}

type ReceiveMsgInfo struct {
	BaseInfo       `json:"baseInfo,omitempty"`
	BaseToUserInfo `json:"baseToUserIno,omitempty"`
	AddFriend      `json:"addFriend,omitempty"`

	MsgType        string `json:"msgType,omitempty"`
	Msg            string `json:"msg,omitempty"`
	MediaTempUrl   string `json:"mediaTempUrl,omitempty"`
	GroupMemberNum int    `json:"groupMemberNum,omitempty"`
}

type CallbackMsgInfo struct {
	RetResponse `json:"retResponse,omitempty"`
	BaseInfo    `json:"baseInfo,omitempty"`

	CallbackMsgs []SendBaseInfo `json:"msg,omitempty"`
}

type SendMsgInfo struct {
	SendMsgs []SendBaseInfo `json:"sendBaseInfo,omitempty"`
}

type SendMsgResponse struct {
	RetResponse `json:"retResponse,omitempty"`
}
