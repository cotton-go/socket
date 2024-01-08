package worker

import (
	"worker/pkg/app"
	"worker/pkg/config"
	"worker/pkg/log"
)

func NewServer(conf config.Config) *app.App {
	logger := log.NewLog(conf.Logger)
	return app.NewApp(
		app.WithServer(
			InitTCPServer(conf.TCP, logger),
			InitHTTPServer(conf.HTTP, logger),
			// InitGRPCServer(conf.GRPC, logger),
		),
	)
}
