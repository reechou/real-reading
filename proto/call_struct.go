package proto

const (
	RESPONSE_OK = iota
	RESPONSE_ERR
	RESPONSE_ERR_EXT
	RESPONSE_REFUND_NOT_AUTO
	RESPONSE_REFUND_HAS_EXECED
	RESPONSE_ERR_COUPON_NOT_FOUND
	RESPONSE_COUPON_HAS_USED
	RESPONSE_COUPON_HAS_LOCKED
)

type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
