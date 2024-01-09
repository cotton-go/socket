package connection

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/bytedance/sonic"

	"worker/pkg/codec"
	"worker/pkg/event"
)

// EventHandle 是一个函数类型，用于处理事件
type EventHandle func(*Connection, event.Event)

// Connection 结构体表示一个连接
type Connection struct {
	ID        int64                    // 连接ID
	WorkID    int64                    // 工作ID
	closed    bool                     // 连接是否关闭
	conn      net.Conn                 // 网络连接
	ctx       context.Context          // 上下文对象
	cancel    context.CancelFunc       // 取消函数
	events    map[string][]EventHandle // 事件处理函数列表
	writeBuf  chan event.Event         // 写缓冲区
	encBuf    *bufio.Writer            // 编码缓冲区
	enc       *gob.Encoder             // 编码器
	dec       *gob.Decoder             // 解码器
	codec     codec.ICodec             // 编解码器接口
	handle    EventHandle              // 事件处理函数
	isClient  bool                     // 是否为客户端连接
	heartbeat *time.Ticker             // 心跳间隔
}

// NewConnection 创建一个新的连接对象，并返回该对象的指针
//
// 参数：
//   - opts ...Options 可变参数，表示连接选项
//
// 返回值：
//   - *Connection 返回一个指向 Connection 类型的指针
func NewConnection(opts ...Options) *Connection {
	// 初始化连接对象
	conn := &Connection{
		writeBuf:  make(chan event.Event, 100),
		events:    make(map[string][]EventHandle),
		heartbeat: time.NewTicker(time.Second * 50),
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

// init 初始化连接
//
// 参数：
//   - c *Connection 连接对象指针
//   - opts ...Options 可变参数，表示可选的配置项
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
		go c.onHeartbeat()
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
//
// 参数：
//   - event string 事件名称
//   - fn EventHandle 事件处理函数
//
// 返回值：
//   - error 返回错误信息，如果成功则返回 nil
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

// Emit 发送事件
//
// 参数：
//   - topic string 主题
//   - data event.Event 事件数据
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

// read 函数用于从连接中读取事件数据，直到连接关闭或发生错误。
//
// 返回值：
//   - error 返回错误信息
//
// read方法用于从连接中读取事件数据
func (c *Connection) read() {
	// 在函数退出前调用 recover 方法，防止 panic 导致的程序崩溃
	defer c.recover("read over")

	// 在函数退出前触发关闭事件
	defer func() {
		c.Emit(event.TopicByClose, event.Event{Topic: event.TopicByClose})
	}()

	// 循环读取事件数据，直到连接关闭或发生错误
	for {
		select {
		case <-c.ctx.Done():
			// 当上下文被取消时，返回。
			return
		default:
			if c.closed {
				// 如果连接已关闭，则返回
				return
			}

			var e event.Event
			// 从连接中解码事件数据
			if err := c.dec.Decode(&e); err != nil {
				fmt.Println("read faild", err)
				// 如果解码失败，则返回
				return
			}

			// 对事件数据进行编解码
			e.Data, _ = c.codec.Decode(e.Data)
			// 触发相应的事件处理函数
			c.Emit(e.Topic, e)
		}
	}
}

// write 函数用于将缓冲区中的数据写入连接。
//
// 参数：无
//
// 返回值：无
func (c *Connection) write() {
	// 在函数退出前调用 recover 方法，防止 panic 导致的程序崩溃。
	defer c.recover("write over")
	for {
		select {
		case <-c.ctx.Done():
			// 当上下文被取消时，返回。
			return
		case buffer := <-c.writeBuf:
			// 如果连接已关闭，则无法发送数据。
			if c.closed {
				fmt.Println("is closed not can send")
				return
			}

			// 对数据进行编解码。
			buffer.Data, _ = c.codec.Encode(buffer.Data)
			// 如果编码失败，则返回错误信息。
			if err := c.enc.Encode(buffer); err != nil {
				fmt.Println("write faild", err)
				return
			}

			// 将缓冲区中的数据写入连接。
			c.encBuf.Flush()
		}
	}
}

// Send 函数用于向指定主题发送数据。
//
// 参数：
//   - topic string 主题名称
//   - data any 任意类型的数据
//
// 返回值：
//   - error 返回错误信息，如果连接已关闭则返回 "is closed" 错误
func (c *Connection) Send(topic string, data any) error {
	if c.closed {
		return errors.New("is closed")
	}

	// 将事件写入缓冲区并返回 nil。
	c.writeBuf <- event.Event{Topic: topic, Data: data}
	return nil
}

// onHeartbeat 处理心跳事件
//
// 参数：空
//
// 返回值：空
func (c *Connection) onHeartbeat() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-c.heartbeat.C:
			c.Send(event.TopicByHeartbeat, nil)
		}
	}
}

// Close 函数用于关闭连接。
//
// 参数：无
//
// 返回值：
//   - error 返回错误信息，如果关闭成功则返回 nil。
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

// recover 函数用于在发生错误时进行恢复操作。
//
// 参数：
//   - event string 事件名称
func (c *Connection) recover(event string) {
	// 如果有错误发生，则打印错误信息。
	if err := recover(); err != nil {
		fmt.Println("connection recover", "event", event, "err", err)
	}
}
