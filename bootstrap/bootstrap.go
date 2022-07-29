package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/duxphp/duxgo/core"
	"github.com/duxphp/duxgo/core/alarm"
	"github.com/duxphp/duxgo/core/cache"
	"github.com/duxphp/duxgo/core/config"
	"github.com/duxphp/duxgo/core/database"
	"github.com/duxphp/duxgo/core/exception"
	"github.com/duxphp/duxgo/core/logger"
	"github.com/duxphp/duxgo/core/register"
	"github.com/duxphp/duxgo/core/task"
	"github.com/duxphp/duxgo/core/util/function"
	"github.com/duxphp/duxgo/core/validator"
	"github.com/duxphp/duxgo/core/websocket"
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
	"strconv"
	"syscall"
	"time"
)

const Version = "0.1.0"

type bootstrap struct {
	app *echo.Echo
	Ch  chan os.Signal
}

// New 启动器
func New() *bootstrap {
	var ch = make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM)
	return &bootstrap{
		Ch: ch,
	}
}

// RegisterCore 注册服务
func (t *bootstrap) RegisterCore() *bootstrap {
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

	// 注册数据库
	database.GormInit()

	// 注册MangoDB
	database.QmgoInit()

	// 注册Redis
	database.RedisInit()

	// 注册队列服务
	task.Init()

	// 注册应用模块
	//t.RegisterApp()

	// 注册websocket服务
	websocket.InitSocket()

	return t
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (t *bootstrap) RegisterHttp() *bootstrap {
	// web服务
	t.app = echo.New()
	core.App = t.app

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
	render := &Template{
		templates: template.Must(template.New("").Delims("${", "}").Funcs(funcMap).ParseFS(core.ViewsFs, "views/*", "app/*/views/*")),
	}

	t.app.Renderer = render

	// 注册异常处理
	t.app.HTTPErrorHandler = func(err error, c echo.Context) {
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
		err = c.JSON(code, map[string]any{
			"code":    code,
			"message": msg,
		})
		if err != nil {
			core.Logger.Error().Err(err).Send()
		}
	}

	// 异常恢复处理
	t.app.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 4 << 10, // 1 KB
		LogLevel:  log.ERROR,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			core.Logger.Error().Err(err).Msg("PANIC RECOVER")
			return exception.Internal(err)
		},
	}))

	// IP 获取规则
	t.app.IPExtractor = func(req *http.Request) string {
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
	t.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
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
	t.app.Static("/uploads", "./uploads")
	t.app.StaticFS("/", echo.MustSubFS(core.StaticFs, "public"))

	// 前端中间件
	t.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{echo.HeaderContentType, echo.HeaderOrigin, echo.HeaderAccept, echo.HeaderXCSRFToken, echo.HeaderAuthorization, "X-dux-sfc", "x-dialog", "AccessKey", "X-Dux-Platform", "Content-MD5", "Content-Date"},
		ExposeHeaders: []string{"*"},
	}))

	// 关闭自带日志
	t.app.Logger.SetLevel(log.OFF)

	t.app.Use(middleware.RequestID())

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
		t.app.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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
func (t *bootstrap) RegisterApp() *bootstrap {
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
			appConfig.Register(t.app)
		}
	}

	// 应用路由
	for _, name := range register.AppIndex {
		appConfig := register.AppList[name]
		if appConfig.AppRoute != nil {
			appConfig.AppRoute(register.AppRouter, t.app)
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
			appConfig.AppRouteAuth(register.AppRouter, t.app)
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
			appConfig.Boot(t.app)
		}
	}

	return t
}

// StartTask 开启任务服务
func (t *bootstrap) StartTask() *bootstrap {
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
		core.Logger.Info().Msg("start scheduler service")
		task.StartScheduler()

	}()
	core.Logger.Info().Msg("start queue service")
	task.StartQueue()
	return t
}

func (t *bootstrap) StopTask() {
	core.Logger.Info().Msg("stop queue service")
	core.Queue.Shutdown()
	core.Logger.Info().Msg("stop scheduler service")
	core.Scheduler.Shutdown()
}

func (t *bootstrap) StartHttp() {
	// ping 队列服务
	task.Add("ping", &map[string]any{})

	prot := core.Config["app"].GetString("server.port")
	debug := core.Config["app"].GetBool("server.debug")

	data, _ := json.MarshalIndent(t.app.Routes(), "", "  ")
	ioutil.WriteFile("./routes.json", data, 0644)

	t.app.HideBanner = true

	var logo string

	const (
		cBlack   = "\u001b[90m"
		cRed     = "\u001b[91m"
		cCyan    = "\u001b[96m"
		cGreen   = "\u001b[92m"
		cYellow  = "\u001b[93m"
		cBlue    = "\u001b[94m"
		cMagenta = "\u001b[95m"
		cWhite   = "\u001b[97m"
		cReset   = "\u001b[0m"
	)

	value := func(s string, width int) string {
		pad := width - len(s)
		str := ""
		for i := 0; i < pad; i++ {
			str += "."
		}
		if s == "Disabled" {
			str += " " + s
		} else {
			str += fmt.Sprintf(" %s%s%s", cCyan, s, cBlack)
		}
		return str
	}

	centerValue := func(s string, width int) string {
		pad := strconv.Itoa((width - len(s)) / 2)
		str := fmt.Sprintf("%"+pad+"s", " ")
		str += fmt.Sprintf("%s%s%s", cCyan, s, cBlack)
		str += fmt.Sprintf("%"+pad+"s", " ")
		if len(str)-10 < width {
			str += " "
		}
		return str
	}

	logo += cBlack + " ┌───────────────────────────────────────────────────┐\n"
	logo += cBlack + " │ " + centerValue(" DuxGO v"+Version, 49) + " │\n"
	logo += cBlack + " │ " + centerValue("simple and fast development framework", 49) + " │\n"
	logo += cBlack + " │                                                   │\n"
	logo += fmt.Sprintf(cBlack+" │ Echo Ver %s  Routes %s │\n", value(echo.Version, 14), value(strconv.Itoa(len(t.app.Routes())), 15))
	logo += fmt.Sprintf(cBlack+" │ Debug .%s  PID ....%s │\n", value(lo.Ternary[string](debug, "enabled", "disabled"), 16), value(strconv.Itoa(os.Getpid()), 14))
	logo += cBlack + " └───────────────────────────────────────────────────┘\n" + cReset
	fmt.Print(logo)

	// 记录启动时间
	core.BootTime = time.Now()

	// 启动http服务
	go func() {
		serverAddr := ":" + prot
		err := t.app.Start(serverAddr)
		if err != nil {
			core.Logger.Error().Err(err).Msg("http stop")
		}
	}()
}

func (t *bootstrap) StopHttp() {
	core.Logger.Info().Msg("trigger a shutdown event")
	err, _ := event.Fire("app.close", event.M{})
	if err != nil {
		core.Logger.Error().Err(err).Msg("event stop")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	core.Logger.Info().Msg("stop the web service")

	if err := t.app.Shutdown(ctx); err != nil {
		core.Logger.Error().Err(err).Msg("http stop")
	}
}

func (t *bootstrap) Release() {
	core.Logger.Info().Msg("stop ants service")
	websocket.ReleaseSocket()
	ants.Release()
	os.Exit(0)

}
