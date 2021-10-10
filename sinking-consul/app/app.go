package app

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/command/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/http/route"
)

func Run() {
	//启动服务
	service.Run()
	//启动应用
	route.Init()
}
