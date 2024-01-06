package worker

import (
	"context"

	"worker/pkg/cache"
	"worker/pkg/codec"
	"worker/pkg/connection"
	"worker/pkg/event"
	"worker/pkg/snowflake"
)

type Options func(w *Worker)

func WithID(value int64) Options {
	return func(w *Worker) {
		if value < 1 {
			value = snowflake.Next()
		}
		w.id = value
	}
}

func WithCache(value cache.ICache) Options {
	return func(w *Worker) {
		if value == nil {
			value = cache.NewMemory()
		}

		w.cache = value
	}
}

func WithCodec(value codec.ICodec) Options {
	return func(w *Worker) {
		if value == nil {
			value = &codec.Default{}
		}

		w.codec = value
	}
}

func WithHandle(value connection.EventHandle) Options {
	return func(w *Worker) {
		if value == nil {
			value = func(c *connection.Connection, e event.Event) {}
		}

		w.handle = value
	}
}

// func WithEvent() Options {
// 	return func(w *Worker) {
// 		w.events = make(map[string]connection.EventHandle)
// 	}
// }

func WithContext(value context.Context) Options {
	return func(w *Worker) {
		if value == nil {
			value = context.Background()
		}

		w.ctx, w.cancel = context.WithCancel(value)
	}
}
