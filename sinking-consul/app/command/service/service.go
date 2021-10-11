package service

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/logs"
)

func registerService() {
	logs.Println("开始启动注册服务....")
}

func registerConfig() {
	logs.Println("开始启动配置同步....")
}

func Run() {
	go registerService()
	go registerConfig()
}
