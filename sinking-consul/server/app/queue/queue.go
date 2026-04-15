package queue

import (
	"server/app/queue/sync"
	"server/app/util/queue"
)

var (
	Sync *queue.Client
)

func Init() {
	Sync = sync.Register()
}
