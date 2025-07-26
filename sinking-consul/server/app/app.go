package app

import (
	"server/app/command"
	"server/app/http/route"
)

func Run() {
	command.Init()
	route.Init()
}
