package service

import (
	"server/app/service/auth"
	"server/app/service/cluster"
)

// service实例
var (
	Auth    = auth.GetIns()
	Cluster = cluster.GetIns()
)

// Init 初始化服务
func Init() {
	Cluster.Init() // 初始化集群服务
}

// Enum 枚举信息
var Enum = map[string]interface{}{
	"cluster": map[string]interface{}{
		"online_status": Cluster.OnlineStatus(), //在线状态
		"status":        Cluster.Status(),       //集群状态
	},
}
