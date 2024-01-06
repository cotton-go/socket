package event

type Event struct {
	Topic string `json:"topic"`
	Data  any    `json:"data"`
}
