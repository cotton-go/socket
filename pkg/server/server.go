package server

import (
	"context"
)

// Server 接口定义了服务器的基本操作
type Server interface {
	// Start 方法用于启动服务器，接收一个 context.Context 类型的参数，返回值为 error 类型
	Start(context.Context) error

	// Stop 方法用于停止服务器，接收一个 context.Context 类型的参数，返回值为 error 类型
	Stop(context.Context) error
}
