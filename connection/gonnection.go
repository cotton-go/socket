package connection

import (
	"bufio"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net"

	"worker/codec"
	"worker/event"
)

type EventHandle func(*Connection, event.Event)

type Connection struct {
	ID       int64
	conn     net.Conn
	ctx      context.Context
	cancel   context.CancelFunc
	events   map[string]EventHandle
	writeBuf chan event.Event
	encBuf   *bufio.Writer
	enc      *gob.Encoder
	dec      *gob.Decoder
	codec    codec.ICodec
	closed   bool
}

func NewConnection(opts ...Options) *Connection {
	c := &Connection{
		writeBuf: make(chan event.Event, 100),
		events:   make(map[string]EventHandle),
	}

	go c.init(opts...)
	return c
}

func (w *Connection) makeOption(opts ...Options) {
	var option = []Options{
		WithID(0),
		WithCodec(nil),
		WithContext(context.Background()),
	}

	option = append(option, opts...)
	for _, opt := range option {
		opt(w)
	}
}

func (c *Connection) init(opts ...Options) {
	c.makeOption(opts...)
	c.On("__close__", func(c *Connection, _ event.Event) {
		c.close()
	})

	go c.write()
	go c.read()
}

func (c *Connection) On(event string, fn EventHandle) error {
	if c.closed {
		return errors.New("is closed")
	}

	c.events[event] = fn
	return nil
}

func (c *Connection) Emit(event string, data event.Event) {
	if fn, ok := c.events[event]; ok {
		go func() {
			defer c.recover("emit")
			fn(c, data)
		}()
	}
}

func (c *Connection) read() {
	defer c.recover("read over")
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			if c.closed {
				return
			}

			var event event.Event
			if err := c.dec.Decode(&event); err != nil {
				if err != io.EOF {
					c.close()
					return
				}
				fmt.Println("read faild", err)
				continue
			}

			event.Data, _ = c.codec.Decode(event.Data)
			c.Emit(event.Topic, event)
		}
	}
}

func (c *Connection) write() {
	defer c.recover("read")
	for {
		select {
		case <-c.ctx.Done():
			return
		case buf := <-c.writeBuf:
			if c.closed {
				return
			}

			buf.Data, _ = c.codec.Encode(buf.Data)
			if err := c.enc.Encode(buf); err != nil {
				if err != io.EOF {
					c.close()
					return
				}

				fmt.Println("write faild", err)
				continue
			}
			c.encBuf.Flush()
		}
	}
}

func (c *Connection) Send(topic string, data any) {
	c.writeBuf <- event.Event{
		Topic: topic,
		Data:  data,
	}
}

func (c *Connection) close() {
	c.cancel()
	c.conn.Close()
	c.closed = true
}

func (c *Connection) recover(event string) {
	if r := recover(); r != nil {
		fmt.Println("r", r)
	}
}
