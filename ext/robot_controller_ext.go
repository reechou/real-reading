package ext

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/config"
	"github.com/reechou/real-reading/robot_proto"
)

type RobotControllerExt struct {
	cfg    *config.Config
	client *http.Client
}

func NewRobotControllerExt(cfg *config.Config) *RobotControllerExt {
	return &RobotControllerExt{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (self *RobotControllerExt) LoginRobot(request *robot_proto.StartWxReq) (interface{}, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return nil, err
	}

	url := "http://" + self.cfg.RobotControllerHost.Host + ROBOT_CONTROLLER_URI_LOGIN_ROBOT
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return nil, err
	}
	var response RobotControllerResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return nil, err
	}
	if response.Code != ROBOT_CONTROLLER_RESPONSE_OK {
		holmes.Error("login robot[%v] result code error: %d %s", request, response.Code, response.Msg)
		return nil, fmt.Errorf("login robot error.")
	}

	return response.Data, nil
}

func (self *RobotControllerExt) GetRobot(request *GetRobotReq) (*RobotInfo, error) {
	reqBytes, err := json.Marshal(request)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return nil, err
	}

	url := "http://" + self.cfg.RobotControllerHost.Host + ROBOT_CONTROLLER_URI_GET_ROBOT
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return nil, err
	}
	var response RCGetRobotRsp
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return nil, err
	}
	if response.Code != ROBOT_CONTROLLER_RESPONSE_OK {
		holmes.Error("login robot[%v] result code error: %d %s", request, response.Code, response.Msg)
		return nil, fmt.Errorf("login robot error.")
	}

	return &response.Data, nil
}
