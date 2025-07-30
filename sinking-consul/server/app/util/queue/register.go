package queue

import (
	"encoding/json"
	"errors"
	"server/app/util/job"
)

type Client struct {
	topic             Queue
	name              string
	consumer          chan interface{}
	consumerThreadNum int
	consumerFunc      func(param string)
	producer          chan interface{}
	producerThreadNum int
}

// Try 错误捕获实现
func (l *Client) Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

func (l *Client) run() error {
	if l.consumer == nil {
		l.consumer = make(chan interface{}, 10000)
	}
	if l.producer == nil {
		l.producer = make(chan interface{}, 10000)
	}
	go (&job.Task{Producer: func(channel chan interface{}) {
		l.consumer = channel
	}, Consumer: func(param interface{}) {
		l.Try(func() {
			l.consumerFunc(param.(string))
		}, nil)
	}, Thread: l.consumerThreadNum}).Run()
	go (&job.Task{Producer: func(channel chan interface{}) {
		l.producer = channel
	}, Consumer: func(param interface{}) {
		l.Try(func() {
			marshal, e := json.Marshal(param)
			if e == nil {
				_, e = l.topic.SendMsg(l.name, string(marshal))
			}
		}, nil)
	}, Thread: l.producerThreadNum}).Run()
	err := l.topic.ListenReceiveMsgDo(l.name, func(mqMsg MqMsg) {
		l.consumer <- mqMsg.Body
	})
	if err != nil {
		return err
	}
	return nil
}

// SendTask 发送任务
func (l *Client) SendTask(param interface{}) {
	l.producer <- param
}

// register 注册消息队列实例
// topicName 队列名称
// producerThreadNum 发送者协程数量
// consumerThreadNum 消费者协程数量
func register(queue Queue, topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string)) (*Client, error) {
	if topicName == "" {
		return nil, errors.New("队列名称不能为空")
	}
	if producerThreadNum <= 0 {
		return nil, errors.New("发送者协程数量必须大于0")
	}
	if consumerThreadNum <= 0 {
		return nil, errors.New("消费者协程数量必须大于0")
	}
	ins := &Client{
		topic:             queue,
		name:              topicName,
		producerThreadNum: producerThreadNum,
		consumerThreadNum: consumerThreadNum,
		consumerFunc:      consumerFunc,
	}
	err := ins.run()
	if err != nil {
		return nil, errors.New("启动" + topicName + "消费实例失败")
	}
	return ins, nil
}

// RegisterMemory 注册redis消息队列实例
// topicName 队列名称
// producerThreadNum 发送者协程数量
// consumerThreadNum 消费者协程数量
func RegisterMemory(topicName string, producerThreadNum int, consumerThreadNum int, consumerFunc func(param string)) (*Client, error) {
	topic, err := RegisterMemoryQueueClient(100) //注册消费实例
	if err != nil {
		return nil, errors.New("注册" + topicName + "消费实例失败")
	}
	return register(topic, topicName, producerThreadNum, consumerThreadNum, consumerFunc)
}
