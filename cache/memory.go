package cache

import "worker/connection"

type Memory struct{}

func (m *Memory) Online(conn *connection.Connection) error {
	return nil
}

func (m *Memory) Offline(conn *connection.Connection) error {
	return nil
}
