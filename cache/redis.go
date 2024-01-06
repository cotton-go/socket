package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"

	"worker/connection"
)

type Redis struct {
	ctx   context.Context
	store *redis.Client
}

func NewRedis(store *redis.Client) *Redis {
	return &Redis{
		ctx:   context.Background(),
		store: store,
	}
}

func (c Redis) Online(conn *connection.Connection) error {
	field, key := c.makeKey(conn.ID)
	value, _ := json.Marshal(conn)
	return c.store.HSet(c.ctx, key, field, string(value)).Err()
}

func (c Redis) Offline(conn *connection.Connection) error {
	field, key := c.makeKey(conn.ID)
	return c.store.HDel(c.ctx, key, field).Err()
}

func (c Redis) Find(id int64) *connection.Connection {
	var value struct {
		ID     int64
		WorkID int64
	}

	field, key := c.makeKey(id)
	bytes, err := c.store.HGet(c.ctx, key, field).Bytes()
	if err != nil {
		fmt.Println("redis find error[1001]", err)
		return nil
	}

	if err := json.Unmarshal(bytes, &value); err != nil {
		fmt.Println("redis find error[1002]", err)
	}

	return &connection.Connection{ID: value.ID, WorkID: value.WorkID}
}

func (c Redis) makeKey(id int64) (string, string) {
	field := strconv.Itoa(int(id))
	key := "connections"
	return field, key
}
