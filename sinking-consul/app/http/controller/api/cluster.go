package api

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/encode"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"time"
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
	info := &service.Cluster{
		Hash:          encode.Md5Encode(cluster.Ip + ":" + cluster.Port),
		Ip:            cluster.Ip,
		Port:          cluster.Port,
		LastHeartTime: time.Now().Unix(),
		Status:        0,
	}
	service.ClustersLock.Lock()
	service.Clusters[info.Hash] = info
	service.ClustersLock.Unlock()
	response.Success(s, "注册集群成功", nil)
}
