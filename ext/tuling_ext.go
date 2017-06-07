package ext

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/config"
)

type TulingExt struct {
	key    string
	client *http.Client
	cfg    *config.Config
}

func NewTulingExt(cfg *config.Config) *TulingExt {
	tuling := &TulingExt{
		client: &http.Client{},
		cfg:    cfg,
	}
	tuling.key = tuling.cfg.Tuling.Key

	return tuling
}

func (self *TulingExt) SimpleCall(msg string, userId int) (string, error) {
	request := &TulingRequest{
		Perception: TulingPerception{
			InputText: TulingInputText{
				Text: msg,
			},
		},
		UserInfo: TulingUserInfo{
			ApiKey: self.key,
			UserId: userId,
		},
	}
	reqBytes, err := json.Marshal(request)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return "", err
	}
	holmes.Debug("tuling request: %s", string(reqBytes))

	req, err := http.NewRequest("POST", self.cfg.Tuling.Url, bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return "", err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return "", err
	}
	var response TulingResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return "", err
	}
	if response.Intent.Code != TULING_RESULT_SUCCESS {
		holmes.Error("tuling simple call result code error: %d %v", response.Intent.Code, response)
		return "", fmt.Errorf("tuling simple call result error.")
	}
	result := ""
	for i := 0; i < len(response.Results); i++ {
		switch response.Results[i].ResultType {
		case TULING_RESULT_TYPE_TEXT:
			result += response.Results[i].Text
		case TULING_RESULT_TYPE_URL:
			result += response.Results[i].Url
		case TULING_RESULT_TYPE_NEWS:
			result += response.Results[i].News
		}
		if i != (len(response.Results) - 1) {
			result += "\n"
		}
	}

	return result, nil
}

func (self *TulingExt) SimpleCallV1(msg string, userId string) (string, error) {
	request := &TulingRequestV1{
		Key:    self.key,
		Info:   msg,
		UserId: userId,
	}
	reqBytes, err := json.Marshal(request)
	if err != nil {
		holmes.Error("json encode error: %v", err)
		return "", err
	}
	holmes.Debug("tuling request: %s", string(reqBytes))

	req, err := http.NewRequest("POST", self.cfg.Tuling.V1Url, bytes.NewBuffer(reqBytes))
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return "", err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return "", err
	}
	var response TulingResponseV1
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return "", err
	}

	result := ""
	switch response.Code {
	case TULING_V1_RESULT_SUCCESS_TEXT:
		result = response.Text
	case TULING_V1_RESULT_SUCCESS_URL:
		result = response.Text + "\n" + response.Url
	default:
		holmes.Error("tuling v1 simple call result code error: %d %v", response.Code, response)
		return "", fmt.Errorf("tuling v1 simple call result error.")
	}

	return result, nil
}
