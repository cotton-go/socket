package connection

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"net"

	"github.com/bytedance/sonic"

	"worker/pkg/codec"
	"worker/pkg/event"
)

// EventHandle 是一个函数类型，用于处理事件
type EventHandle func(*Connection, event.Event)

// Connection 结构体表示一个连接
type Connection struct {
	ID       int64                    // 连接ID
	WorkID   int64                    // 工作ID
	closed   bool                     // 连接是否关闭
	conn     net.Conn                 // 网络连接
	ctx      context.Context          // 上下文对象
	cancel   context.CancelFunc       // 取消函数
	events   map[string][]EventHandle // 事件处理函数列表
	writeBuf chan event.Event         // 写缓冲区
	encBuf   *bufio.Writer            // 编码缓冲区
	enc      *gob.Encoder             // 编码器
	dec      *gob.Decoder             // 解码器
	codec    codec.ICodec             // 编解码器接口
	handle   EventHandle              // 事件处理函数
	isClient bool                     // 是否为客户端连接
}

// NewConnection 创建一个新的连接对象，并返回该对象的指针
func NewConnection(opts ...Options) *Connection {
	// 初始化连接对象
	conn := &Connection{
		writeBuf: make(chan event.Event, 100),
		events:   make(map[string][]EventHandle),
	}

	// 调用 makeOption 方法设置连接选项
	conn.makeOption(opts...)
	// 启动连接初始化协程
	go conn.init(opts...)
	// 返回连接对象指针
	return conn
}

// makeOption 根据传入的选项参数设置连接对象的属性
func (w *Connection) makeOption(opts ...Options) {
	// 定义默认选项
	var option = []Options{
		WithID(0),
		WithCodec(nil),
		WithHandle(nil),
		WithContext(context.Background()),
	}

	// 将传入的选项参数合并到默认选项中
	option = append(option, opts...)
	// 遍历选项列表，对每个选项进行处理
	for _, opt := range option {
		opt(w)
	}
}

// init 初始化连接对象，并根据参数设置连接选项。
func (c *Connection) init(opts ...Options) {
	// 打印连接 ID
	// fmt.Println("Connection init id=", c.ID)

	// 如果是客户端，则注册连接初始化事件的回调函数
	if c.isClient {
		c.On(event.TopicByInitID, func(_ *Connection, e event.Event) {
			// 对传入的数据进行 base64 解码
			b, err := base64.StdEncoding.DecodeString(e.Data.(string))
			if err != nil {
				fmt.Println("on connection init error[1001]", err)
				return
			}

			// 将解码后的数据反序列化为 Connection 对象
			var conn Connection
			if err := sonic.Unmarshal(b, &conn); err != nil {
				fmt.Println("on connection init error[1002]", err)
				return
			}

			// 将新连接的 ID 和工作 ID 赋值给当前连接对象
			c.ID = conn.ID
			c.WorkID = conn.WorkID
		})
	} else {
		// 如果是服务端，则将当前连接对象序列化为字节数组，并发送给客户端
		b, _ := sonic.Marshal(c)
		c.Send(event.TopicByInitID, b)
	}

	// 启动写入数据协程
	go c.write()
	// 启动读取数据协程
	go c.read()
}

// On 注册事件处理函数
func (c *Connection) On(event string, fn EventHandle) error {
	// 如果连接已关闭，则返回错误
	if c.closed {
		return errors.New("is closed")
	}

	// 将事件处理函数添加到对应事件的回调函数列表中
	c.events[event] = append(c.events[event], fn)

	// 返回 nil 表示成功注册事件处理函数
	return nil
}

// Emit 发送事件消息
func (c *Connection) Emit(topic string, data event.Event) {
	// 定义一个函数 handle,用于处理事件回调函数列表中的每个函数
	handle := func(handles ...EventHandle) {
		for _, fn := range handles {
			go func(fn EventHandle) {
				defer c.recover("connection emit")
				fn(c, data)
			}(fn)
		}
	}

	// 如果存在对应主题的事件回调函数列表，则执行 handle 函数
	if handles, ok := c.events[topic]; ok {
		go handle(handles...)
	}

	// 根据事件类型进行不同的处理
	switch data.Topic {
	case event.TopicByInitID:
		// 如果是客户端连接，则执行相应的操作
		if !c.isClient {
			return
		}

		// 如果存在对应主题的事件回调函数列表，则执行 handle 函数
		if handles, ok := c.events[event.TopicByInitID]; ok {
			go handle(handles...)
		}
	default:
		// 如果存在自定义的事件处理函数，并且连接未关闭，则执行该函数
		if c.handle != nil && !c.closed {
			c.handle(c, data)
		}
	}
}

// read 从连接中读取事件数据，并触发相应的事件处理函数。
func (c *Connection) read() {
	// 在函数退出前调用 recover 方法，防止 panic 导致的程序崩溃。
	defer c.recover("read over")

	// 在函数退出前触发关闭事件。
	defer func() {
		c.Emit(event.TopicByClose, event.Event{Topic: event.TopicByClose})
	}()

	// 循环读取事件数据，直到连接关闭或发生错误。
	for {
		select {
		case <-c.ctx.Done():
			// 当上下文被取消时，返回。
			return
		default:
			if c.closed {
				// 如果连接已关闭，则返回。
				return
			}

			var event event.Event
			// 从连接中解码事件数据。
			if err := c.dec.Decode(&event); err != nil {
				fmt.Println("read faild", err)
				// 如果解码失败，则返回。
				return
			}

			// 对事件数据进行编解码。
			event.Data, _ = c.codec.Decode(event.Data)
			// 触发相应的事件处理函数。
			c.Emit(event.Topic, event)
		}
	}
}

// write 方法用于向连接中写入数据。
func (c *Connection) write() {
	// 在函数退出前调用 recover 方法，防止 panic 导致的程序崩溃。
	defer c.recover("write over")
	for {
		select {
		case <-c.ctx.Done():
			// 当上下文被取消时，返回。
			return
		case buf := <-c.writeBuf:
			// 如果连接已关闭，则无法发送数据。
			if c.closed {
				fmt.Println("is closed not can send")
				return
			}

			// 对数据进行编解码。
			buf.Data, _ = c.codec.Encode(buf.Data)
			// 如果编码失败，则返回错误信息。
			if err := c.enc.Encode(buf); err != nil {
				// if err != io.EOF {
				//  // c.Close()
				//  return
				// }

				fmt.Println("write faild", err)
				return
			}
			// 将缓冲区中的数据写入连接。
			c.encBuf.Flush()
		}
	}
}

// Send 方法用于向连接发送数据。
//
// 参数：
// topic string - 要发送数据的话题名称。
// data any - 要发送的数据，可以是任意类型。
//
// 返回值：
// error - 如果连接已关闭，则返回错误信息；否则返回 nil。
func (c *Connection) Send(topic string, data any) error {
	if c.closed {
		return errors.New("is closed")
	}

	// 将事件写入缓冲区并返回 nil。
	c.writeBuf <- event.Event{
		Topic: topic,
		Data:  data,
	}
	return nil
}

// Close 方法用于关闭连接。
func (c *Connection) Close() error {
	// 如果连接已经关闭，则直接返回 nil。
	if c.closed {
		return nil
	}

	// 发送关闭事件。
	c.handle(c, event.Event{Topic: event.TopicByClose, Data: nil})

	// 锁定连接对象，保证线程安全。
	// c.locker.Lock()

	// 将连接状态设置为已关闭。
	c.closed = true

	// 取消所有未完成的请求。
	c.cancel()

	// 关闭底层的网络连接。
	c.conn.Close()

	// 返回错误信息。
	return nil
}

// recover 方法用于在连接中恢复错误处理。
func (c *Connection) recover(event string) {
	// 如果有错误发生，则打印错误信息。
	if err := recover(); err != nil {
		fmt.Println("connection recover", "event", event, "err", err)
	}
}
