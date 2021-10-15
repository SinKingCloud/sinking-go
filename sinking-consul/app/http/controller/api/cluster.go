package api

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/model"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// ClusterRegister 注册集群
func ClusterRegister(s *sinking_web.Context) {
	type register struct {
		Ip   string `form:"ip" json:"ip"`
		Port string `form:"port" json:"port"`
	}
	cluster := &register{}
	err := s.BindJson(&cluster)
	if err != nil || cluster.Ip == "" || cluster.Port == "" {
		response.Error(s, "参数不足", nil)
		return
	}
	service.ClustersRegister(cluster.Ip, cluster.Port)
	response.Success(s, "注册集群成功", nil)
}

// ClusterList 集群列表
func ClusterList(s *sinking_web.Context) {
	response.Success(s, "获取集群列表成功", service.ClustersList())
}

// ClusterServiceList 集群服务列表
func ClusterServiceList(s *sinking_web.Context) {
	list := service.GetAllServiceList()
	response.Success(s, "获取集群服务信息成功", list)
}

// ClusterConfigList 集群配置列表
func ClusterConfigList(s *sinking_web.Context) {
	list := (&model.Config{}).SelectAllCache()
	response.Success(s, "获取集群配置信息成功", list)
}
