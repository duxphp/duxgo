package global

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"time"
)

var (
	// App fiber应用
	App *fiber.App
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
