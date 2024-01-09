package worker

import (
	"context"
	"fmt"
	"net"
	"sync"

	"worker/pkg/cache"
	"worker/pkg/codec"
	"worker/pkg/connection"
	"worker/pkg/event"
	"worker/pkg/registry"
)

// Worker代表一个具有其属性和方法的工作对象。
type Worker struct {
	id          int64                            // 工作对象的ID
	count       int64                            // 工作对象处理的任务数
	lock        sync.RWMutex                     // 读写锁，用于线程安全
	ctx         context.Context                  // 用于取消和超时的上下文
	cancel      context.CancelFunc               // 用于停止工作对象的取消函数
	connections map[int64]*connection.Connection // 活动连接的映射表
	cbuffer     chan *connection.Connection      // 传入连接的缓冲区
	dbuffer     chan *connection.Connection      // 传出连接的缓冲区
	cache       cache.ICache                     // 存储数据的缓存接口
	codec       codec.ICodec                     // 编码和解码数据的编解码器接口
	handle      connection.EventHandle           // 事件处理器，用于处理事件
	registry    registry.Registry                // 注册中心处理器，用于注册服务
}

// NewWorker 方法用于创建一个新的 Worker 实例。
//
// 参数：
// - opts ...Options:可变参数，表示可选的配置项。
//
// 返回值：
// - *Worker:新创建的 Worker 实例。
func NewWorker(opts ...Options) *Worker {
	// 初始化 Worker 实例
	w := &Worker{
		connections: make(map[int64]*connection.Connection),
		cbuffer:     make(chan *connection.Connection, 100),
		dbuffer:     make(chan *connection.Connection, 100),
	}

	// 启动初始化函数
	go w.init(opts...)

	// 返回新创建的 Worker 实例
	return w
}

// makeOption 方法用于设置 Worker 的选项。
//
// 参数：
// - opts ...Options:可变参数，表示要设置的选项。
//
// 返回值：无
func (w *Worker) makeOption(opts ...Options) {
	// 初始化默认选项
	var option = []Options{
		WithID(0),
		WithCache(nil),
		WithCodec(nil),
		WithContext(context.Background()),
	}

	// 将传入的选项添加到默认选项中
	option = append(option, opts...)

	// 遍历选项并应用到 Worker 上
	for _, opt := range option {
		opt(w)
	}
}

// init 方法用于初始化 Worker。
//
// 参数：
// - opts ...Options:可变参数，表示要设置的选项。
//
// 返回值：无
func (w *Worker) init(opts ...Options) {
	// 调用 makeOption 方法设置选项
	w.makeOption(opts...)

	// 监听服务注册事件处理函数
	go w.onRegister()

	// 启动连接事件处理函数
	go w.onConnection()

	// 启动断开连接事件处理函数
	go w.onDisconnect()
}

func (w *Worker) onRegister() {
	// w.registry.Register()

	// 循环处理断开连接事件
	for {
		select {
		case <-w.ctx.Done():
			w.registry.Deregister()
			return
		}
	}
}

// onConnection 方法用于处理连接事件。
//
// 参数：无
//
// 返回值：无
func (w *Worker) onConnection() {
	// 循环处理连接事件
	for {
		select {
		case <-w.ctx.Done():
			// 如果上下文被取消，则退出函数
			return
		case conn := <-w.cbuffer:
			// 从缓冲区中获取连接对象
			w.lock.Lock()
			id := conn.ID
			// 计数器加一
			w.count += 1
			// 将连接对象添加到连接列表中
			w.connections[id] = conn
			// 将连接对象设置为在线状态
			if err := w.cache.Online(conn); err != nil {
				// 如果设置在线状态失败，则输出错误信息
				fmt.Println("cache online error", err)
			}
			w.handle(conn, event.Event{Topic: event.TopicByLogin})
			w.lock.Unlock()
		}
	}
}

// onDisconnect 方法用于处理断开连接事件。
//
// 参数：无
//
// 返回值：无
func (w *Worker) onDisconnect() {
	// 循环处理断开连接事件
	for {
		select {
		case <-w.ctx.Done():
			// 如果上下文被取消，则将所有连接对象设置为离线状态并退出函数
			for _, conn := range w.connections {
				w.cache.Offline(conn)
			}
			return
		case conn := <-w.dbuffer:
			id := conn.ID
			// 如果连接对象存在于连接列表中，则将其设置为离线状态
			if _, ok := w.connections[id]; ok {
				w.lock.Lock()
				w.count -= 1
				delete(w.connections, id)
				if err := w.cache.Offline(conn); err != nil {
					fmt.Println("cache offline error", err)
				}

				w.lock.Unlock()
			}
		}
	}
}

// Connection 方法用于创建一个新的连接对象，并将其添加到工作器中。
//
// 参数：
// conn net.Conn:一个 net.Conn 类型的连接对象。
//
// 返回值：
// *connection.Connection:一个 connection.Connection 类型的连接对象。
func (w *Worker) Connection(conn net.Conn) *connection.Connection {
	// 创建一个新的连接对象，并设置其属性
	c := connection.NewConnection(
		connection.WithID(0),
		connection.WithConn(conn),
		connection.WithWorkID(w.id),
		connection.WithCodec(w.codec),
		connection.WithContext(w.ctx),
		connection.WithHandle(w._handle),
	)

	// 当连接关闭时，将连接对象发送到工作器的缓冲区中
	c.On(event.TopicByClose, func(_ *connection.Connection, e event.Event) {
		w.dbuffer <- c
	})

	// 将连接对象发送到工作器的缓冲区中
	w.cbuffer <- c
	return c
}

// _handle 方法用于处理连接事件，并根据事件类型执行相应的操作。
//
// 参数：
// conn *connection.Connection: 一个指向 connection.Connection 类型的指针，表示要处理的连接。
// e event.Event: 一个 event.Event 类型的变量，表示要处理的事件。
//
// 返回值：
// 无返回值。
func (w *Worker) _handle(conn *connection.Connection, e event.Event) {
	// 如果事件主题是关闭连接，将连接添加到缓冲区中，并打印关闭信息。
	if e.Topic == event.TopicByClose {
		w.dbuffer <- conn
		fmt.Println("connection close", conn.ID)
	}

	// 如果存在 handle 方法，则调用该方法处理事件。
	if w.handle != nil {
		w.handle(conn, e)
	}
}

// Disconnect 方法用于断开与指定连接的连接。
//
// 参数：
// conn *connection.Connection: 一个指向 connection.Connection 类型的指针，表示要断开的连接。
func (w *Worker) Disconnect(conn *connection.Connection) {
	w.dbuffer <- conn
}

// ID 方法用于获取 Worker 实例的唯一标识符。
//
// 返回值：
// int64: Worker 实例的唯一标识符。
func (w *Worker) ID() int64 {
	return w.id
}

// Count 方法用于获取 Worker 实例的任务数量。
//
// 返回值：
// int64: Worker 实例的任务数量。
func (w *Worker) Count() int64 {
	return w.count
}

// Connections 方法用于获取 Worker 实例的所有连接。
//
// 返回值：
// map[int64]*connection.Connection: 一个映射，其中键为 int64 类型的连接 ID,值为对应的 connection.Connection 类型指针。
func (w *Worker) Connections() map[int64]*connection.Connection {
	return w.connections
}

// Find 方法用于在 Worker 实例的连接中查找指定 ID 的连接。
//
// 参数：
// id int64: 要查找的连接的 ID。
//
// 返回值：
// *connection.Connection: 如果找到了指定 ID 的连接，则返回对应的 connection.Connection 类型指针；否则返回 nil。
func (w *Worker) Find(id int64) *connection.Connection {
	w.lock.RLock()
	defer w.lock.RUnlock()

	if conn, ok := w.connections[id]; ok {
		return conn
	}

	if conn := w.cache.Find(id); conn != nil {
		return conn
	}

	fmt.Println("connection not found", id)
	return nil
}

// Close 方法用于关闭 Worker 实例的所有连接。
func (w *Worker) Close() {
	w.cancel()
}
