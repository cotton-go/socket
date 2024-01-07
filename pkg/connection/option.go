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

// 定义一个Options类型的函数，用于配置Connection结构体
type Options func(*Connection)

// WithID函数，用于设置Connection的ID属性。
//
// 参数：
//   - value int64 ID值
//
// 返回值：
//   - Options 一个函数，该函数接收一个*Connection类型的参数w,并对其进行操作。
func WithID(value int64) Options {
	// 返回一个新的Options函数
	return func(w *Connection) {
		// 如果value小于1,则调用snowflake.Next()方法生成新的值
		if value < 1 {
			value = snowflake.Next()
		}

		// 将value赋值给w的ID属性
		w.ID = value
	}
}

// WithWorkID函数，用于设置Connection的WorkID属性。
//
// 参数：
//   - value int64工作ID值
//
// 返回值：
//   - Options 一个函数，该函数接收一个*Connection类型的参数w,并对其进行操作。
func WithWorkID(value int64) Options {
	// 返回一个新的Options函数
	return func(w *Connection) {
		// 如果value小于1,则调用snowflake.Next()方法生成新的值
		if value < 1 {
			value = snowflake.Next()
		}

		// 将value赋值给w的WorkID属性
		w.WorkID = value
	}
}

// WithConn函数，用于设置Connection的连接属性。
//
// 参数：
//   - conn net.Conn网络连接实例
//
// 返回值：
//   - Options 一个函数，该函数接收一个*Connection类型的参数c,并对其进行操作。
func WithConn(conn net.Conn) Options {
	// 返回一个新的Options函数
	return func(c *Connection) {
		// 为c的encBuf属性分配一个新的bufio.Writer实例，并将conn作为参数传入
		c.encBuf = bufio.NewWriter(conn)
		// 为c的enc属性分配一个新的gob.Encoder实例，并将c的encBuf属性作为参数传入
		c.enc = gob.NewEncoder(c.encBuf)
		// 为c的dec属性分配一个新的gob.Decoder实例，并将conn作为参数传入
		c.dec = gob.NewDecoder(conn)
	}
}

// WithCodec函数，用于设置Connection的codec属性。
//
// 参数：
//   - value codec.ICodec编解码器实例
//
// 返回值：
//   - Options 一个函数，该函数接收一个*Connection类型的参数w,并对其进行操作。
func WithCodec(value codec.ICodec) Options {
	// 返回一个新的Options函数
	return func(w *Connection) {
		// 如果传入的编解码器为空，则为w的codec属性分配一个新的codec.Default实例，并将&codec.Default{}作为参数传入
		if value == nil {
			value = &codec.Default{}
		}

		// 将传入的编解码器赋值给w的codec属性
		w.codec = value
	}
}

// WithClient函数，用于设置Connection的isClient属性。
//
// 参数：
//   - value bool 是否为客户端
//
// 返回值：
//   - Options 一个函数，该函数接收一个*Connection类型的参数c,并对其进行操作。
func WithClient(value bool) Options {
	// 返回一个新的Options函数
	return func(c *Connection) {
		// 将传入的isClient属性赋值给c的isClient属性
		c.isClient = value
	}
}

// WithHandle函数，用于设置Connection的handle属性。
//
// 参数：
//   - value EventHandle 事件句柄
//
// 返回值：
//   - Options 一个函数，该函数接收一个*Connection类型的参数w,并对其进行操作。
func WithHandle(value EventHandle) Options {
	// 返回一个新的Options函数
	return func(w *Connection) {
		// 如果传入的事件句柄为空，则为w的handle属性分配一个新的空函数，并将func(c *Connection, e event.Event) {}作为参数传入
		if value == nil {
			value = func(c *Connection, e event.Event) {}
		}

		// 将传入的事件句柄赋值给w的handle属性
		w.handle = value
	}
}

// WithContext函数，用于设置Connection的上下文。
//
// 参数：
//   - value context.Context 上下文
//
// 返回值：
//   - Options 一个函数，该函数接收一个*Connection类型的参数c,并对其进行操作。
func WithContext(value context.Context) Options {
	return func(c *Connection) {
		// 如果传入的上下文为空，则使用context.Background()创建一个新的上下文
		if value == nil {
			value = context.Background()
		}

		// 使用context.WithCancel方法创建一个新的上下文，该上下文可以被取消
		c.ctx, c.cancel = context.WithCancel(value)
	}
}
