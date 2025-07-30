package util

import (
	"github.com/spf13/viper"
	"log"
	"server/app/constant"
	"server/app/util/cache"
	"server/app/util/database"
)

var (
	Database *database.Database //数据库
	Conf     *viper.Viper       //文件配置
	Cache    *cache.Cache
	Log      *log.Logger
)

// IsDebug 是否debug模式
func IsDebug() bool {
	return Conf.GetString(constant.ServerMode) == "dev"
}

// ServerAddr 获取服务监听地址
func ServerAddr() (host string, port int) {
	host = Conf.GetString(constant.ServerHost)
	port = Conf.GetInt(constant.ServerPort)
	if host == "" {
		host = "0.0.0.0"
	}
	if port <= 0 {
		port = 5678
	}
	return host, port
}
