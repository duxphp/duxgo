package registry

import (
	"context"
	"github.com/labstack/echo/v4"
	"time"
)

var (
	// App echo应用
	App *echo.Echo
	// Version 版本号
	Version = "v2.0.0"
	// BootTime 启动时间
	BootTime time.Time
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
	DirList   []string
	ConfigDir = "./config/"
)
