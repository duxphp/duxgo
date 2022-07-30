package duxgo

import (
	"github.com/duxphp/duxgo/bootstrap"
	"github.com/duxphp/duxgo/core"
)

type Server struct {
	registerApp     []func(*bootstrap.Bootstrap)
	registerService []func(*bootstrap.Bootstrap)
}

func New() *Server {
	return &Server{}
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

// Start 启动服务
func (s *Server) Start() {
	// 初始化启动服务
	t := bootstrap.New()
	// 注册核心服务
	t.RegisterCore()
	// 注册应用
	for _, call := range s.registerApp {
		call(t)
	}
	// 注册WEB服务
	t.RegisterHttp()
	// 注册服务
	for _, call := range s.registerService {
		call(t)
	}
	// 注册应用服务
	t.RegisterApp()
	// 启动队列服务
	go t.StartTask()
	// 启动WEB服务
	t.StartHttp()
	<-t.Ch
	// 停止队列服务
	t.StopTask()
	// 停止WEB服务
	t.StopHttp()
	// 释放其他服务
	t.Release()
}
