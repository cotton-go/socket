package worker

import (
	"worker/pkg/app"
	"worker/pkg/config"
	"worker/pkg/log"
	"worker/pkg/worker"
)

func NewServer(conf config.Config, opts ...worker.Options) *app.App {
	logger := log.NewLog(conf.Logger)
	return app.NewApp(app.WithServer(
		InitTCPServer(conf.TCP, logger, opts...),
		InitHTTPServer(conf.HTTP, logger),
		// InitGRPCServer(conf.GRPC, logger),
	))
}
