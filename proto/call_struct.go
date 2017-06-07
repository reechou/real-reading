package proto

const (
	RESPONSE_OK = iota
	RESPONSE_ERR
)

type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
