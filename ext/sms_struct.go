package ext

type SMSNotifyRsp struct {
	Code   int64  `json:"error_code"`
	Reason string `json:"reason"`
}
