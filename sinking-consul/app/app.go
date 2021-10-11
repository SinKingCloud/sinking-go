package app

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/command"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/http/route"
)

func Run() {
	//启动服务
	command.Run()
	//启动应用
	route.Init()
}
