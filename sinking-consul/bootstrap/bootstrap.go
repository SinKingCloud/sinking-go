package bootstrap

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/bootstrap/config"
	"github.com/SinKingCloud/sinking-go/sinking-consul/bootstrap/database"
	"github.com/SinKingCloud/sinking-go/sinking-consul/bootstrap/log"
)

func Init() {
	//初始化系统配置
	config.Init()
	//初始化日志
	log.Init()
	//初始化数据库
	database.Init()
}
