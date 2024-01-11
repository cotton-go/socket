package event

import (
	"reflect"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

// Event 结构体定义了一个事件，包含主题和数据两个字段
type Event struct {
	Topic string `json:"topic"` // 事件主题，字符串类型
	Data  any    `json:"data"`  // 事件数据，任意类型
}

func (e *Event) Scan(value any) error {
	to := reflect.TypeOf(e.Data)
	if to.Kind() == reflect.Pointer {
		to = to.Elem()
	}

	switch to.Kind() {
	case reflect.Struct, reflect.Map:
		body, err := sonic.Marshal(e.Data)
		if err != nil {
			return errors.Wrap(err, "failed to parse json data")
		}

		return sonic.Unmarshal(body, value)
	default:
		// 直接赋值给指针
		reflect.ValueOf(value).Elem().Set(reflect.ValueOf(e.Data))
		return nil
	}
}
