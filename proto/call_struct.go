package proto

const (
	RESPONSE_OK = iota
	RESPONSE_ERR
	RESPONSE_ERR_EXT
	RESPONSE_REFUND_NOT_AUTO
)

type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
