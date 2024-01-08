package cache

import (
	"sync"

	"worker/pkg/connection"
)

// Memory 是一个结构体，包含一个读写锁和一个存储连接的 map
type Memory struct {
	lock  sync.RWMutex                     // 读写锁
	store map[int64]*connection.Connection // 存储连接的 map
}

// NewMemory 函数返回一个新的 Memory 实例
func NewMemory() ICache {
	return &Memory{
		store: make(map[int64]*connection.Connection), // 初始化存储连接的 map
	}
}

// Online 方法接受一个连接对象作为参数，将其添加到存储连接的 map 中，并返回 nil
func (m *Memory) Online(conn *connection.Connection) error {
	m.lock.Lock()           // 加锁
	defer m.lock.Unlock()   // 解锁
	m.store[conn.ID] = conn // 将连接对象添加到存储连接的 map 中
	return nil
}

// Offline 方法接受一个连接对象作为参数，从存储连接的 map 中删除该连接对象，并返回 nil
func (m *Memory) Offline(conn *connection.Connection) error {
	m.lock.Lock()         // 加锁
	defer m.lock.Unlock() // 解锁

	delete(m.store, conn.ID) // 从存储连接的 map 中删除该连接对象
	return nil
}

// Find 方法接受一个整型 id 作为参数，从存储连接的 map 中查找对应的连接对象，并返回该连接对象指针
func (m *Memory) Find(id int64) *connection.Connection {
	m.lock.RLock()         // 加读锁
	defer m.lock.RUnlock() // 解锁
	return m.store[id]     // 返回存储连接的 map 中对应 id 的连接对象指针
}
