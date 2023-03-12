package duxgo

import (
	"embed"
	"github.com/duxphp/duxgo/v2/app"
	"github.com/duxphp/duxgo/v2/database"
	"github.com/duxphp/duxgo/v2/registry"
	"github.com/duxphp/duxgo/v2/service"
	"github.com/duxphp/duxgo/v2/task"
	"github.com/duxphp/duxgo/v2/views"
	"github.com/duxphp/duxgo/v2/web"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/cobra"
	"os"
	"time"
)

type Dux struct {
	registerApp []func()
	registerCmd []func(command *cobra.Command)
}

func New() *Dux {
	return &Dux{}
}

// RegisterApp 注册应用
func (t *Dux) RegisterApp(calls ...func()) {
	t.registerApp = append(t.registerApp, calls...)
}

// RegisterCmd 注册命令
func (t *Dux) RegisterCmd(calls ...func(command *cobra.Command)) {
	t.registerCmd = append(t.registerCmd, calls...)
}

// RegisterDir 注册目录
func (t *Dux) RegisterDir(dirs ...string) {
	app.DirList = append(app.DirList, dirs...)
}

//go:embed template/*
var tplFs embed.FS

// 创建通用服务
func (t *Dux) create() {

	// 设置时区
	registry.TimeLocation = time.FixedZone("CST", 8*3600)
	time.Local = registry.TimeLocation

	// 设置模板
	views.TplFs = tplFs

	// 注册应用
	for _, call := range t.registerApp {
		call()
	}

	// 注册命令
	t.RegisterCmd(app.Command, web.Command, task.Command, database.Command)
}

// Run 运行命令
func (t *Dux) Run() {
	// 构架功能
	t.create()

	// 注册命令
	var rootCmd = &cobra.Command{Use: "dux"}
	for _, cmd := range t.registerCmd {
		cmd(rootCmd)
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// SetTablePrefix 设置数据表前缀
func (t *Dux) SetTablePrefix(prefix string) {
	registry.TablePrefix = prefix
}

// SetConfigDir 设置配置目录
func (t *Dux) SetConfigDir(dir string) {
	registry.ConfigDir = dir
}

// SetDatabaseStatus 设置数据库状态
func (t *Dux) SetDatabaseStatus(status bool) {
	service.Server.Database = status
}

// SetRedisStatus 设置redis状态
func (t *Dux) SetRedisStatus(status bool) {
	service.Server.Redis = status
}

// SetMongodbStatus 设置mongodb状态
func (t *Dux) SetMongodbStatus(status bool) {
	service.Server.Mongodb = status
}

// 释放服务
func (t *Dux) release() {
	ants.Release()
}
