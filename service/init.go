package service

import (
	"github.com/duxphp/duxgo/v2/cache"
	"github.com/duxphp/duxgo/v2/config"
	"github.com/duxphp/duxgo/v2/database"
	"github.com/duxphp/duxgo/v2/logger"
	"github.com/duxphp/duxgo/v2/validator"
	"github.com/duxphp/duxgo/v2/views"
)

var Server = ServerStatus{}

// ServerStatus 服务状态
type ServerStatus struct {
	Database bool // 数据库服务
	Redis    bool // redis服务
	Mongodb  bool // mongodb服务
}

func Init() {
	// 配置服务
	config.Init()
	// 日志服务
	logger.Init()
	// 加载缓存器
	cache.Init()
	// 注册验证器
	validator.Init()
	// 注册模板引擎
	views.Init()
	// 注册数据库
	if Server.Database {
		database.GormInit()
	}
	// 注册redis
	if Server.Redis {
		database.RedisInit()
	}
	// 注册mongodb
	if Server.Mongodb {
		database.QmgoInit()
	}

}
