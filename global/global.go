package duxgo

import (
	"context"
	"embed"
	"github.com/coocood/freecache"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/qiniu/qmgo"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"time"
)

var (
	// Version 版本号
	Version string
	// BootTime 启动时间
	BootTime time.Time
	// App echo
	App *echo.Echo
	// Debug 调试模式
	Debug bool
	// DebugMsg 屏蔽消息
	DebugMsg string
	// TimeLocation 时区
	TimeLocation *time.Location
	// Ctx Context
	Ctx      = context.Background()
	StaticFs embed.FS
	ViewsFs  embed.FS
	// Config 通用配置
	Config    = map[string]*viper.Viper{}
	ConfigDir = "./config"
	// ConfigManifest map[string]any

	// Logger 日志服务
	Logger zerolog.Logger
	// Alarm 消息通知
	Alarm any
	// Cache 公共缓存
	Cache *freecache.Cache
	// Db 数据库
	Db *gorm.DB
	// Mgo 数据库
	Mgo *qmgo.Database
	// Redis 数据
	Redis *redis.Client
	// Queue 队列服务端
	Queue *asynq.Server
	// QueueMux 队列调度复用
	QueueMux *asynq.ServeMux
	// QueueClient 队列客户端
	QueueClient *asynq.Client
	// QueueInspector 队列检查器
	QueueInspector *asynq.Inspector
	// Scheduler 调度服务端
	Scheduler *asynq.Scheduler
	// Validator 验证器
	Validator *validator.Validate
)
