package worker

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"worker/pkg/cache"
	"worker/pkg/codec"
	"worker/pkg/connection"
	"worker/pkg/event"
	"worker/pkg/log"
	"worker/pkg/server"
	"worker/pkg/server/grpc"
	httpx "worker/pkg/server/http"
	"worker/pkg/server/tcp"
	"worker/pkg/worker"
)

var (
	work *worker.Worker
)

func InitTCPServer(conf tcp.Config, logger *log.Logger) server.Server {
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

	var cachex = cache.NewMemory()
	// if conf.Redis != nil {
	// 	cachex = cache.NewRedis(redis.NewClient(&redis.Options{
	// 		Addr:       conf.Redis.Addr,
	// 		Username:   conf.Redis.Username,
	// 		Password:   conf.Redis.Password,
	// 		MaxRetries: conf.Redis.MaxRetries,
	// 		DB:         conf.Redis.DB,
	// 	}))
	// }

	work = worker.NewWorker(
		worker.WithCache(cachex),
		worker.WithCodec(icodec),
		worker.WithContext(context.Background()),
		worker.WithHandle(func(c *connection.Connection, e event.Event) {
			fmt.Println("client", c.ID, "topic", e.Topic, "data", e.Data, "count", work.Count())
		}),
	)

	socket := tcp.NewServer(
		logger,
		tcp.WithServerWorker(work),
		tcp.WithServerHost(conf.Host),
		tcp.WithServerPort(conf.Port),
	)

	return socket
}

func InitHTTPServer(conf httpx.Config, logger *log.Logger) server.Server {
	router := gin.Default()
	router.GET("/v1/find", func(ctx *gin.Context) {
		var req struct {
			ID int64 `json:"id" form:"id"`
		}

		if err := ctx.Bind(&req); err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": "获取参数错误"})
			return
		}

		conn := work.Find(req.ID)
		if conn == nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": "用户不在线"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": conn, "code": 0, "msg": "ok"})
	})

	router.GET("/v1/send", func(ctx *gin.Context) {
		var req struct {
			ID    int64  `json:"id" form:"id"`
			Topic string `json:"topic" form:"topic"`
			Data  any    `json:"data" form:"data"`
		}

		if err := ctx.Bind(&req); err != nil {
			return
		}

		conn := work.Find(req.ID)
		if conn == nil {
			return
		}

		if err := conn.Send(req.Topic, req.Data); err != nil {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"msg": "ok", "code": 0})
	})

	return httpx.NewServer(
		logger,
		router,
		httpx.WithServerHost(conf.Host),
		httpx.WithServerPort(conf.Port),
	)
}

func InitGRPCServer(conf grpc.Config, logger *log.Logger) server.Server {
	return grpc.NewServer(logger)
}