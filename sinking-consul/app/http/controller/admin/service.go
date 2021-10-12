package admin

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// ServiceList 服务列表
func ServiceList(s *sinking_web.Context) {
	var list []*service.Service
	for _, v := range service.Services {
		list = append(list, v)
	}
	response.Success(s, "获取服务列表成功", list)
}
