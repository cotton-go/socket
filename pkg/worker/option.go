package worker

import (
	"context"

	"github.com/cotton-go/socket/pkg/cache"
	"github.com/cotton-go/socket/pkg/codec"
	"github.com/cotton-go/socket/pkg/connection"
	"github.com/cotton-go/socket/pkg/event"
	"github.com/cotton-go/socket/pkg/snowflake"
)

// Options 是一个函数类型，用于接收一个 *Worker 实例作为参数，并对其进行配置。
type Options func(*Worker)

// WithID 函数用于设置 Worker 实例的 ID。
//
// 参数：
// value int64: 要设置的 ID 值。如果小于 1,则会自动生成一个新的 ID。
//
// 返回值：
// Options: 一个闭包函数，接受一个 Worker 实例作为参数，并将其 ID 设置为指定的值。
func WithID(value int64) Options {
	return func(w *Worker) {
		if value < 1 {
			value = snowflake.Next()
		}
		w.id = value
	}
}

// WithCache 函数用于设置 Worker 实例的缓存。
//
// 参数：
// value cache.ICache: 要设置的缓存对象，实现了 cache.ICache 接口。如果为 nil,则会创建一个新的内存缓存对象。
//
// 返回值：
// Options: 一个闭包函数，接受一个 Worker 实例作为参数，并将其缓存设置为指定的值。
func WithCache(value cache.ICache) Options {
	return func(w *Worker) {
		if value == nil {
			value = cache.NewMemory()
		}

		w.cache = value
	}
}

// WithCodec 函数用于设置 Worker 实例的编解码器。
//
// 参数：
// value codec.ICodec: 要设置的编解码器对象，实现了 codec.ICodec 接口。如果为 nil,则会使用默认的编解码器。
//
// 返回值：
// Options: 一个闭包函数，接受一个 Worker 实例作为参数，并将其编解码器设置为指定的值。
func WithCodec(value codec.ICodec) Options {
	return func(w *Worker) {
		if value == nil {
			value = &codec.Default{}
		}

		w.codec = value
	}
}

// WithHandle 函数用于设置 Worker 实例的事件处理器。
//
// 参数：
// value connection.EventHandle: 要设置的事件处理器对象，实现了 connection.EventHandle 接口。如果为 nil,则会使用空函数作为默认处理器。
//
// 返回值：
// Options: 一个闭包函数，接受一个 Worker 实例作为参数，并将其事件处理器设置为指定的值。
func WithHandle(value connection.EventHandle) Options {
	return func(w *Worker) {
		if value == nil {
			value = func(c *connection.Connection, e event.Event) {}
		}

		w.handle = value
	}
}

// WithContext 函数用于设置 Worker 实例的上下文。
//
// 参数：
// value context.Context: 要设置的上下文对象，实现了 context.Context 接口。如果为 nil,则会使用默认的上下文对象。
//
// 返回值：
// Options: 一个闭包函数，接受一个 Worker 实例作为参数，并将其上下文设置为指定的值。
func WithContext(value context.Context) Options {
	return func(w *Worker) {
		if value == nil {
			value = context.Background()
		}

		// 使用 WithCancel 方法创建一个新的上下文对象，并将取消函数赋值给 w.cancel
		w.ctx, w.cancel = context.WithCancel(value)
	}
}
