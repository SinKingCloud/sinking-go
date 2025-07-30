package queue

import (
	"server/app/command/queue/sync"
	"server/app/util/queue"
)

var (
	Sync *queue.Client
)

func Init() {
	Sync = sync.Register()
}
