package sync

import (
	"encoding/json"
	"runtime"
	"server/app/util/queue"
)

const TopicName = "log"

var Instance *queue.Client

type Task struct {
	Type          Type   //任务类型
	RemoteAddress string //远程地址
}

func Register() *queue.Client {
	ins, err := queue.RegisterMemory(TopicName, 1, runtime.NumCPU(), func(param string) {
		var task *Task
		err := json.Unmarshal([]byte(param), &task)
		if err == nil {
			exec(task)
		}
	})
	if err != nil {
		panic(err)
		return nil
	}
	Instance = ins
	return ins
}
