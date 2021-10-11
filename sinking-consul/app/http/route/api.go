package route

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/http/controller/api"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/http/middleware"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

func InitApiRouter(route *sinking_web.Engine) {
	loadClusterRoute(route)
}

func loadClusterRoute(route *sinking_web.Engine) {
	apiGroup := route.Group("/api/cluster")
	apiGroup.Use(middleware.ClusterAuth())
	{
		apiGroup.GET("/register", api.ClusterRegister)
	}
}
