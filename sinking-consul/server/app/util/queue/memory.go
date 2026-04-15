package queue

import (
	"errors"
	"math/rand"
	"sync"
)

type MemoryQueue struct {
	unicastChans   map[string][]chan string // 单播消息通道 (key: topic, value: channel 列表)
	broadcastChans map[string][]chan string // 广播消息通道 (key: topic, value: channel 列表)
	closed         bool
	mu             sync.RWMutex
}

// ListenReceiveUnicast 监听接收单播消息
func (r *MemoryQueue) ListenReceiveUnicast(topic string, receiveDo func(mqMsg MqMsg) bool) error {
	r.mu.Lock()
	if r.unicastChans == nil {
		r.unicastChans = make(map[string][]chan string)
	}

	// 为这个监听器创建独立的 channel
	ch := make(chan string, 10000)
	r.unicastChans[topic] = append(r.unicastChans[topic], ch)
	r.mu.Unlock()

	go func(taskChan chan string, topic string) {
		for value := range taskChan {
			receiveDo(MqMsg{
				RunType: ReceiveUnicast,
				Topic:   topic,
				Body:    value,
			})
		}
	}(ch, topic)
	return nil
}

// ListenReceiveBroadcast 监听接收广播消息
func (r *MemoryQueue) ListenReceiveBroadcast(topic string, receiveDo func(mqMsg MqMsg) bool) error {
	r.mu.Lock()
	if r.broadcastChans == nil {
		r.broadcastChans = make(map[string][]chan string)
	}

	// 为这个监听器创建独立的 channel
	ch := make(chan string, 10000)
	r.broadcastChans[topic] = append(r.broadcastChans[topic], ch)
	r.mu.Unlock()

	go func(taskChan chan string, topic string) {
		for value := range taskChan {
			receiveDo(MqMsg{
				RunType: ReceiveBroadcast,
				Topic:   topic,
				Body:    value,
			})
		}
	}(ch, topic)
	return nil
}

// SendUnicast 发送单播消息（负载均衡，随机发送给一个订阅者）
func (r *MemoryQueue) SendUnicast(topic string, body string) (mqMsg MqMsg, err error) {
	r.mu.RLock()
	channels := r.unicastChans[topic]
	r.mu.RUnlock()

	if len(channels) == 0 {
		return mqMsg, errors.New("unicast topic not subscribed")
	}
	if r.closed {
		return mqMsg, errors.New("queue is closed")
	}

	// 随机选择一个 channel 发送
	ch := channels[rand.Intn(len(channels))]

	select {
	case ch <- body:
		return MqMsg{
			RunType: SendUnicast,
			Topic:   topic,
			Body:    body,
		}, nil
	default:
		return mqMsg, errors.New("queue is full or closed")
	}
}

// SendBroadcast 发送广播消息（所有订阅者都能收到）
func (r *MemoryQueue) SendBroadcast(topic string, body string) (mqMsg MqMsg, err error) {
	r.mu.RLock()
	channels := r.broadcastChans[topic]
	r.mu.RUnlock()

	if len(channels) == 0 {
		return mqMsg, errors.New("broadcast topic not subscribed")
	}
	if r.closed {
		return mqMsg, errors.New("queue is closed")
	}

	// 发送给所有订阅者
	for _, ch := range channels {
		select {
		case ch <- body:
		default:
			// 队列满时跳过
		}
	}

	return MqMsg{
		RunType: SendBroadcast,
		Topic:   topic,
		Body:    body,
	}, nil
}

// Close 关闭内存队列
func (r *MemoryQueue) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.closed {
		r.closed = true
		// 关闭所有单播通道
		for _, channels := range r.unicastChans {
			for _, ch := range channels {
				close(ch)
			}
		}
		// 关闭所有广播通道
		for _, channels := range r.broadcastChans {
			for _, ch := range channels {
				close(ch)
			}
		}
	}
}

// RegisterMemoryQueueClient 注册queue实例
func RegisterMemoryQueueClient() (Queue, error) {
	return &MemoryQueue{
		unicastChans:   make(map[string][]chan string),
		broadcastChans: make(map[string][]chan string),
		closed:         false,
	}, nil
}
