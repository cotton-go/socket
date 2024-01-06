package tcp

import (
	"context"
	"fmt"
	"net"

	"worker/pkg/log"
	"worker/pkg/worker"
)

type Option func(s *Server)

type Server struct {
	port   int
	host   string
	logger *log.Logger
	worker *worker.Worker
	Server net.Listener
	ctx    context.Context
	cancel context.CancelFunc
}

func NewServer(logger *log.Logger, opts ...Option) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Server{
		ctx:    ctx,
		cancel: cancel,
		logger: logger,
		worker: worker.NewWorker(),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

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

func WithServerWorker(worker *worker.Worker) Option {
	return func(s *Server) {
		s.worker = worker
	}
}

func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		s.logger.Sugar().Fatalf("Failed to listen: %v", err)
	}

	defer listener.Close()
	defer s.worker.Close()

	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			// 接受连接
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				continue
			}

			// 启动一个 goroutine 处理连接
			s.worker.Connection(conn)
		}
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.cancel()

	s.logger.Info("Server exiting")
	return nil
}
