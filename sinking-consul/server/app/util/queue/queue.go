package queue

type Queue interface {
	SendMsg(topic string, body string) (mqMsg MqMsg, err error)
	ListenReceiveMsgDo(topic string, receiveDo func(mqMsg MqMsg)) (err error)
}

const (
	_ = iota
	SendMsg
	ReceiveMsg
)

type MqMsg struct {
	RunType int    `json:"run_type"`
	Topic   string `json:"topic"`
	Body    string `json:"body"`
}
