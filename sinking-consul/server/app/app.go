package app

import (
	"server/app/http/route"
	"server/app/queue"
	"server/app/service"
	"server/app/task"
)

func Run() {
	service.Init()
	queue.Init()
	task.Init()
	route.Init()
}
