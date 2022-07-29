package alarm

import (
	"github.com/duxphp/duxgo/global"
	"github.com/go-resty/resty/v2"
	"time"
)

func Init() {
	global.Alarm = New(global.Config["app"].GetStringSlice("Alarm.urls"))
}

type Alarm struct {
	Url []string
}

// New 发送服务告警
func New(url []string) *Alarm {
	return &Alarm{url}
}

// Send 发送消息
func (t *Alarm) Send(title string, body string) {
	urlParams := "/" + title + "/" + body
	for _, s := range t.Url {
		resty.New().SetTimeout(10 * time.Second).R().Get(s + urlParams)
	}
}
