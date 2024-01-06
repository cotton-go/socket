package tcp

import (
	"context"
	"fmt"
	"testing"
	"time"

	"worker/constant"
	"worker/pkg/client"
	"worker/pkg/codec"
	"worker/pkg/connection"
	"worker/pkg/event"
	"worker/pkg/log"
	"worker/pkg/worker"
)

func TestNewTCP(t *testing.T) {
	auth := false
	prot := 8080
	host := "127.0.0.1"
	icodec := codec.NewDESECB("1234567890123456")
	ctx, cancel := context.WithCancel(context.Background())
	work := worker.NewWorker(
		worker.WithContext(ctx),
		worker.WithCodec(icodec),
		worker.WithHandle(func(c *connection.Connection, e event.Event) {
			fmt.Println("on handle", "topic", e.Topic, "value", e.Data)
			if e.Topic == constant.TopicByClose {
				fmt.Println("收到断开连接请求", c.ID)
				return
			}

			if e.Topic == constant.TopicByLogin {
				auth = true
				fmt.Println("收到登陆认证请求", c.ID)
				c.Send("logind", e.Data)
				return
			}

			//
			if !auth {
				fmt.Println("未收到登陆认证请求，立即退出", c.ID)
				c.Close()
				return
			}

			// fmt.Println("on handle", "topic", e.Topic, "value", e.Data)
			c.On("msg", func(c *connection.Connection, e event.Event) {
				fmt.Println("on msg", e.Data, "conn", c.ID)
				c.Send("rev", e.Data)
			})
		}),
	)

	defer work.Close()

	server := NewServer(log.NewLog(log.Config{}), WithServerWorker(work), WithServerHost(host), WithServerPort(prot))
	go server.Start(ctx)

	t.Run("client", func(t *testing.T) {
		time.Sleep(time.Second * 10)

		handle := connection.WithHandle(func(c *connection.Connection, e event.Event) {
			fmt.Println("on handle 1", "topic", e.Topic, "value", e.Data)
			if e.Topic == "logind" {
				fmt.Println("登陆成功", c.ID)
				return
			}

			c.On("rev", func(c *connection.Connection, e event.Event) {
				fmt.Println("on rev", e.Data)
				fmt.Println("count", work.Count())
				fmt.Println()
			})
		})

		addr := fmt.Sprintf("%v:%v", host, prot)
		client, err := client.New(addr, connection.WithContext(ctx), connection.WithCodec(icodec), handle)
		if err != nil {
			fmt.Println("Error connecting to server:", err)
			return
		}

		// client.Send(constant.TopicByLogin, "")
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg := time.Now().String()
				if err := client.Send("msg", msg); err != nil {
					fmt.Println("send error", err)
					return
				}

				// fmt.Println("send", msg)
				time.Sleep(time.Second * 5)
			}
		}

		// // 连接服务器
		// conn, err := net.Dial("tcp", addr)
		// if err != nil {
		// 	fmt.Println("Error connecting to server:", err)
		// 	return
		// }

		// defer conn.Close()
		// c := connection.NewConnection(connection.WithContext(ctx), connection.WithConn(conn), connection.WithCodec(icodec))
		// c.On("rev", func(c *connection.Connection, e event.Event) {
		// 	fmt.Println("on rev", e.Data)
		// 	fmt.Println("count", work.Count())
		// 	// fmt.Println("conn", work.Find(c.ID))
		// 	fmt.Println()
		// })

		// for {
		// 	select {
		// 	case <-ctx.Done():
		// 		return
		// 	default:
		// 		msg := time.Now().String()
		// 		c.Send("msg", msg)
		// 		fmt.Println("send", msg)
		// 		time.Sleep(time.Second * 5)
		// 	}
		// }
	})

	time.Sleep(time.Minute * 1)
	cancel()
}