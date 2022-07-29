package duxgo

import (
	"github.com/duxphp/duxgo/bootstrap"
	"github.com/duxphp/duxgo/global"
)

type Server struct {
	registerList []func(*Server)
}

func New() *Server {
	return &Server{}
}

// Register 注册启动服务
func (s *Server) Register(call func(*Server)) {
	s.registerList = append(s.registerList, call)
}

// SetConfigDir 设置配置目录
func (s *Server) SetConfigDir(dir string) {
	global.ConfigDir = dir
}

// Start 启动服务
func (s *Server) Start() {
	// 初始化启动服务
	t := bootstrap.New()
	// 注册核心服务
	t.RegisterCore()
	// 注册服务应用
	for _, call := range s.registerList {
		call(s)
	}
	// 注册WEB服务
	t.RegisterHttp()
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
