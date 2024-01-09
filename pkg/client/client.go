package client

import (
	"fmt"
	"net"

	"github.com/cotton-go/socket/pkg/connection"
	"github.com/cotton-go/socket/pkg/event"
)

// Client 结构体表示一个客户端，包含一个连接对象
type Client struct {
	conn *connection.Connection
}

// New 方法用于创建一个新的客户端实例。
//
// 参数
//   - addr 表示服务器地址
//   - ...connection.Options 表示连接选项
//
// 返回值
//   - client 客户端实例
//   - error 错误信息
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

// Send方法用于发送消息
//
// 参数：
// - topic 主题
// - data 数据
//
// 返回值
// - error
func (c Client) Send(topic string, data any) error {
	// 调用连接对象的Send方法发送消息
	return c.conn.Send(topic, data)
}

// Subscription方法用于订阅主题
//
// 参数:
// - topic 主题
// - handdle 事件处理器
//
// 返回值: 无
func (c Client) Subscription(topic string, handdle connection.EventHandle) {
	c.conn.On(topic, handdle)
}
