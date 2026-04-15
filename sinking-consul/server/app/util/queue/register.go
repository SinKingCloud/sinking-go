package queue

import (
	"encoding/json"
	"errors"
)

// UnicastClient 单播消息队列客户端（泛型版）
type UnicastClient[T any] struct {
	topic        Queue
	name         string
	consumerFunc func(param T) error
	taskChan     chan MqMsg
}

// Try 错误捕获实现
func (l *UnicastClient[T]) Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

// SendUnicastTask 发送单播任务（自动序列化）
func (l *UnicastClient[T]) SendUnicastTask(param T) {
	marshal, e := json.Marshal(param)
	if e == nil {
		_, _ = l.topic.SendUnicast(l.name, string(marshal))
	}
}

// Stop 停止客户端
func (l *UnicastClient[T]) Stop() error {
	if l.taskChan != nil {
		close(l.taskChan)
	}
	if l.topic != nil {
		l.topic.Close()
	}
	return nil
}

// RegisterUnicast 注册单播消息队列实例
func RegisterUnicast[T any](queue Queue, topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param T) error) (*UnicastClient[T], error) {
	if topicName == "" {
		return nil, errors.New("队列名称不能为空")
	}
	if producerThreadNum <= 0 {
		return nil, errors.New("发送者协程数量必须大于0")
	}
	if consumerThreadNum <= 0 {
		return nil, errors.New("消费者协程数量必须大于0")
	}

	ins := &UnicastClient[T]{
		topic:        queue,
		name:         topicName,
		consumerFunc: consumerFunc,
		taskChan:     make(chan MqMsg, consumerThreadNum*100),
	}

	// 预创建 worker pool，消费时自动反序列化
	for i := 0; i < consumerThreadNum; i++ {
		go func() {
			for mqMsg := range ins.taskChan {
				ins.Try(func() {
					var param T
					if err := json.Unmarshal([]byte(mqMsg.Body), &param); err == nil {
						_ = ins.consumerFunc(param)
					}
				}, func(e interface{}) {})
			}
		}()
	}

	// 监听单播消息
	err := queue.ListenReceiveUnicast(topicName, func(mqMsg MqMsg) bool {
		ins.taskChan <- mqMsg
		return true
	})
	if err != nil {
		return nil, errors.New("启动" + topicName + "单播消费实例失败: " + err.Error())
	}

	return ins, nil
}

// BroadcastClient 广播消息队列客户端（泛型版）
type BroadcastClient[T any] struct {
	topic        Queue
	name         string
	consumerFunc func(param T) error
	taskChan     chan MqMsg
}

// Try 错误捕获实现
func (l *BroadcastClient[T]) Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

// SendBroadcastTask 发送广播任务（自动序列化）
func (l *BroadcastClient[T]) SendBroadcastTask(param T) {
	marshal, e := json.Marshal(param)
	if e == nil {
		_, _ = l.topic.SendBroadcast(l.name, string(marshal))
	}
}

// Stop 停止客户端
func (l *BroadcastClient[T]) Stop() error {
	if l.taskChan != nil {
		close(l.taskChan)
	}
	if l.topic != nil {
		l.topic.Close()
	}
	return nil
}

// RegisterBroadcast 注册广播消息队列实例
func RegisterBroadcast[T any](queue Queue, topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param T) error) (*BroadcastClient[T], error) {
	if topicName == "" {
		return nil, errors.New("队列名称不能为空")
	}
	if producerThreadNum <= 0 {
		return nil, errors.New("发送者协程数量必须大于0")
	}
	if consumerThreadNum <= 0 {
		return nil, errors.New("消费者协程数量必须大于0")
	}

	ins := &BroadcastClient[T]{
		topic:        queue,
		name:         topicName,
		consumerFunc: consumerFunc,
		taskChan:     make(chan MqMsg, consumerThreadNum*100),
	}

	// 预创建 worker pool，消费时自动反序列化
	for i := 0; i < consumerThreadNum; i++ {
		go func() {
			for mqMsg := range ins.taskChan {
				ins.Try(func() {
					var param T
					if err := json.Unmarshal([]byte(mqMsg.Body), &param); err == nil {
						_ = ins.consumerFunc(param)
					}
				}, func(e interface{}) {})
			}
		}()
	}

	// 监听广播消息
	err := queue.ListenReceiveBroadcast(topicName, func(mqMsg MqMsg) bool {
		ins.taskChan <- mqMsg
		return true
	})
	if err != nil {
		return nil, errors.New("启动" + topicName + "广播消费实例失败: " + err.Error())
	}

	return ins, nil
}

// RegisterMemoryUnicast 注册内存单播消息队列实例
func RegisterMemoryUnicast[T any](topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param T) error) (*UnicastClient[T], error) {
	topic, err := RegisterMemoryQueueClient()
	if err != nil {
		return nil, errors.New("注册" + topicName + "单播消费实例失败")
	}
	return RegisterUnicast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}

// RegisterMemoryBroadcast 注册内存广播消息队列实例
func RegisterMemoryBroadcast[T any](topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param T) error) (*BroadcastClient[T], error) {
	topic, err := RegisterMemoryQueueClient()
	if err != nil {
		return nil, errors.New("注册" + topicName + "广播消费实例失败")
	}
	return RegisterBroadcast(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}
