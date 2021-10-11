package api

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/model"
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
	if cluster.Ip == "" || cluster.Port == "" {
		response.Error(s, "注册集群失败,参数不足", nil)
		return
	}
	info := &service.Cluster{
		Hash:          encode.Md5Encode(cluster.Ip + ":" + cluster.Port),
		Ip:            cluster.Ip,
		Port:          cluster.Port,
		LastHeartTime: model.DateTime(time.Now()),
		Status:        0,
	}
	service.Clusters[info.Hash] = *info
	response.Success(s, "注册集群成功", info)
}
