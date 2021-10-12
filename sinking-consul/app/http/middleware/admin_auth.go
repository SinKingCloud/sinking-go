package middleware

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/constant/cachePrefix"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/jwt"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"strings"
	"time"
)

// AdminAuth 用户鉴权
func AdminAuth() sinking_web.HandlerFunc {
	return func(c *sinking_web.Context) {
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			response.Return(c, 403, "登陆超时，请重新登录", nil)
			c.Abort()
			return
		}
		checkToken := strings.SplitN(tokenHeader, " ", 2)
		if len(checkToken) != 2 && checkToken[0] != "Bearer" {
			response.Return(c, 403, "登陆超时，请重新登录", nil)
			c.Abort()
			return
		}
		key := jwt.CheckToken(checkToken[1])
		if key == nil || time.Now().Unix() > key.ExpiresAt {
			response.Return(c, 403, "登陆超时，请重新登录", nil)
			c.Abort()
			return
		}
		c.Set(cachePrefix.User, key.User)
		c.Next()
	}
}
