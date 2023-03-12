package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duxphp/duxgo/v2/config"
	"github.com/duxphp/duxgo/v2/exception"
	"github.com/duxphp/duxgo/v2/function"
	"github.com/duxphp/duxgo/v2/logger"
	"github.com/duxphp/duxgo/v2/registry"
	"github.com/duxphp/duxgo/v2/views"
	"github.com/duxphp/duxgo/v2/websocket"
	"github.com/gookit/color"
	"github.com/gookit/event"
	"github.com/gookit/goutil/fsutil"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/samber/lo"
	"net"
	"net/http"
	"os"
	"time"
)

func Init() {
	// 注册 web 服务
	registry.App = echo.New()

	// 注册模板引擎
	render := &views.Template{
		Templates: views.Tpl(),
	}
	registry.App.Renderer = render

	// 注册异常处理
	registry.App.HTTPErrorHandler = func(err error, c echo.Context) {
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
			logger.Log().Error().Bytes("body", body).Err(err).Msg("error")
		}

		// AJAX请求
		if function.IsJson(c) {
			err = c.JSON(code, map[string]any{
				"code":    code,
				"message": msg,
			})
			if err != nil {
				logger.Log().Error().Err(err).Send()
			}
			return
		}
		// WEB请求
		if code == http.StatusNotFound {
			err = views.Tpl().ExecuteTemplate(c.Response(), "404.gohtml", nil)
		} else {
			err = views.Tpl().ExecuteTemplate(c.Response(), "500.gohtml", map[string]any{
				"code":    code,
				"message": msg,
			})
		}

		if err != nil {
			logger.Log().Error().Err(err).Send()
		}
	}

	// 异常恢复处理
	registry.App.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 4 << 10, // 1 KB
		LogLevel:  log.ERROR,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logger.Log().Error().Err(err).Bytes("stack", stack).Send()
			return exception.Internal(err)
		},
	}))

	// IP 获取规则
	registry.App.IPExtractor = func(req *http.Request) string {
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

	// 超时处理
	timeout := config.Get("app").GetInt("server.timeout")
	registry.App.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
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
	registry.App.Static("/uploads", "./uploads")

	// 注册虚拟目录
	//t.App.StaticFS("/", echo.MustSubFS(core.StaticFs, "public"))

	// cors 跨域处理
	registry.App.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
	}))

	// 关闭自带日志
	registry.App.Logger.SetLevel(log.OFF)

	// 设置默认页面
	registry.App.GET("/", func(c echo.Context) error {
		err := views.Tpl().ExecuteTemplate(c.Response(), "welcome.gohtml", nil)
		if err != nil {
			return err
		}
		return nil
	})

	// 注册请求ID
	registry.App.Use(middleware.RequestID())

	// 访问日志

	if config.Get("app").GetBool("logger.request.status") {
		vLog := logger.New(
			logger.GetWriter(
				config.Get("app").GetString("logger.request.level"),
				config.Get("app").GetString("logger.request.path")+"/web.log",
				config.Get("app").GetInt("logger.request.maxSize"),
				config.Get("app").GetInt("logger.request.maxBackups"),
				config.Get("app").GetInt("logger.request.maxAge"),
				config.Get("app").GetBool("logger.request.compress"),
				true,
			),
		).With().Timestamp().Logger()
		registry.App.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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

		// 注册websocket
		websocket.Init()
	}
}

func Start() {
	port := config.Get("app").GetString("server.port")
	data, _ := json.MarshalIndent(registry.App.Routes(), "", "  ")
	_ = fsutil.WriteFile("./routes.json", data, 0644)

	registry.App.HideBanner = true
	banner()

	// 记录启动时间
	registry.BootTime = time.Now()

	// 启动服务
	serverAddr := ":" + port
	err := registry.App.Start(serverAddr)

	// 退出服务
	if errors.Is(err, http.ErrServerClosed) {
		color.Print("\n⇨ <red>Server closed</>\n")
		return
	}
	if err != nil {
		logger.Log().Error().Err(err).Msg("web")
	}
	// 退出事件
	err, _ = event.Fire("app.close", event.M{})
	if err != nil {
		logger.Log().Error().Err(err).Msg("event stop")
	}

	// 关闭服务
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = registry.App.Shutdown(ctx)

	// 释放websocket服务
	websocket.Release()
}

func banner() {
	debug := config.Get("app").GetBool("server.debug")

	var banner string
	banner += `   _____           ____ ____` + "\n"
	banner += `  / __  \__ ______/ ___/ __ \` + "\n"
	banner += ` / /_/ / /_/ /> </ (_ / /_/ /` + "\n"
	banner += `/_____/\_,__/_/\_\___/\____/  v` + registry.Version + "\n"

	type item struct {
		Name  string
		Value any
	}

	var sysMaps []item
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
		Value: len(registry.App.Routes()),
	})

	banner += "⇨ "
	for _, v := range sysMaps {
		banner += v.Name + " <green>" + fmt.Sprintf("%v", v.Value) + "</>  "
	}
	banner += "\n"
	color.Print(banner)
}
