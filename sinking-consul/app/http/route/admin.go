package route

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/http/controller/admin"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

func InitAdminRouter(route *sinking_web.Engine) {
	apiGroup := route.Group("/admin")
	apiGroup.Use()
	{
		apiGroup.GET("/index", admin.Index)
	}
}
