package admin

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// ClusterList 集群列表
func ClusterList(s *sinking_web.Context) {
	response.Success(s, "获取集群列表成功", service.Clusters)
}
