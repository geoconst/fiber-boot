package app

import (
	"github.com/google/wire"
)

var Provider = wire.NewSet(NewServer, NewApp)

type App struct {
	Server *Server
}

func NewApp(server *Server) *App {
	return &App{
		Server: server,
	}
}

// 启动应用
func (a *App) Run() {
	beforeRun()
	a.Server.Start()
}

// 自定义初始化代码
func beforeRun() {
	// 后续定时任务之类可以在这里初始化
}
