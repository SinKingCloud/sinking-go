package task

import (
	"server/app/task/save"
	"server/app/task/sync"
)

func Init() {
	sync.Init()
	save.Init()
}
