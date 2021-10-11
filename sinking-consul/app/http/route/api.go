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
		loadClusterRoute(apiGroup)
	}
}

func loadClusterRoute(route *sinking_web.RouterGroup) {
	apiGroup := route.Group("/cluster")
	apiGroup.GET("/register", api.ClusterRegister)
	apiGroup.GET("/get", api.ClusterGet)
}
