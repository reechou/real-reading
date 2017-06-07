package ext

type TulingRequestV1 struct {
	Key    string `json:"key"`
	Info   string `json:"info"`
	UserId string `json:"userid"`
}

type TulingResponseV1 struct {
	Code int    `json:"code"`
	Text string `json:"text"`
	Url  string `json:"url"`
}
