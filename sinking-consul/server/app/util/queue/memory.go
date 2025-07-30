package queue

import (
	"errors"
)

type MemoryQueue struct {
	taskChan chan string
}

func (r *MemoryQueue) ListenReceiveMsgDo(topic string, receiveDo func(mqMsg MqMsg)) error {
	go func(taskChan chan string, topic string) {
		for value := range taskChan {
			receiveDo(MqMsg{
				RunType: ReceiveMsg,
				Topic:   topic,
				Body:    value,
			})
		}
	}(r.taskChan, topic)
	return nil
}

func (r *MemoryQueue) SendMsg(topic string, body string) (mqMsg MqMsg, err error) {
	if r.taskChan == nil {
		return mqMsg, errors.New("未初始化")
	}
	r.taskChan <- body
	return MqMsg{
		RunType: SendMsg,
		Topic:   topic,
		Body:    body,
	}, nil
}

// RegisterMemoryQueueClient 注册queue实例
func RegisterMemoryQueueClient(num int) (Queue, error) {
	return &MemoryQueue{
		taskChan: make(chan string, num),
	}, nil
}
