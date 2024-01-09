package app

import "github.com/cotton-go/socket/pkg/server"

// WithServer 根据传入的服务器列表设置 App 实例的服务器属性。
func WithServer(servers ...server.Server) Option {
	return func(a *App) {
		a.servers = servers
	}
}

// WithName 根据传入的名称设置 App 实例的名称属性。
func WithName(name string) Option {
	return func(a *App) {
		a.name = name
	}
}
