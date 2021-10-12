package admin

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// ClusterList 集群列表
func ClusterList(s *sinking_web.Context) {
	var list []*service.Cluster
	for _, v := range service.Clusters {
		list = append(list, v)
	}
	response.Success(s, "获取集群列表成功", list)
}
