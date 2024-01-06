package cache

import (
	"sync"

	"worker/connection"
)

type Memory struct {
	lock  sync.RWMutex
	store map[int64]*connection.Connection
}

func NewMemory() *Memory {
	return &Memory{
		store: make(map[int64]*connection.Connection),
	}
}

func (m *Memory) Online(conn *connection.Connection) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.store[conn.ID] = conn
	return nil
}

func (m *Memory) Offline(conn *connection.Connection) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.store, conn.ID)
	return nil
}

func (m *Memory) Find(id int64) *connection.Connection {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.store[id]
}
