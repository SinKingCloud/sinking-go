package command

import (
	"server/app/command/queue"
	"server/app/command/task"
)

func Init() {
	queue.Init()
	task.Init()
}
