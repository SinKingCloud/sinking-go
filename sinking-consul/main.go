package main

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/bootstrap/config"
	"github.com/SinKingCloud/sinking-go/sinking-consul/bootstrap/database"
	"github.com/SinKingCloud/sinking-go/sinking-consul/bootstrap/log"
)

func main() {
	//加载系统配置
	config.Init()
	//链接elk日志
	log.Init()
	//链接数据库
	database.Init()
}
