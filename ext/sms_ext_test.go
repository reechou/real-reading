package ext

import (
	"github.com/reechou/real-reading/config"
	"testing"
	"fmt"
)

func TestNewSMSNotifyExt(t *testing.T) {
	sms := NewSMSNotifyExt(&config.Config{SMSNotify: config.SMSNotify{
		Host:       "http://v.juhe.cn/sms/send",
		TemplateId: "36745",
		Key:        "",
	}})
	err := sms.SMSNotify("", "周林栋")
	fmt.Println(err)
}
