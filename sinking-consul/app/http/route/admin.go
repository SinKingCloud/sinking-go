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
	apiGroup.POST("/login", admin.UserLogin)
}

func loadAdminClusterRoute(route *sinking_web.RouterGroup) {
	apiGroup := route.Group("/cluster")
	apiGroup.POST("/list", admin.ClusterList)
}

func loadAdminServiceRoute(route *sinking_web.RouterGroup) {
	apiGroup := route.Group("/service")
	apiGroup.POST("/list", admin.ServiceList)
}
