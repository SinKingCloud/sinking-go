package app

import (
	"server/app/command"
	"server/app/http/route"
	"server/app/service"
)

func Run() {
	command.Init()
	service.Init()
	route.Init()
}
