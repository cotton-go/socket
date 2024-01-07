package subscribe

import "sync"

// 定义一个订阅者接口
type Subscriber interface {
	Notify(message string)
}

// 定义发布者结构体
type Publisher struct {
	subscribers []Subscriber
	mutex       sync.Mutex
}

// 添加订阅者
func (p *Publisher) AddSubscriber(subscriber Subscriber) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.subscribers = append(p.subscribers, subscriber)
}

// 移除订阅者
func (p *Publisher) RemoveSubscriber(subscriber Subscriber) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for i, sub := range p.subscribers {
		if sub == subscriber {
			p.subscribers = append(p.subscribers[:i], p.subscribers[i+1:]...)
			break
		}
	}
}

// 发送消息给所有订阅者
func (p *Publisher) SendMessage(message string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, sub := range p.subscribers {
		sub.Notify(message)
	}
}
