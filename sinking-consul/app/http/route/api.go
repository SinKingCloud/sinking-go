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
	apiGroup.POST("/register", api.ClusterRegister)
}

func loadApiServiceRoute(route *sinking_web.RouterGroup) {
	apiGroup := route.Group("/service")
	apiGroup.POST("/register", api.ServiceRegister)
	apiGroup.POST("/status", api.ServiceStatus)
	apiGroup.POST("/list", api.ServiceList)
}

func loadApiConfigRoute(route *sinking_web.RouterGroup) {
	apiGroup := route.Group("/config")
	apiGroup.POST("/list", api.ConfigList)
}
