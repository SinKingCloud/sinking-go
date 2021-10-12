package route

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/http/controller/admin"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/http/middleware"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

func InitAdminRouter(route *sinking_web.Engine) {
	apiGroup := route.Group("/admin")
	apiGroup.Use(middleware.AdminAuth())
	{
		loadAdminIndexRoute(apiGroup)
		loadAdminClusterRoute(apiGroup)
		loadAdminServiceRoute(apiGroup)
	}
}

func loadAdminIndexRoute(route *sinking_web.RouterGroup) {
	apiGroup := route.Group("/index")
	apiGroup.GET("/login", admin.UserLogin)
}

func loadAdminClusterRoute(route *sinking_web.RouterGroup) {
	apiGroup := route.Group("/cluster")
	apiGroup.GET("/list", admin.ClusterList)
}

func loadAdminServiceRoute(route *sinking_web.RouterGroup) {
	apiGroup := route.Group("/service")
	apiGroup.GET("/list", admin.ServiceList)
}
