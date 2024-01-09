package socket

import (
	"github.com/cotton-go/socket/pkg/app"
	"github.com/cotton-go/socket/pkg/config"
	"github.com/cotton-go/socket/pkg/log"
	"github.com/cotton-go/socket/pkg/worker"
)

func NewServer(conf config.Config, opts ...worker.Options) *app.App {
	logger := log.NewLog(conf.Logger)
	return app.NewApp(app.WithServer(
		InitTCPServer(conf.TCP, logger, opts...),
		InitHTTPServer(conf.HTTP, logger),
		// InitGRPCServer(conf.GRPC, logger),
	))
}
