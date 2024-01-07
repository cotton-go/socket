package config

import (
	"worker/pkg/log"
	"worker/pkg/server/grpc"
	"worker/pkg/server/http"
	"worker/pkg/server/tcp"
)

type Config struct {
	Env    string
	Logger log.Config
	HTTP   http.Config
	GRPC   grpc.Config
	TCP    tcp.Config
}
