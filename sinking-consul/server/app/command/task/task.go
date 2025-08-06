package task

import (
	"server/app/command/task/save"
	"server/app/command/task/sync"
)

func Init() {
	sync.Init()
	save.Init()
}
