package config

import (
	"worker/pkg/log"
	"worker/pkg/server/grpc"
	"worker/pkg/server/http"
	"worker/pkg/server/tcp"
)

type Config struct {
	Env    string
	Logger log.Config  `yaml:"Logger"`
	HTTP   http.Config `yaml:"HTTP"`
	GRPC   grpc.Config `yaml:"GRPC"`
	TCP    tcp.Config  `yaml:"TCP"`
}
