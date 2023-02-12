package register

import (
	"github.com/duxphp/duxgo/bootstrap"
	"github.com/duxphp/duxgo/route"
	"github.com/duxphp/duxgo/util"
)

var (
	AppList  = make(map[string]*AppConfig)
	AppIndex []string
)

type Router map[string]*route.RouterData
type Menu map[string]*util.MenuData

// AppConfig 注册规则
type AppConfig struct {
	Name     string                       //应用名称
	Config   any                          //应用配置
	Init     func(t *bootstrap.Bootstrap) //初始化应用
	Register func(t *bootstrap.Bootstrap) // 注册应用
	Boot     func(t *bootstrap.Bootstrap) // 启动应用
}

// App 注册应用
func App(opt *AppConfig) {
	AppList[opt.Name] = opt
	AppIndex = append(AppIndex, opt.Name)
}
