package app

import (
	"github.com/duxphp/duxgo/v2/function"
	"github.com/duxphp/duxgo/v2/registry"
)

var DirList = []string{
	"./uploads",
	"./data",
	"./config",
	"./app",
	"./tmp",
	"./data/logs",
	"./data/logs/default",
	"./data/logs/request",
	"./data/logs/service",
	"./data/logs/database"}

func Init() {

	// 自动创建目录
	for _, path := range registry.DirList {
		if !function.IsExist(path) {
			if !function.CreateDir(path) {
				panic("failed to create " + path + " directory")
			}
		}
	}

	// 初始化
	for _, name := range Indexes {
		appConfig := List[name]
		if appConfig.Init != nil {
			appConfig.Init()
		}
	}

	// 注册
	for _, name := range Indexes {
		appConfig := List[name]
		if appConfig.Register != nil {
			appConfig.Register()
		}
	}

	// 启动
	for _, name := range Indexes {
		appConfig := List[name]
		if appConfig.Boot != nil {
			appConfig.Boot()
		}
	}

}
