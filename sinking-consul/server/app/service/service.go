package service

import (
	"server/app/service/auth"
	"server/app/service/config"
	"server/app/service/log"
)

// service实例
var (
	Config = config.GetIns()
	Auth   = auth.GetIns()
	Log    = log.GetIns()
)

// Enum 枚举信息
var Enum = map[string]interface{}{
	"log": map[string]interface{}{
		"type": Log.Types(), //日志类型
	},
}
