package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	"worker/pkg/log"
)

// Server 结构体表示一个 gRPC 服务器
type Server struct {
	*grpc.Server             // gRPC 服务器实例
	host         string      // 服务器主机名
	port         int         // 服务器端口号
	logger       *log.Logger // 日志记录器实例
}

// NewServer 创建一个新的 Server 实例。
//
// 参数：
// - logger *log.Logger 日志记录器实例
// - opts ...Option 一个或多个 Option 类型的函数，用于对 Server 进行配置。
//
// 返回值：
// - *Server 返回一个新创建的 Server 实例。
func NewServer(logger *log.Logger, opts ...Option) *Server {
	// 创建一个新的 Server 实例。
	s := &Server{
		Server: grpc.NewServer(),
		logger: logger,
	}

	// 对 Server 进行配置。
	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start 是 Server 的启动方法。它会启动服务器并监听指定的主机和端口。如果启动失败，将记录错误信息并终止程序。
//
// 参数：
// - ctx context.Context 一个上下文对象，用于控制服务器的生命周期。在调用此方法时传入该对象。
//
// 返回值：
// - error 如果启动过程中发生错误，则返回一个错误信息；否则返回 nil。
func (s *Server) Start(ctx context.Context) error {
	// 在指定的主机和端口上监听连接请求。
	li, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		// 如果监听失败，记录错误信息并终止程序。
		s.logger.Sugar().Fatalf("Failed to listen: %v", err)
	}
	defer li.Close() // 确保在函数结束时关闭监听器。

	// 将监听器传递给 gRPC 服务器以便接受连接请求并处理请求。如果处理失败，记录错误信息并终止程序。
	if err = s.Server.Serve(li); err != nil {
		// 如果处理失败，记录错误信息并终止程序。
		s.logger.Sugar().Fatalf("Failed to serve: %v", err)
	}
	return nil // 如果一切正常，返回 nil。
}

// Stop 停止服务器
// 参数：
//   - ctx context.Context 上下文
//
// 返回值：
//   - error 返回错误信息
func (s *Server) Stop(ctx context.Context) error {
	// 使用 WithTimeout 函数设置超时时间为5秒
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	// 调用 GracefulStop 方法来优雅地停止服务器
	s.Server.GracefulStop()

	// 记录日志信息
	s.logger.Info("Server exiting")
	// 返回 nil 表示没有错误发生
	return nil
}
