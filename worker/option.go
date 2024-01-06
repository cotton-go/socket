package worker

import (
	"context"

	"worker/cache"
	"worker/codec"
	"worker/snowflake"
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
			value = &cache.Memory{}
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

func WithContext(value context.Context) Options {
	return func(w *Worker) {
		if value == nil {
			value = context.Background()
		}

		w.ctx, w.cancel = context.WithCancel(value)
	}
}
