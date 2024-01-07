package event

// Event 结构体定义了一个事件，包含主题和数据两个字段
type Event struct {
	Topic string `json:"topic"` // 事件主题，字符串类型
	Data  any    `json:"data"`  // 事件数据，任意类型
}
