package service

import (
	repositoryCluster "server/app/repository/cluster"
	repositoryConfig "server/app/repository/config"
	log2 "server/app/repository/log"
	repositoryNode "server/app/repository/node"
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
	Cluster cluster.Service
	Node    node.Service
	Config  config.Service
)

// Init 初始化服务
func Init() {
	conf := util.Conf
	cache := util.Cache
	database := util.Database
	repositoryLog := log2.NewRepository(database)
	clusterRepository := repositoryCluster.NewRepository(database)
	nodeRepository := repositoryNode.NewRepository(database)
	configRepository := repositoryConfig.NewRepository(database)

	Setting = setting.NewService(conf)
	Auth = auth.NewService(Setting, cache)
	Log = log.NewService(repositoryLog)
	Config = config.NewService(configRepository, cache)
	Node = node.NewService(nodeRepository, cache)
	Cluster = cluster.NewService(clusterRepository, cache, Setting, Node, Config)
}
