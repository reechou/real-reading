package robot_proto

type RobotInfo struct {
	RobotWxNick string `json:"robot"`
	RunTime     int64  `json:"runTime"`
}

type StartWxReq struct {
	RobotType int `json:"robotType"`

	IfInvite           bool   `json:"ifInvite,omitempty"`
	IfInviteEndExit    bool   `json:"inviteEndExit,omitempty"`
	InviteMsg          string `json:"inviteMsg,omitempty"`
	IfClearWx          bool   `json:"ifClearWx,omitempty"`
	ClearWxMsg         string `json:"clearWxMsg,omitempty"`
	ClearWxPrefix      string `json:"clearWxPrefix,omitempty"`
	IfSaveRobotFriends bool   `json:"ifSaveRobotFriends,omitempty"`
	IfSaveRobotGroups  bool   `json:"ifSaveRobotGroups,omitempty"`
	// 是否替换emoji
	IfNotReplaceEmoji bool `json:"ifNotReplaceEmoji,omitempty"`
	// 建群逻辑
	IfCreateGroup     bool     `json:"ifCreateGroup,omitempty"`
	CreateGroupPrefix string   `json:"createGroupPrefix,omitempty"`
	CreateGroupStart  int      `json:"createGroupStart,omitempty"`
	CreateGroupNum    int      `json:"createGroupNum,omitempty"`
	CreateGroupUsers  []string `json:"createGroupUsers,omitempty"`
	// 不准修改群名
	IfNotChangeGroupName bool `json:"ifNotChangeGroupName,omitempty"`
	// 群加人逻辑
	IfSaveGroupMember         bool  `json:"ifSaveGroupMember,omitempty"`
	AddGroupMemberCycleOfTime int64 `json:"addGroupMemberCycleOfTime,omitempty"`
	AddGroupMemberCycleOfNum  int64 `json:"addGroupMemberCycleOfNum,omitempty"`
}

type RobotFindFriendReq struct {
	WechatNick string `json:"wechatNick"`
	UserName   string `json:"username"`
	NickName   string `json:"nickname"`
}

type RobotRemarkFriendReq struct {
	WechatNick string `json:"wechatNick"`
	UserName   string `json:"username"`
	NickName   string `json:"nickname"`
	Remark     string `json:"remark"`
}

type RobotGroupTirenReq struct {
	WechatNick     string `json:"wechatNick"`
	GroupUserName  string `json:"groupUserName"`
	GroupNickName  string `json:"groupNickName"`
	MemberUserName string `json:"memberUserName"`
	MemberNickName string `json:"memberNickName"`
}

type RobotGetGroupMemberListReq struct {
	WechatNick    string `json:"wechatNick"`
	GroupUserName string `json:"groupUserName"`
	GroupNickName string `json:"groupNickName"`
}

type RobotAddFriendReq struct {
	WechatNick    string `json:"wechatNick"`
	UserName      string `json:"userName"`
	VerifyContent string `json:"verifyContent"`
}

type RobotGetLoginsReq struct {
	RobotType int `json:"robotType"`
}

// about response
type WxFindFriendResponse struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data UserFriend `json:"data"`
}

type WxGroupTirenResponse struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data GroupUserInfo `json:"data"`
}

type WxGroupMemberListResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data []GroupUserInfo `json:"data"`
}

type WxResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
