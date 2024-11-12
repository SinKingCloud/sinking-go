package route

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/http/controller/admin"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/http/middleware"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

func InitAdminRouter(route *sinking_web.Engine) {
	loadAdminIndexRoute(route)
	apiGroup := route.Group("/admin")
	apiGroup.Use(middleware.AdminAuth())
	{
		loadAdminClusterRoute(apiGroup)
		loadAdminServiceRoute(apiGroup)
	}
}

func loadAdminIndexRoute(route *sinking_web.Engine) {
	apiGroup := route.Group("/index")
	apiGroup.ANY("/login", admin.UserLogin)
}

func loadAdminClusterRoute(route *sinking_web.RouterGroup) {
	apiGroup := route.Group("/cluster")
	apiGroup.ANY("/list", admin.ClusterList)
}

func loadAdminServiceRoute(route *sinking_web.RouterGroup) {
	apiGroup := route.Group("/service")
	apiGroup.ANY("/list", admin.ServiceList)
}
