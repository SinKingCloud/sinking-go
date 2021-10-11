package middleware

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// ApiAuth 通信鉴权
func ApiAuth() sinking_web.HandlerFunc {
	return func(c *sinking_web.Context) {
		name := setting.GetSystemConfig().Servers.TokenName
		token := setting.GetSystemConfig().Servers.Token
		requestToken := c.Request.Header.Get(name)
		if requestToken != token {
			response.TokenError(c, "token验证失败！", nil)
			c.Abort()
		}
		c.Next()
	}
}
