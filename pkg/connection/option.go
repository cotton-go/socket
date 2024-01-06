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

type Options func(*Connection)

func WithID(value int64) Options {
	return func(w *Connection) {
		if value < 1 {
			value = snowflake.Next()
		}

		w.ID = value
	}
}

func WithWorkID(value int64) Options {
	return func(w *Connection) {
		if value < 1 {
			value = snowflake.Next()
		}

		w.WorkID = value
	}
}

func WithConn(conn net.Conn) Options {
	return func(c *Connection) {
		c.encBuf = bufio.NewWriter(conn)
		c.enc = gob.NewEncoder(c.encBuf)
		c.dec = gob.NewDecoder(conn)
	}
}

func WithCodec(value codec.ICodec) Options {
	return func(w *Connection) {
		if value == nil {
			value = &codec.Default{}
		}

		w.codec = value
	}
}

func WithHandle(value EventHandle) Options {
	return func(w *Connection) {
		if value == nil {
			value = func(c *Connection, e event.Event) {}
		}

		w.handle = value
	}
}

func WithContext(value context.Context) Options {
	return func(c *Connection) {
		if value == nil {
			value = context.Background()
		}

		c.ctx, c.cancel = context.WithCancel(value)
	}
}
