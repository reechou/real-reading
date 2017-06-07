package ext

const (
	TULING_RESULT_SUCCESS = 10005
)

const (
	TULING_V1_RESULT_SUCCESS_TEXT = 100000
	TULING_V1_RESULT_SUCCESS_URL  = 200000
	TULING_V1_RESULT_SUCCESS_NEWS = 302000
	TULING_V1_RESULT_SUCCESS_CAI  = 308000
)

const (
	TULING_RESULT_TYPE_TEXT = "text"
	TULING_RESULT_TYPE_URL  = "url"
	TULING_RESULT_TYPE_NEWS = "news"
)

type TulingInputText struct {
	Text string `json:"text"`
}

type TulingLocation struct {
	City           string `json:"city,omitempty"`
	Latitude       string `json:"latitude,omitempty"`
	Longitude      string `json:"longitude,omitempty"`
	NearestPoiName string `json:"nearest_poi_name,omitempty"`
	Province       string `json:"province,omitempty"`
	Street         string `json:"street,omitempty"`
}

type TulingSelfInfo struct {
	Location TulingLocation `json:"location"`
}

type TulingPerception struct {
	InputText TulingInputText `json:"inputText"`
	SelfInfo  TulingSelfInfo  `json:"selfInfo,omitempty"`
}

type TulingUserInfo struct {
	ApiKey string `json:"apiKey"`
	UserId int    `json:"userId"`
}

type TulingRequest struct {
	Perception TulingPerception `json:"perception"`
	UserInfo   TulingUserInfo   `json:"userInfo"`
}

// about response
type TulingParameters struct {
	NearbyPlace string `json:"nearby_place"`
}

type TulingResultInfo struct {
	ResultType string `json:"resultType"`
	Url        string `json:"url"`
	Text       string `json:"text"`
	News       string `json:"news"`
}

type TulingIntent struct {
	Code       int              `json:"code"`
	IntentName string           `json:"intentName"`
	ActionName string           `json:"actionName"`
	Parameters TulingParameters `json:"parameters"`
}

type TulingResponse struct {
	Intent  TulingIntent       `json:"intent"`
	Results []TulingResultInfo `json:"results"`
}
