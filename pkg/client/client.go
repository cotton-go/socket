package client

import (
	"fmt"
	"net"

	"worker/pkg/connection"
	"worker/pkg/event"
)

// Client 结构体表示一个客户端，包含一个连接对象
type Client struct {
	conn *connection.Connection
}

// New 方法用于创建一个新的客户端实例。
func New(addr string, opts ...connection.Options) (*Client, error) {
	// 连接服务器
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return nil, err
	}

	// 创建新的连接对象
	connectiond := connection.NewConnection(connection.WithConn(conn), connection.WithClient(true))
	for _, fn := range opts {
		fn(connectiond)
	}

	// 初始化 ID 和 WorkID
	connectiond.ID = 0
	connectiond.WorkID = 0

	// 监听连接初始化事件，设置 ID 和 WorkID
	connectiond.On(event.TopicByInitID, func(c *connection.Connection, _ event.Event) {
		connection.WithID(c.ID)
		connection.WithWorkID(c.WorkID)
	})

	// 创建新的客户端对象
	client := Client{
		conn: connectiond,
	}

	// 返回新创建的客户端对象
	return &client, nil
}

// Send 方法用于向指定主题发送消息。
func (c Client) Send(topic string, data any) error {
	// 调用连接对象的 Send 方法发送消息
	return c.conn.Send(topic, data)
}

// Subscription 方法用于订阅指定主题的消息，并在收到消息时调用 handdle 函数进行处理。
func (c Client) Subscription(topic string, handdle connection.EventHandle) {
	c.conn.On(topic, handdle)
}
