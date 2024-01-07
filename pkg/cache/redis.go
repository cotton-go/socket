package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"

	"worker/pkg/connection"
)

// Redis 结构体，包含一个上下文和一个 Redis 客户端
type Redis struct {
	ctx   context.Context
	store *redis.Client
}

// NewRedis 函数用于创建一个新的 Redis 实例
func NewRedis(store *redis.Client) *Redis {
	return &Redis{
		ctx:   context.Background(),
		store: store,
	}
}

// Online 方法将连接对象存储到 Redis 中，并返回错误信息
func (c Redis) Online(conn *connection.Connection) error {
	field, key := c.makeKey(conn.ID)
	value, _ := sonic.Marshal(conn)
	return c.store.HSet(c.ctx, key, field, string(value)).Err()
}

// Offline 方法从 Redis 中删除指定的连接对象，并返回错误信息
func (c Redis) Offline(conn *connection.Connection) error {
	field, key := c.makeKey(conn.ID)
	return c.store.HDel(c.ctx, key, field).Err()
}

// Find 方法用于在 Redis 中查找指定 ID 的连接对象，并返回该连接对象的指针。
//
// 参数：
// id int64:要查找的连接对象的 ID。
//
// 返回值：
// *connection.Connection:如果找到了指定 ID 的连接对象，则返回该连接对象的指针；否则返回 nil。
func (c Redis) Find(id int64) *connection.Connection {
	// 定义一个结构体变量 value,用于存储从 Redis 中获取到的连接对象的信息。
	var value struct {
		ID     int64
		WorkID int64
	}

	// 调用 makeKey 方法生成 Redis 中的 key 和 field。
	field, key := c.makeKey(id)

	// 从 Redis 中获取指定 key 对应的 bytes 数据。
	bytes, err := c.store.HGet(c.ctx, key, field).Bytes()
	if err != nil {
		fmt.Println("redis find error[1001]", err)
		return nil
	}

	// 将 bytes 数据反序列化为 value 结构体。
	if err := sonic.Unmarshal(bytes, &value); err != nil {
		fmt.Println("redis find error[1002]", err)
	}

	// 根据 value 结构体中的信息创建一个新的 connection.Connection 对象，并返回其指针。
	return &connection.Connection{ID: value.ID, WorkID: value.WorkID}
}

// makeKey 方法用于根据给定的 ID 生成 Redis 中的 key 和 field。
//
// 参数：
// id int64:要生成 key 和 field 的 ID。
//
// 返回值：
// string, string:生成的 Redis 中的 key 和 field。
func (c Redis) makeKey(id int64) (string, string) {
	// 将 ID 转换为字符串类型，并作为 field。
	field := strconv.Itoa(int(id))

	// 将 "connections" 作为 key。
	key := "connections"

	// 返回生成的 key 和 field。
	return field, key
}
