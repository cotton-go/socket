package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"

	"go.uber.org/zap"

	"github.com/cotton-go/socket/pkg/log"
	"github.com/cotton-go/socket/pkg/worker"
)

// Server 结构体表示一个服务器实例
type Server struct {
	port        int                   // 服务器监听的端口号
	host        string                // 服务器监听的主机名或IP地址
	logger      *log.Logger           // 日志记录器
	worker      *worker.Worker        // 工作线程池
	Server      net.Listener          // 网络监听器
	ctx         context.Context       // 上下文对象
	cancel      context.CancelFunc    // 取消函数
	startBefore func(context.Context) // 在启动前执行的回调函数
	startAfter  func(context.Context) // 在启动后执行的回调函数
	stopBefore  func(context.Context) // 在停止前执行的回调函数
	stopAfter   func(context.Context) // 在停止后执行的回调函数
}

// NewServer 创建一个新的服务器实例
//
// 参数:
//   - logger: *log.Logger,日志记录器
//   - opts: []Option,可选参数列表
//
// 返回值：
//   - *Server,新创建的服务器实例
func NewServer(logger *log.Logger, opts ...Option) *Server {
	// 创建一个可取消的上下文对象
	ctx, cancel := context.WithCancel(context.Background())
	// 初始化服务器实例
	s := &Server{
		ctx:    ctx,
		cancel: cancel,
		logger: logger,
		worker: worker.NewWorker(),
	}

	// 遍历传入的选项函数，并执行它们
	for _, opt := range opts {
		opt(s)
	}

	// 返回新创建的服务器实例
	return s
}

// Start 启动服务器，并在指定的上下文中运行。
//
// 参数：
//   - ctx context.Context - 用于控制服务器启动和停止的上下文。
//
// 返回值：
//   - error - 如果启动过程中出现错误，则返回错误；否则返回 nil。
func (s *Server) Start(ctx context.Context) error {
	// 如果服务器为空，则返回错误
	if s == nil {
		return errors.New("Server is nil")
	}

	// 如果 startBefore 不为空，则在启动之前执行该函数
	if s.startBefore != nil {
		s.startBefore(ctx)
	}

	// 监听指定的地址和端口
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		s.logger.Sugar().Fatalf("Failed to listen: %v", err)
		return err
	}

	// 如果 startAfter 不为空，则在启动之后执行该函数
	if s.startAfter != nil {
		s.startAfter(ctx)
	}

	// 关闭监听器和工作线程
	defer listener.Close()
	defer s.worker.Close()

	s.logger.Info("TCP Server startd listener", zap.String("host", s.host), zap.Int("port", s.port))

	// 循环处理客户端连接
	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("Server exiting[1002]")
			return nil
		default:
			// 以非阻塞方式接受新的连接
			conn, err := listener.Accept()
			if err != nil {
				s.logger.Error("Error accepting connection:", zap.Error(err))
				continue
			}

			// 启动一个 goroutine 来处理连接
			go func(conn net.Conn) {
				s.worker.Connection(conn)
			}(conn)
		}
	}
}

// Stop 停止服务器并执行必要的操作
//
// 参数：
//   - ctx context.Context - 用于控制服务器启动和停止的上下文。
//
// 返回值：
//   - error - 如果停止过程中出现错误，则返回错误；否则返回 nil。
func (s *Server) Stop(ctx context.Context) error {
	// 在函数退出前执行的代码块，用于确保在停止服务器之前执行必要的操作
	defer func() {
		if s.stopAfter != nil {
			s.stopAfter(ctx)
		}
	}()

	// 如果 stopBefore 不为空，则在停止服务器之前执行该函数
	if s.stopBefore != nil {
		s.stopBefore(ctx)
	}

	// 取消服务器的所有 goroutine
	s.cancel()

	// fmt.Println("count", s.worker.Count())

	// 记录错误日志，表示服务器正在退出
	s.logger.Error("Server exiting[1001]")

	// 返回 nil 表示没有错误发生
	return nil
}
