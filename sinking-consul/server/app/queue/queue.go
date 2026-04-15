package queue

import (
	"server/app/queue/sync"
	"server/app/util/queue"
)

var (
	Sync *queue.UnicastClient[*sync.Task]
)

func Init() {
	Sync = sync.Register()
}
