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

type ServerStatus struct {
	Database bool
	Redis    bool
	Mongodb  bool
}

func Init() {
	config.Init()
	logger.Init()
	cache.Init()
	validator.Init()
	views.Init()
	if Server.Database {
		database.GormInit()
	}
	if Server.Redis {
		database.RedisInit()
	}
	if Server.Mongodb {
		database.QmgoInit()
	}

}
