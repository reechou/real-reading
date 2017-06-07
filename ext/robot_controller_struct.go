package ext

const (
	ROBOT_CONTROLLER_RESPONSE_OK = iota
	ROBOT_CONTROLLER_RESPONSE_ERR
)

const (
	ROBOT_CONTROLLER_URI_LOGIN_ROBOT                = "/manager/login_robot"
	ROBOT_CONTROLLER_URI_GET_ROBOT                  = "/manager/get_robot"
	ROBOT_CONTROLLER_URI_GET_ROBOTS_FROM_TYPE       = "/manager/get_robots_from_type"
	ROBOT_CONTROLLER_URI_GET_LOGIN_ROBOTS_FROM_TYPE = "/manager/get_login_robots_from_type"
)

type RobotInfo struct {
	ID        int64  `json:"id"`
	RobotType int    `json:"robotType"`
	Ip        string `json:"ip"`
	OfPort    string `json:"ofPort"`
}

type RCGetRobotRsp struct {
	Code int64     `json:"code"`
	Msg  string    `json:"msg,omitempty"`
	Data RobotInfo `json:"data,omitempty"`
}

type RobotControllerResponse struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type CreateRobotHost struct {
	Host     string `json:"host"`
	RobotNum int    `json:"robotNum"`
}

type GetRobotReq struct {
	RobotWx string `json:"robotWx"`
}

type GetRobotListFromTypeReq struct {
	RobotType int `json:"robotType"`
}

type GetLoginRobotListFromTypeReq struct {
	RobotType int `json:"robotType"`
}
