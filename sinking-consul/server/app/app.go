package app

import (
	"server/app/command"
	"server/app/http/route"
	"server/app/service"
)

func Run() {
	service.Init()
	command.Init()
	route.Init()
}
