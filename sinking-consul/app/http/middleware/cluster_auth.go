package middleware

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// ClusterAuth 集群通信鉴权
func ClusterAuth() sinking_web.HandlerFunc {
	return func(c *sinking_web.Context) {
		name := setting.GetConfig().GetString("servers.token-name")
		token := setting.GetConfig().GetString("servers.token")
		requestToken := c.Request.Header.Get(name)
		if requestToken != token {
			response.TokenError(c, "token验证失败！", nil)
			c.Abort()
		}
		c.Next()
	}
}
