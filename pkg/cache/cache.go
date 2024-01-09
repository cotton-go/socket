package cache

import "github.com/cotton-go/socket/pkg/connection"

// ICache 是一个接口，定义了在线、离线和查找操作的方法
type ICache interface {
	// Online 方法接受一个连接对象作为参数，返回一个错误信息
	Online(conn *connection.Connection) error

	// Offline 方法接受一个连接对象作为参数，返回一个错误信息
	Offline(conn *connection.Connection) error

	// Find 方法接受一个整型 id 作为参数，返回一个连接对象指针
	Find(id int64) *connection.Connection
}
