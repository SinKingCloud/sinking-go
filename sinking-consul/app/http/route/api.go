package route

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/http/controller/api"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/http/middleware"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

func InitApiRouter(route *sinking_web.Engine) {
	apiGroup := route.Group("/api")
	apiGroup.Use(middleware.ApiAuth())
	{
		loadApiClusterRoute(apiGroup)
		loadApiServiceRoute(apiGroup)
		loadApiConfigRoute(apiGroup)
	}
}

func loadApiClusterRoute(route *sinking_web.RouterGroup) {
	apiGroup := route.Group("/cluster")
	apiGroup.ANY("/list", api.ClusterList)
	apiGroup.ANY("/register", api.ClusterRegister)
	apiGroup.ANY("/services", api.ClusterServiceList)
	apiGroup.ANY("/configs", api.ClusterConfigList)
}

func loadApiServiceRoute(route *sinking_web.RouterGroup) {
	apiGroup := route.Group("/service")
	apiGroup.ANY("/register", api.ServiceRegister)
	apiGroup.ANY("/status", api.ServiceStatus)
	apiGroup.ANY("/list", api.ServiceList)
	apiGroup.ANY("/all_list", api.GetProjectAllServiceList)
}

func loadApiConfigRoute(route *sinking_web.RouterGroup) {
	apiGroup := route.Group("/config")
	apiGroup.ANY("/list", api.ConfigList)
}
