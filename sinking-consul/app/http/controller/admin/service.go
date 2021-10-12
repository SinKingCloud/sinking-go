package admin

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// ServiceList 服务列表
func ServiceList(s *sinking_web.Context) {
	response.Success(s, "获取服务列表成功", service.Services)
}
