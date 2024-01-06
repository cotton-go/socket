package cache

import "worker/connection"

type ICache interface {
	Online(conn *connection.Connection) error
	Offline(conn *connection.Connection) error
	// Find(id ...int64)
}
