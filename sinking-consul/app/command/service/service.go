package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/logs"
)

func start() {
	logs.Println("开始注册服务....")
	logs.Println("注册服务成功....")
}

func Run() {
	go start()
}
