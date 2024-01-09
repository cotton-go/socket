package config

import (
	"github.com/cotton-go/socket/pkg/log"
	"github.com/cotton-go/socket/pkg/server/grpc"
	"github.com/cotton-go/socket/pkg/server/http"
	"github.com/cotton-go/socket/pkg/server/tcp"
)

type Config struct {
	Env    string
	Logger log.Config  `yaml:"Logger"`
	HTTP   http.Config `yaml:"HTTP"`
	GRPC   grpc.Config `yaml:"GRPC"`
	TCP    tcp.Config  `yaml:"TCP"`
}
