package web

import (
	"errors"
	"fmt"
	"github.com/duxphp/duxgo/v2/config"
	"github.com/duxphp/duxgo/v2/handlers"
	"github.com/duxphp/duxgo/v2/logger"
	"github.com/duxphp/duxgo/v2/registry"
	"github.com/duxphp/duxgo/v2/views"
	"github.com/duxphp/duxgo/v2/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/gookit/color"
	"github.com/gookit/event"
	"github.com/samber/lo"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

func Init() {
	// 注册 web 服务
	engine := views.Tpl()

	proxyHeader := config.Get("app").GetString("app.proxyHeader")
	registry.App = fiber.New(fiber.Config{
		AppName:               "DuxGO",
		Prefork:               false,
		CaseSensitive:         false,
		StrictRouting:         false,
		DisableStartupMessage: true,
		ProxyHeader:           lo.Ternary[string](proxyHeader != "", proxyHeader, "X-Real-IP"),
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var msg any
			if e, ok := err.(*handlers.CoreError); ok {
				// 程序错误
				msg = e.Message
			} else if e, ok := err.(*fiber.Error); ok {
				// http错误
				code = e.Code
				msg = e.Message
			} else {
				// 其他错误
				msg = err.Error()
				logger.Log().Error().Bytes("body", ctx.Body()).Err(err).Msg("error")
			}
			// 异步请求
			if ctx.Is("json") || ctx.XHR() {
				ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
				return ctx.Status(code).JSON(handlers.New(code, err.Error()))
			}

			// Web 请求
			if code == http.StatusNotFound {
				return ctx.Render("404.gohtml", fiber.Map{})
			} else {
				return ctx.Render("500.gohtml", fiber.Map{
					"code":    code,
					"message": msg,
				})
			}
		},
		Views: engine,
	})

	// 异常恢复处理
	registry.App.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			logger.Log().Error().Interface("err", e).Bytes("stack", debug.Stack()).Send()
		},
	}))

	// 超时处理
	t := config.Get("app").GetDuration("server.timeout")
	registry.App.Use(timeout.New(func(c *fiber.Ctx) error {
		if c.Get("upgrade") == "websocket" {
			return nil
		}
		return fiber.ErrRequestTimeout
	}, t*time.Second))

	// 注册静态路由
	registry.App.Static("/", "./public")

	// cors 跨域处理
	registry.App.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowHeaders:  "*",
		ExposeHeaders: "*",
	}))

	// 设置日志
	webLog := logger.New(
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
	registry.App.Use(fiberLogger.New(fiberLogger.Config{Output: webLog}))

	// 设置默认页面
	registry.App.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("welcome.gohtml", fiber.Map{})
	})

	// 注册请求ID
	registry.App.Use(requestid.New())

	// 注册websocket
	websocket.Init()
}

func Start() {
	port := config.Get("app").GetString("server.port")

	// 启动信息
	banner()

	// 记录启动时间
	registry.BootTime = time.Now()

	// 启动服务
	err := registry.App.Listen(":" + port)

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
	_ = registry.App.Shutdown()

	// 释放websocket服务
	websocket.Release()
}

func banner() {
	debugBool := config.Get("app").GetBool("server.debug")

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
		Name:  "Fiber",
		Value: fiber.Version,
	})
	sysMaps = append(sysMaps, item{
		Name:  "Debug",
		Value: lo.Ternary[string](debugBool, "enabled", "disabled"),
	})
	sysMaps = append(sysMaps, item{
		Name:  "PID",
		Value: os.Getpid(),
	})
	sysMaps = append(sysMaps, item{
		Name:  "Routes",
		Value: len(registry.App.Stack()),
	})

	banner += "⇨ "
	for _, v := range sysMaps {
		banner += v.Name + " <green>" + fmt.Sprintf("%v", v.Value) + "</>  "
	}
	banner += "\n"
	color.Print(banner)
}
