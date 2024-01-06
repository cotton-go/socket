package cache

import "worker/pkg/connection"

type ICache interface {
	Online(conn *connection.Connection) error
	Offline(conn *connection.Connection) error
	Find(id int64) *connection.Connection
}
