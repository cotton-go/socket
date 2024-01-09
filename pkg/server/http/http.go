package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/cotton-go/socket/pkg/log"
)

// Server 结构体表示一个 HTTP 服务器
type Server struct {
	handler http.Handler // HTTP 请求处理函数
	httpSrv *http.Server // HTTP 服务器实例
	host    string       // 服务器主机名
	port    int          // 服务器端口号
	logger  *log.Logger  // 日志记录器实例
}

// Option 类型表示对 Server 结构体的选项设置函数，参数为指向 Server 的指针。
type Option func(s *Server)

// NewServer 创建一个新的服务器实例
// 参数：
// - logger *log.Logger 用于记录日志的日志记录器
// - handler http.Handler 用于处理HTTP请求的处理器
// - opts ...Option 可选的配置选项
// 返回值：
// - *Server 返回新创建的服务器实例
func NewServer(logger *log.Logger, handler http.Handler, opts ...Option) *Server {
	s := &Server{
		handler: handler,
		logger:  logger,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithServerHost 设置服务器主机名的选项
// 参数：
// - host string 主机名
// 返回值：
// - Option 一个函数，接受一个 *Server 参数并修改其主机名
func WithServerHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

// WithServerPort 设置服务器端口号的选项
// 参数：
// - port int 端口号
// 返回值：
// - Option 一个函数，接受一个 *Server 参数并修改其端口号
func WithServerPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

// Start 启动服务器的函数
// 参数：
// - ctx context.Context 上下文
// 返回值：
// - error 返回错误信息
func (s *Server) Start(ctx context.Context) error {
	// 创建一个新的 http.Server 并设置其地址和处理程序
	s.httpSrv = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: s.handler,
	}

	s.logger.Info("HTTP Server started listener", zap.String("host", s.host), zap.Int("port", s.port))
	// 开始监听并处理请求
	if err := s.httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		// 如果出现错误且不是服务器关闭错误，则记录错误并退出程序
		s.logger.Sugar().Fatalf("listen: %s\n", err)
	}

	// 返回 nil 表示没有错误
	return nil
}

// Stop 停止服务器的函数
// 参数：
// - ctx context.Context 上下文
// 返回值：
// - error 返回错误信息
func (s *Server) Stop(ctx context.Context) error {
	// 记录日志，表示正在关闭服务器
	s.logger.Sugar().Info("Shutting down server...")

	// 使用上下文和超时时间创建一个新的上下文
	// 这个新的上下文会告诉服务器它有5秒钟的时间来完成当前正在处理的请求
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 尝试关闭服务器
	if err := s.httpSrv.Shutdown(ctx); err != nil {
		// 如果在关闭服务器时发生错误，记录错误并强制关闭服务器
		s.logger.Sugar().Fatal("Server forced to shutdown:", err)
	}

	// 记录日志，表示服务器正在退出
	s.logger.Sugar().Info("Server exiting")

	// 返回nil表示没有错误
	return nil
}
