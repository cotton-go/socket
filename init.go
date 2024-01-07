package worker

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin/binding"

	"worker/pkg/app"
	"worker/pkg/codec"
	"worker/pkg/log"
	"worker/pkg/server/grpc"
	http1 "worker/pkg/server/http"
	"worker/pkg/server/tcp"
	"worker/pkg/worker"
)

var (
	work *worker.Worker
)

func InitTCPServer(conf tcp.Config, logger *log.Logger) app.Option {
	var icodec codec.ICodec
	switch strings.ToUpper(conf.Codec) {
	case "AESCBC":
		icodec = codec.NewAESCBC(conf.Secret)
	case "AESECB":
		icodec = codec.NewAESECB(conf.Secret)
	case "DESCBC":
		icodec = codec.NewDESCBC(conf.Secret)
	case "DESECB":
		icodec = codec.NewDESECB(conf.Secret)
	default:
		icodec = codec.NewDefault()
	}

	work = worker.NewWorker(
		// worker.WithCache(),
		worker.WithCodec(icodec),
		worker.WithContext(context.Background()),
	)

	return app.WithServer(tcp.NewServer(
		logger,
		tcp.WithServerWorker(work),
		tcp.WithServerHost(conf.Host),
		tcp.WithServerPort(conf.Port),
	))
}

func InitHTTPServer(conf http1.Config, logger *log.Logger) app.Option {
	router := http.NewServeMux()
	router.HandleFunc("/v1/find", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID int64 `json:"id" form:"id"`
		}

		b := binding.Default(r.Method, r.Header.Get("Content-Type"))
		if err := b.Bind(r, &req); err != nil {
			return
		}

		conn := work.Find(req.ID)
		if conn == nil {
			return
		}

		http1.JSON(w, http.StatusOK, http1.H{"data": conn, "code": 0, "msg": "ok"})
	})

	router.HandleFunc("/v1/send", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID    int64  `json:"id" form:"id"`
			Topic string `json:"topic" form:"topic"`
			Data  any    `json:"data" form:"data"`
		}

		b := binding.Default(r.Method, r.Header.Get("Content-Type"))
		if err := b.Bind(r, &req); err != nil {
			return
		}

		conn := work.Find(req.ID)
		if conn == nil {
			return
		}

		if err := conn.Send(req.Topic, req.Data); err != nil {
			return
		}

		http1.JSON(w, http.StatusOK, http1.H{"msg": "ok", "code": 0})
	})

	return app.WithServer(http1.NewServer(logger, router))
}

func InitGRPCServer(conf grpc.Config, logger *log.Logger) app.Option {
	return app.WithServer(grpc.NewServer(logger))
}
