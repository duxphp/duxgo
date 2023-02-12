package register

import (
	util2 "github.com/duxphp/duxgo/v2/util"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

var (
	AppList   = make(map[string]*AppConfig)
	AppIndex  []string
	AppRouter = make(Router)
	AppMenu   = make(Menu)
)

type Router map[string]*util2.RouterData
type Menu map[string]*util2.MenuData

// AppConfig 注册规则
type AppConfig struct {
	Name         string                           //应用名称
	Config       any                              //应用配置
	Model        func()                           // 模型注册
	Register     func(*echo.Echo)                 //注册函数
	AppRoute     func(Router, *echo.Echo)         // 普通路由服务
	AppRouteAuth func(Router, *echo.Echo)         // 授权路由服务
	Route        func(Router)                     // 普通路由
	RouteAuth    func(Router)                     // 授权路由
	AppMenu      func(Menu)                       // 菜单注册
	Menu         func(Menu)                       // 菜单注册
	Event        func()                           // 事件注册
	Queue        func(queue *asynq.ServeMux)      // 队列注册
	Scheduler    func(scheduler *asynq.Scheduler) // 定时调度注册
	Websocket    func()                           // Socket服务注册
	Boot         func(*echo.Echo)                 //启动函数
}

// App 注册应用
func App(opt *AppConfig) {
	AppList[opt.Name] = opt
	AppIndex = append(AppIndex, opt.Name)
}
