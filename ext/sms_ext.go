package ext

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/reechou/holmes"
	"github.com/reechou/real-reading/config"
)

type SMSNotifyExt struct {
	client *http.Client
	cfg    *config.Config
}

func NewSMSNotifyExt(cfg *config.Config) *SMSNotifyExt {
	sms := &SMSNotifyExt{
		client: &http.Client{},
		cfg:    cfg,
	}

	return sms
}

func (self *SMSNotifyExt) SMSNotifyNormal(mobile, tmpId, params string) error {
	if mobile == "" {
		return fmt.Errorf("not found mobile")
	}

	requestUrl := self.cfg.SMSNotify.Host
	parseRequestUrl, _ := url.Parse(requestUrl)

	tplValue := url.QueryEscape(params)
	extraParams := url.Values{
		"mobile":    {mobile},
		"tpl_id":    {tmpId},
		"tpl_value": {tplValue},
		"key":       {self.cfg.SMSNotify.Key},
	}
	parseRequestUrl.RawQuery = extraParams.Encode()

	req, err := http.NewRequest("GET", parseRequestUrl.String(), nil)
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return err
	}
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return err
	}
	var response SMSNotifyRsp
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return err
	}
	if response.Code != 0 {
		holmes.Error("sms notify error: %s", response.Reason)
		return fmt.Errorf("sms notify error: %s", response.Reason)
	}

	holmes.Info("sms send notify to mobile[%s] success", mobile)

	return nil
}

func (self *SMSNotifyExt) SMSNotify(mobile, realName string) error {
	if mobile == "" {
		return fmt.Errorf("not found mobile")
	}

	requestUrl := self.cfg.SMSNotify.Host
	parseRequestUrl, _ := url.Parse(requestUrl)

	tplValue := url.QueryEscape(fmt.Sprintf("#RealName#=%s", realName))
	extraParams := url.Values{
		"mobile":    {mobile},
		"tpl_id":    {self.cfg.SMSNotify.TemplateId},
		"tpl_value": {tplValue},
		"key":       {self.cfg.SMSNotify.Key},
	}
	parseRequestUrl.RawQuery = extraParams.Encode()

	req, err := http.NewRequest("GET", parseRequestUrl.String(), nil)
	if err != nil {
		holmes.Error("http new request error: %v", err)
		return err
	}
	//req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := self.client.Do(req)
	if err != nil {
		holmes.Error("http do request error: %v", err)
		return err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		holmes.Error("ioutil ReadAll error: %v", err)
		return err
	}
	var response SMSNotifyRsp
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		holmes.Error("json decode error: %v [%s]", err, string(rspBody))
		return err
	}
	if response.Code != 0 {
		holmes.Error("sms notify error: %s", response.Reason)
		return fmt.Errorf("sms notify error: %s", response.Reason)
	}

	holmes.Info("sms send notify to mobile[%s] realname[%s] success", mobile, realName)

	return nil
}
