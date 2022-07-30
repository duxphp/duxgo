package duxgo

import (
	"embed"
	"github.com/duxphp/duxgo/bootstrap"
	"github.com/duxphp/duxgo/core"
	"github.com/duxphp/duxgo/database"
	"github.com/duxphp/duxgo/task"
	"github.com/duxphp/duxgo/websocket"
)

// Server 服务管理
type Server struct {
	registerApp     []func(*bootstrap.Bootstrap) // 注册应用数据
	registerService []func(*bootstrap.Bootstrap) // 注册服务数据
	ServerStatus    ServerStatus                 // 系统服务状态
}

// ServerStatus 服务状态
type ServerStatus struct {
	database  bool // 数据库服务
	redis     bool // redis服务
	mongodb   bool // mongodb服务
	queue     bool // 队列服务
	websocket bool // websocket服务
}

// New 创建服务管理
func New() *Server {
	return &Server{
		ServerStatus: ServerStatus{
			database:  true,
			redis:     true,
			mongodb:   true,
			queue:     true,
			websocket: true,
		},
	}
}

// RegisterApp 注册应用
func (s *Server) RegisterApp(call func(*bootstrap.Bootstrap)) {
	s.registerApp = append(s.registerApp, call)
}

// RegisterService 注册服务
func (s *Server) RegisterService(call func(*bootstrap.Bootstrap)) {
	s.registerService = append(s.registerService, call)
}

// SetConfigDir 设置配置目录
func (s *Server) SetConfigDir(dir string) {
	core.ConfigDir = dir
}

// SetDatabaseStatus 设置数据库状态
func (s *Server) SetDatabaseStatus(status bool) {
	s.ServerStatus.database = status
}

// SetRedisStatus 设置redis状态
func (s *Server) SetRedisStatus(status bool) {
	s.ServerStatus.redis = status
}

// SetMongodbStatus 设置mongodb状态
func (s *Server) SetMongodbStatus(status bool) {
	s.ServerStatus.mongodb = status
}

// SetQueueStatus 设置队列状态
func (s *Server) SetQueueStatus(status bool) {
	s.ServerStatus.queue = status
}

// SetWebsocketStatus 设置websocket状态
func (s *Server) SetWebsocketStatus(status bool) {
	s.ServerStatus.websocket = status
}

//go:embed template/*
var tplFs embed.FS

// Start 启动服务
func (s *Server) Start() {
	// 设置系统模板
	core.TplFs = tplFs

	// 初始化启动服务
	t := bootstrap.New()

	// 注册核心服务
	t.RegisterCore()

	// 注册数据库
	if s.ServerStatus.database {
		database.GormInit()
	}
	// 注册redis
	if s.ServerStatus.redis {
		database.RedisInit()
	}
	// 注册mongodb
	if s.ServerStatus.mongodb {
		database.QmgoInit()
	}
	// 注册队列
	if s.ServerStatus.queue {
		task.Init()
	}
	// 注册websocket
	if s.ServerStatus.websocket {
		websocket.Init()
	}

	// 注册WEB服务
	t.RegisterHttp()

	// 注册应用
	for _, call := range s.registerApp {
		call(t)
	}

	// 注册服务
	for _, call := range s.registerService {
		call(t)
	}

	// 注册应用服务
	t.RegisterApp()

	// 启动队列服务
	if s.ServerStatus.queue {
		go t.StartTask()
		task.Add("ping", &map[string]any{})
	}
	// 启动WEB服务
	t.StartHttp()

	<-t.Ch

	// 停止队列服务
	if s.ServerStatus.queue {
		t.StopTask()
	}

	// 停止WEB服务
	t.StopHttp()

	// 释放websocket
	if s.ServerStatus.websocket {
		websocket.ReleaseSocket()
	}

	// 释放系统服务
	t.Release()
}
