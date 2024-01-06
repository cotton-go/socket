package connection

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"

	"worker/constant"
	"worker/pkg/codec"
	"worker/pkg/event"
)

type EventHandle func(*Connection, event.Event)

type Connection struct {
	ID       int64
	WorkID   int64
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
	handle   EventHandle
}

func NewConnection(opts ...Options) *Connection {
	conn := &Connection{
		writeBuf: make(chan event.Event, 100),
		events:   make(map[string]EventHandle),
	}

	conn.makeOption(opts...)
	go conn.init(opts...)
	b, _ := json.Marshal(conn)
	conn.Send(constant.TopicByInitID, b)
	return conn
}

func (w *Connection) makeOption(opts ...Options) {
	var option = []Options{
		WithID(0),
		WithCodec(nil),
		WithHandle(nil),
		WithContext(context.Background()),
	}

	option = append(option, opts...)
	for _, opt := range option {
		opt(w)
	}
}

func (c *Connection) init(opts ...Options) {
	c.On(constant.TopicByClose, func(c *Connection, _ event.Event) {
		c.close()
	})

	c.On(constant.TopicByInitID, func(c *Connection, e event.Event) {
		b, err := base64.StdEncoding.DecodeString(e.Data.(string))
		if err != nil {
			fmt.Println("on connection init error[1001]", err)
			return
		}

		if err := json.Unmarshal(b, c); err != nil {
			fmt.Println("on connection init error[1002]", err)
		}
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

	switch data.Topic {
	case constant.TopicByInitID:
	default:
		c.handle(c, data)
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
	defer c.recover("write over")
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
	if err := recover(); err != nil {
		fmt.Println("connection recover", "event", event, "err", err)
	}
}
