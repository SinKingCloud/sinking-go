package sync

import (
	"runtime"
	"server/app/util/queue"
)

const TopicName = "sync"

type Task struct {
	Type          Type   //任务类型
	RemoteAddress string //远程地址
}

func Register() *queue.UnicastClient[*Task] {
	ins, err := queue.RegisterMemoryUnicast[*Task](TopicName, 1, runtime.NumCPU(), func(param *Task) error {
		exec(param)
		return nil
	})
	if err != nil {
		panic(err)
		return nil
	}
	return ins
}
