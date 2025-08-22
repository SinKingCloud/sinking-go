package service

import (
	"server/app/service/auth"
	"server/app/service/cluster"
	"server/app/service/config"
	"server/app/service/node"
)

// service实例
var (
	Auth    = auth.GetIns()
	Cluster = cluster.GetIns()
	Node    = node.GetIns()
	Config  = config.GetIns()
)

// Init 初始化服务
func Init() {
	Cluster.Init() // 初始化集群服务
	Config.Init()  // 初始化配置服务
	Node.Init()    // 初始化节点服务
}

// Enum 枚举信息
var Enum = map[string]interface{}{
	"cluster": map[string]interface{}{
		"status": Cluster.Status(), //在线状态
	},
	"node": map[string]interface{}{
		"online_status": Node.OnlineStatus(), //在线状态
		"status":        Node.Status(),       //集群状态
	},
	"config": map[string]interface{}{
		"type":   Config.Types(),  //配置类型
		"status": Config.Status(), //是否启用
	},
}
