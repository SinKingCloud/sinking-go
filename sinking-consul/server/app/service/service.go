package service

import (
	log2 "server/app/repository/log"
	"server/app/service/auth"
	"server/app/service/cluster"
	"server/app/service/config"
	"server/app/service/log"
	"server/app/service/node"
	"server/app/service/setting"
	"server/app/util"
)

// service实例
var (
	Log     log.Service
	Auth    auth.Service
	Setting setting.Service
	Cluster = cluster.GetIns()
	Node    = node.GetIns()
	Config  = config.GetIns()
)

// Init 初始化服务
func Init() {
	conf := util.Conf
	cache := util.Cache
	database := util.Database
	repositoryLog := log2.NewRepository(database)

	Setting = setting.NewService(conf)
	Auth = auth.NewService(Setting, cache)
	Log = log.NewService(repositoryLog)

	Cluster.Init() // 初始化集群服务
	Config.Init()  // 初始化配置服务
	Node.Init()    // 初始化节点服务
}
