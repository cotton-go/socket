package tcp

import (
	"context"

	"worker/pkg/worker"
)

// Option 服务器配置选项类型
type Option func(s *Server)

// WithServerHost 设置服务器主机名
//
// 参数：
//   - host string 主机名
//
// 返回值：
//   - Option 返回一个配置选项，用于链式调用
func WithServerHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

// WithServerPort 设置服务器端口号
//
// 参数：
//   - port int 端口号
//
// 返回值：
//   - Option 返回一个配置选项，用于链式调用
func WithServerPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

// WithServerStartBefore 在服务器启动前执行的函数
//
// 参数：
//   - fn func(context.Context) 函数
//
// 返回值：
//   - Option 返回一个配置选项，用于链式调用
func WithServerStartBefore(fn func(context.Context)) Option {
	return func(s *Server) {
		s.startBefore = fn
	}
}

// WithServerStartAfter 在服务器启动后执行的函数
//
// 参数：
//   - fn func(context.Context) 函数
//
// 返回值：
//   - Option 返回一个配置选项，用于链式调用
func WithServerStartAfter(fn func(context.Context)) Option {
	return func(s *Server) {
		s.startAfter = fn
	}
}

// WithServerWorker 设置服务器的 worker,并返回一个 Option 类型的函数。
//
// 参数：
// - worker *worker.Worker:要设置的 worker 对象。
//
// 返回值：
// - func(*Server):一个接受 Server 类型参数的函数，用于将传入的 worker 对象赋值给 Server 的 worker 属性。
func WithServerWorker(worker *worker.Worker) Option {
	return func(s *Server) {
		s.worker = worker
	}
}
