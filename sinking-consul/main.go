package main

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app"
	"github.com/SinKingCloud/sinking-go/sinking-consul/bootstrap"
)

func main() {
	//初始化
	bootstrap.Init()
	//启动应用
	app.Run()
}
