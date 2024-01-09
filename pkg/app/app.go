package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cotton-go/socket/pkg/server"
)

// App 结构体表示一个应用程序，包含名称和服务器列表。
type App struct {
	name    string          // 应用程序的名称
	servers []server.Server // 应用程序关联的服务器列表
}

// Option 类型表示一个函数，用于修改 App 结构体的属性。
type Option func(*App)

// NewApp 根据传入的选项创建一个新的 App 实例。
func NewApp(opts ...Option) *App {
	a := &App{}
	for _, opt := range opts {
		opt(a)
	}

	return a
}

// Run启动应用程序并等待终止信号。
// 如果任何服务器启动失败或停止失败，则返回错误。
func (a *App) Run(ctx context.Context) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx) // 创建一个可以被取消的新上下文
	defer cancel()                        // 当函数返回时取消上下文

	signals := make(chan os.Signal, 1)                      // 创建一个接收操作系统信号的通道
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM) // 通知 goroutine 应该监听哪些信号

	// 在单独的 goroutine中启动每个服务器
	for _, srv := range a.servers {
		go func(srv server.Server) {
			err := srv.Start(ctx) // 启动服务器
			if err != nil {
				log.Printf("Server start err: %v", err) // 如果服务器启动失败，记录错误日志
			}
		}(srv)
	}

	// 等待终止信号或上下文取消
	select {
	case <-signals:
		// 收到终止信号
		log.Println("Received termination signal")
	case <-ctx.Done():
		// 上下文已取消
		log.Println("Context canceled")
	}

	// 以优雅的方式停止每个服务器
	for _, srv := range a.servers {
		err := srv.Stop(ctx) // 停止服务器
		if err != nil {
			log.Printf("Server stop err: %v", err) // 如果服务器停止失败，记录错误日志
		}
	}

	return nil // 在成功时返回nil
}
