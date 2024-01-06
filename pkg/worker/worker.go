package worker

import (
	"context"
	"fmt"
	"net"
	"sync"

	"worker/pkg/cache"
	"worker/pkg/codec"
	"worker/pkg/connection"
)

type Worker struct {
	id          int64
	count       int64
	lock        sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	connections map[int64]*connection.Connection
	cbuffer     chan *connection.Connection
	dbuffer     chan *connection.Connection
	cache       cache.ICache
	codec       codec.ICodec
	handle      connection.EventHandle
}

func NewWorker(opts ...Options) *Worker {
	w := &Worker{
		connections: make(map[int64]*connection.Connection),
		cbuffer:     make(chan *connection.Connection, 100),
		dbuffer:     make(chan *connection.Connection, 100),
	}

	go w.init(opts...)
	return w
}

func (w *Worker) makeOption(opts ...Options) {
	var option = []Options{
		WithID(0),
		WithCache(nil),
		WithCodec(nil),
		WithContext(context.Background()),
	}

	option = append(option, opts...)
	for _, opt := range option {
		opt(w)
	}
}

func (w *Worker) init(opts ...Options) {
	w.makeOption(opts...)
	go w.onConnection()
	go w.onDisconnect()
}

func (w *Worker) onConnection() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case conn := <-w.cbuffer:
			w.lock.Lock()
			id := conn.ID
			w.count += 1
			w.connections[id] = conn
			if err := w.cache.Online(conn); err != nil {
				fmt.Println("cache online error", err)
			}
			w.lock.Unlock()
		}
	}
}

func (w *Worker) onDisconnect() {
	for {
		select {
		case <-w.ctx.Done():
			for _, conn := range w.connections {
				w.cache.Offline(conn)
			}
			return
		case conn := <-w.dbuffer:
			id := conn.ID
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

func (w *Worker) Connection(conn net.Conn) *connection.Connection {
	c := connection.NewConnection(
		connection.WithID(0),
		connection.WithConn(conn),
		connection.WithWorkID(w.id),
		connection.WithCodec(w.codec),
		connection.WithContext(w.ctx),
		connection.WithHandle(w.handle),
	)

	w.cbuffer <- c
	return c
}

func (w *Worker) Disconnect(conn *connection.Connection) {
	w.dbuffer <- conn
}

func (w *Worker) ID() int64 {
	return w.id
}

func (w *Worker) Count() int64 {
	return w.count
}

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

func (w *Worker) Close() {
	w.cancel()
	// for _, v := range w.connections {

	// }
}