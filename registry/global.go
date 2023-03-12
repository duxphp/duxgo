package registry

import (
	"context"
	"embed"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"html/template"
	"time"
)

var (
	// App echo应用
	App *echo.Echo
	// Version 版本号
	Version = "v2.0.0"
	// BootTime 启动时间
	BootTime time.Time
	// Di 依赖注入
	Di *do.Injector
	// TablePrefix 表前缀
	TablePrefix = "app_"
	// Debug 调试模式
	Debug bool
	// DebugMsg 屏蔽消息
	DebugMsg string
	// TimeLocation 时区
	TimeLocation *time.Location
	// Ctx Context
	Ctx = context.Background()
	// DirList 目录列表
	DirList []string
	// TplFs 系统视图
	TplFs embed.FS
	// Tpl 系统模板视图
	Tpl       *template.Template
	ConfigDir = "./config/"
	// Logger 日志服务
	//Logger zerolog.Logger
	// Validator 验证器
	Validator *validator.Validate
)
