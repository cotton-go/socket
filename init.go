package socket

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cotton-go/socket/pkg/cache"
	"github.com/cotton-go/socket/pkg/codec"
	"github.com/cotton-go/socket/pkg/log"
	"github.com/cotton-go/socket/pkg/server"
	"github.com/cotton-go/socket/pkg/server/grpc"
	httpx "github.com/cotton-go/socket/pkg/server/http"
	"github.com/cotton-go/socket/pkg/server/tcp"
	"github.com/cotton-go/socket/pkg/worker"
)

var (
	work *worker.Worker
)

func InitTCPServer(conf tcp.Config, logger *log.Logger, opts ...worker.Options) server.Server {
	var cachex = cache.NewMemory()
	// if conf.Redis != nil {
	// cachex = cache.NewRedis(redis.NewClient(&redis.Options{
	// 	Addr:       conf.Redis.Addr,
	// 	Username:   conf.Redis.Username,
	// 	Password:   conf.Redis.Password,
	// 	MaxRetries: conf.Redis.MaxRetries,
	// 	DB:         conf.Redis.DB,
	// }))
	// }

	opts = append(opts,
		worker.WithCache(cachex),
		worker.WithCodec(codec.New(conf.Codec, conf.Secret)),
	)

	work = worker.NewWorker(opts...)
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

	router.POST("/v1/send", func(ctx *gin.Context) {
		var req struct {
			ID    int64  `form:"id" json:"id"`
			Topic string `form:"topic" json:"topic"`
			Data  any    `form:"data" json:"data"`
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

		if err := conn.Send(req.Topic, req.Data); err != nil {
			ctx.JSON(http.StatusOK, gin.H{"code": 1, "msg": "发送数据失败"})
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
