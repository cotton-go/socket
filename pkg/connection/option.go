package connection

import (
	"bufio"
	"context"
	"encoding/gob"
	"net"

	"worker/pkg/codec"
	"worker/pkg/event"
	"worker/pkg/snowflake"
)

// Options 类型表示连接的选项。
type Options func(*Connection)

// WithID 方法用于设置 ID,如果值小于 1,则生成一个新的 snowflake ID。
func WithID(value int64) Options {
	return func(w *Connection) {
		if value < 1 {
			value = snowflake.Next()
		}

		w.ID = value
	}
}

// WithWorkID 方法用于设置工作 ID,如果值小于 1,则生成一个新的 snowflake ID。
func WithWorkID(value int64) Options {
	return func(w *Connection) {
		if value < 1 {
			value = snowflake.Next()
		}

		w.WorkID = value
	}
}

// WithConn 方法用于设置网络连接，并初始化编解码器和事件处理器。
func WithConn(conn net.Conn) Options {
	return func(c *Connection) {
		c.encBuf = bufio.NewWriter(conn)
		c.enc = gob.NewEncoder(c.encBuf)
		c.dec = gob.NewDecoder(conn)
	}
}

// WithCodec 方法用于设置编解码器。如果没有传入参数，则使用默认的编解码器。
func WithCodec(value codec.ICodec) Options {
	return func(w *Connection) {
		if value == nil {
			value = &codec.Default{}
		}

		w.codec = value
	}
}

// WithClient 方法用于设置是否为客户端连接。如果传入 true,则表示是客户端连接；否则为服务器端连接。
func WithClient(value bool) Options {
	return func(c *Connection) {
		c.isClient = value
	}
}

// WithHandle 方法用于设置事件处理器。如果没有传入参数，则使用空函数。
func WithHandle(value EventHandle) Options {
	return func(w *Connection) {
		if value == nil {
			value = func(c *Connection, e event.Event) {}
		}

		w.handle = value
	}
}

// WithContext 方法用于设置上下文，以便在需要时取消操作。如果没有传入参数，则使用默认的上下文。
func WithContext(value context.Context) Options {
	return func(c *Connection) {
		if value == nil {
			value = context.Background()
		}

		c.ctx, c.cancel = context.WithCancel(value)
	}
}
