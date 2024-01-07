package tcp

import (
	"context"

	"worker/pkg/worker"
)

type Option func(s *Server)

func WithServerHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

func WithServerPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

func WithServerStartBefore(fn func(context.Context)) Option {
	return func(s *Server) {
		s.startBefore = fn
	}
}

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
