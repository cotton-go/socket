package worker

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"worker/pkg/cache"
	"worker/pkg/client"
	"worker/pkg/codec"
	"worker/pkg/connection"
	"worker/pkg/event"
)

func TestWork(t *testing.T) {
	addr := "127.0.0.1:8080"
	icodec := codec.NewDESECB("1234567890123456")
	ctx, cancel := context.WithCancel(context.Background())
	work := NewWorker(
		WithCache(cache.NewRedis(redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"}))),
		WithContext(ctx),
		WithCodec(icodec),
		WithHandle(func(c *connection.Connection, e event.Event) {
			fmt.Println("on handle", "topic", e.Topic, "value", e.Data)
		}),
	)
	defer work.Close()

	timer := time.NewTimer(time.Second * 60)
	var wg sync.WaitGroup
	wg.Add(1)

	go t.Run("server", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("server recover", err)
			}
		}()

		// 监听端口
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			fmt.Println("Error listening:", err)
			return
		}

		defer listener.Close()
		wg.Done()
		fmt.Println("TCP Server is listening on", listener.Addr())

		for {
			select {
			case <-timer.C:
				cancel()
				return
			default:
				// 接受连接
				conn, err := listener.Accept()
				if err != nil {
					fmt.Println("Error accepting connection:", err)
					continue
				}

				// 启动一个 goroutine 处理连接
				c := work.Connection(conn)
				c.On("msg", func(c *connection.Connection, e event.Event) {
					fmt.Println("on msg", e.Data, "conn", c.ID)
					c.Send("rev", e.Data)
				})
			}
		}
	})

	go t.Run("client", func(t *testing.T) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("client recover", err)
			}
		}()

		wg.Wait()

		// 连接服务器
		c, err := client.New(addr, connection.WithContext(ctx), connection.WithCodec(icodec))
		if err != nil {
			return
		}

		c.Subscription("rev", func(c *connection.Connection, e event.Event) {
			fmt.Println("on rev", e.Data)
			fmt.Println("count", work.Count())
			// fmt.Println("work", work)
			// fmt.Println("connID", c.ID)
			fmt.Println("conn", work.Find(c.ID))
			fmt.Println()
		})

		// for {
		// 	select {
		// 	case <-timer.C:
		// 		// c.cancel()
		// 		return
		// 	default:
		msg := time.Now().String()
		c.Send("msg", msg)
		fmt.Println("send", msg)
		time.Sleep(time.Second * 5)
		// 	}
		// }
	})

	for {
		select {
		case <-timer.C:
			return
		default:
			fmt.Println("当前在线人数", work.Count())
			time.Sleep(time.Second)
		}
	}
}
