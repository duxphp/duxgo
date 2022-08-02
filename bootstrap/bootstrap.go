package bootstrap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duxphp/duxgo/alarm"
	"github.com/duxphp/duxgo/cache"
	"github.com/duxphp/duxgo/config"
	"github.com/duxphp/duxgo/core"
	"github.com/duxphp/duxgo/exception"
	"github.com/duxphp/duxgo/logger"
	"github.com/duxphp/duxgo/register"
	"github.com/duxphp/duxgo/task"
	"github.com/duxphp/duxgo/util/function"
	"github.com/duxphp/duxgo/validator"
	"github.com/gookit/color"
	"github.com/gookit/event"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/panjf2000/ants/v2"
	"github.com/samber/lo"
	"html/template"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const Version = "0.1.0"

type Bootstrap struct {
	App *echo.Echo
	Ch  chan os.Signal
}

// New 启动器
func New() *Bootstrap {
	var ch = make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM)
	return &Bootstrap{
		Ch: ch,
	}
}

// RegisterCore 注册服务
func (t *Bootstrap) RegisterCore() *Bootstrap {
	// 设置时区
	core.TimeLocation = time.FixedZone("CST", 8*3600)
	core.Version = Version
	time.Local = core.TimeLocation

	// 配置服务
	config.Init()

	// 日志服务
	logger.Init()

	// 告警服务
	alarm.Init()

	// 加载缓存器
	cache.Init()

	// 注册验证器
	validator.Init()

	// 注册模板引擎
	funcMap := template.FuncMap{
		"unescape": func(s string) template.HTML {
			return template.HTML(s)
		},
		"marshal": func(v interface{}) template.JS {
			a, _ := json.Marshal(v)
			return template.JS(a)
		},
	}
	tpl := template.Must(template.New("").Delims("${", "}").Funcs(funcMap).ParseFS(core.TplFs, "template/*"))
	core.Tpl = tpl

	// 注册目录
	core.DirList = []string{"./uploads"}

	return t
}

// Template 模板服务
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// RegisterHttp 注册http服务
func (t *Bootstrap) RegisterHttp() *Bootstrap {
	// web服务
	t.App = echo.New()
	core.App = t.App

	// 注册模板
	render := &Template{
		templates: core.Tpl,
	}
	t.App.Renderer = render

	// 注册异常处理
	t.App.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		var msg any
		if e, ok := err.(*exception.CoreError); ok {
			msg = e.Message
		} else if e, ok := err.(*echo.HTTPError); ok {
			code = e.Code
			msg = e.Message
		} else {
			msg = err.Error()
			body := function.CtxBody(c)
			core.Logger.Error().Bytes("body", body).Err(err).Msg("error")
		}

		// AJAX请求
		if function.IsAjax(c) {
			err = c.JSON(code, map[string]any{
				"code":    code,
				"message": msg,
			})
			if err != nil {
				core.Logger.Error().Err(err).Send()
			}
			return
		}
		// WEB请求
		if code == http.StatusNotFound {
			err = core.Tpl.ExecuteTemplate(c.Response(), "404.gohtml", nil)
		} else {
			err = core.Tpl.ExecuteTemplate(c.Response(), "500.gohtml", map[string]any{
				"code":    code,
				"message": msg,
			})
		}
		if err != nil {
			core.Logger.Error().Err(err).Send()
		}

	}

	// 异常恢复处理
	t.App.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 4 << 10, // 1 KB
		LogLevel:  log.ERROR,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			return exception.Internal(err)
		},
	}))

	// IP 获取规则
	t.App.IPExtractor = func(req *http.Request) string {
		remoteAddr := req.RemoteAddr
		if ip := req.Header.Get(echo.HeaderXRealIP); ip != "" {
			remoteAddr = ip
		} else if ip = req.Header.Get(echo.HeaderXForwardedFor); ip != "" {
			remoteAddr = ip
		} else {
			remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
		}
		if remoteAddr == "::1" {
			remoteAddr = "127.0.0.1"
		}
		return remoteAddr
	}

	// 链接超时
	timeout := core.Config["app"].GetInt("server.timeout")
	t.App.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper: func(c echo.Context) bool {
			if c.IsWebSocket() {
				return true
			} else {
				return false
			}
		},
		Timeout: time.Duration(timeout) * time.Second,
	}))

	// 注册静态路由
	t.App.Static("/uploads", "./uploads")
	//t.App.StaticFS("/", echo.MustSubFS(core.StaticFs, "public"))

	// 前端中间件
	t.App.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{echo.HeaderContentType, echo.HeaderOrigin, echo.HeaderAccept, echo.HeaderXCSRFToken, echo.HeaderAuthorization, "X-dux-sfc", "x-dialog", "AccessKey", "X-Dux-Platform", "Content-MD5", "Content-Date"},
		ExposeHeaders: []string{"*"},
	}))

	// 关闭自带日志
	t.App.Logger.SetLevel(log.OFF)

	// 设置默认页面
	t.App.GET("/", func(c echo.Context) error {
		err := core.Tpl.ExecuteTemplate(c.Response(), "welcome.gohtml", nil)
		if err != nil {
			return err
		}
		return nil
	})

	// 注册请求ID
	t.App.Use(middleware.RequestID())

	// 访问日志
	if core.Config["app"].GetBool("logger.request.status") {
		vLog := logger.New(
			core.Config["app"].GetString("logger.request.level"),
			core.Config["app"].GetString("logger.request.path"),
			core.Config["app"].GetInt("logger.request.maxSize"),
			core.Config["app"].GetInt("logger.request.maxBackups"),
			core.Config["app"].GetInt("logger.request.maxAge"),
			core.Config["app"].GetBool("logger.request.compress"),
		).With().Timestamp().Logger()
		t.App.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:       true,
			LogHost:      true,
			LogStatus:    true,
			LogMethod:    true,
			LogLatency:   true,
			LogRemoteIP:  true,
			LogError:     true,
			LogRequestID: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				logT := vLog.Info()
				if v.Latency > 1*time.Second {
					logT = vLog.Warn()
				}
				logT.Str("id", v.RequestID).
					Int("status", v.Status).
					Str("method", v.Method).
					Str("uri", v.URI).
					Str("ip", v.RemoteIP).
					Dur("latency", v.Latency).
					Err(v.Error).
					Msg("request")
				return nil
			},
		}))
	}

	return t
}

// RegisterApp 注册应用
func (t *Bootstrap) RegisterApp() *Bootstrap {
	// 注册模型
	for _, name := range register.AppIndex {
		appConfig := register.AppList[name]
		if appConfig.Model != nil {
			appConfig.Model()
		}
	}

	// 注册服务
	for _, name := range register.AppIndex {
		appConfig := register.AppList[name]
		if appConfig.Register != nil {
			appConfig.Register(t.App)
		}
	}

	// 应用路由
	for _, name := range register.AppIndex {
		appConfig := register.AppList[name]
		if appConfig.AppRoute != nil {
			appConfig.AppRoute(register.AppRouter, t.App)
		}
	}

	// 注册路由
	for _, name := range register.AppIndex {
		appConfig := register.AppList[name]
		if appConfig.Route != nil {
			appConfig.Route(register.AppRouter)
		}
	}

	// 应用授权路由
	for _, name := range register.AppIndex {
		appConfig := register.AppList[name]
		if appConfig.AppRouteAuth != nil {
			appConfig.AppRouteAuth(register.AppRouter, t.App)
		}
	}

	// 注册授权路由
	for _, name := range register.AppIndex {
		appConfig := register.AppList[name]
		if appConfig.RouteAuth != nil {
			appConfig.RouteAuth(register.AppRouter)
		}
	}

	// 应用菜单
	for _, name := range register.AppIndex {
		appConfig := register.AppList[name]
		if appConfig.AppMenu != nil {
			appConfig.AppMenu(register.AppMenu)
		}
	}

	// 菜单注册
	for _, name := range register.AppIndex {
		appConfig := register.AppList[name]
		if appConfig.Menu != nil {
			appConfig.Menu(register.AppMenu)
		}
	}

	// 事件注册
	for _, name := range register.AppIndex {
		appConfig := register.AppList[name]
		if appConfig.Event != nil {
			appConfig.Event()
		}
	}

	// Socket注册
	for _, name := range register.AppIndex {
		appConfig := register.AppList[name]
		if appConfig.Websocket != nil {
			appConfig.Websocket()
		}
	}

	// 启动服务
	for _, name := range register.AppIndex {
		appConfig := register.AppList[name]
		if appConfig.Boot != nil {
			appConfig.Boot(t.App)
		}
	}

	return t
}

// StartTask 开启任务服务
func (t *Bootstrap) StartTask() *Bootstrap {
	// 队列注册
	for _, name := range register.AppIndex {
		appConfig := register.AppList[name]
		if appConfig.Queue != nil {
			appConfig.Queue(core.QueueMux)
		}
	}

	// 调度注册
	for _, name := range register.AppIndex {
		appConfig := register.AppList[name]
		if appConfig.Scheduler != nil {
			appConfig.Scheduler(core.Scheduler)
		}
	}
	//启动队列与调度服务
	go func() {
		task.StartScheduler()

	}()
	task.StartQueue()
	return t
}

// StopTask 停止任务服务
func (t *Bootstrap) StopTask() {
	core.Queue.Shutdown()
	core.Scheduler.Shutdown()
}

// StartHttp 启动http服务

func (t *Bootstrap) StartHttp() {

	// 自动创建目录
	for _, path := range core.DirList {
		if !function.IsExist(path) {
			if !function.CreateDir(path) {
				panic("failed to create " + path + " directory")
			}
		}
	}

	prot := core.Config["app"].GetString("server.port")
	debug := core.Config["app"].GetBool("server.debug")

	data, _ := json.MarshalIndent(t.App.Routes(), "", "  ")
	ioutil.WriteFile("./routes.json", data, 0644)

	t.App.HideBanner = true

	var banner string
	banner += `   _____           ____ ____` + "\n"
	banner += `  / __  \__ ______/ ___/ __ \` + "\n"
	banner += ` / /_/ / /_/ /> </ (_ / /_/ /` + "\n"
	banner += `/_____/\_,__/_/\_\___/\____/  v` + Version + "\n"

	type item struct {
		Name  string
		Value any
	}

	sysMaps := []item{}
	sysMaps = append(sysMaps, item{
		Name:  "Echo",
		Value: echo.Version,
	})
	sysMaps = append(sysMaps, item{
		Name:  "Debug",
		Value: lo.Ternary[string](debug, "enabled", "disabled"),
	})
	sysMaps = append(sysMaps, item{
		Name:  "PID",
		Value: os.Getpid(),
	})
	sysMaps = append(sysMaps, item{
		Name:  "Routes",
		Value: len(t.App.Routes()),
	})

	banner += "⇨ "
	for _, v := range sysMaps {
		banner += v.Name + " <green>" + fmt.Sprintf("%v", v.Value) + "</>  "
	}
	banner += "\n"
	color.Print(banner)

	// 记录启动时间
	core.BootTime = time.Now()

	// 启动http服务
	go func() {
		serverAddr := ":" + prot
		err := t.App.Start(serverAddr)
		if errors.Is(err, http.ErrServerClosed) {
			color.Print("\n⇨ <red>Server closed</>\n")
			return
		}
		if err != nil {
			core.Logger.Error().Err(err).Msg("web")
		}
	}()
}

// StopHttp 停止http服务
func (t *Bootstrap) StopHttp() {
	err, _ := event.Fire("App.close", event.M{})
	if err != nil {
		core.Logger.Error().Err(err).Msg("event stop")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := t.App.Shutdown(ctx); err != nil {
		core.Logger.Error().Err(err).Msg("web")
	}
}

// Release 释放服务
func (t *Bootstrap) Release() {
	ants.Release()
	os.Exit(0)

}
