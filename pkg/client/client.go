package client

import (
	"fmt"
	"net"

	"worker/pkg/connection"
)

type Client struct {
	conn *connection.Connection
}

func New(addr string, opts ...connection.Options) (*Client, error) {
	// 连接服务器
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return nil, err
	}

	client := Client{
		conn: connection.NewConnection(connection.WithConn(conn)),
	}

	for _, fn := range opts {
		fn(client.conn)
	}

	return &client, nil
}

func (c Client) Send(topic string, data any) error {
	return c.conn.Send(topic, data)
}
