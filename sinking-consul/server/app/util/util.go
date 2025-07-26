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
