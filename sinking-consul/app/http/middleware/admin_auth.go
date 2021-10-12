package middleware

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/constant/cachePrefix"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/model"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/cache"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/jwt"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"strconv"
	"strings"
	"time"
)

// AdminAuth 用户鉴权
func AdminAuth() sinking_web.HandlerFunc {
	return func(c *sinking_web.Context) {
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			response.NotLogin(c, "登陆超时，请重新登录(-1)", nil)
			c.Abort()
			return
		}
		checkToken := strings.SplitN(tokenHeader, " ", 2)
		if len(checkToken) != 2 && checkToken[0] != "Bearer" {
			response.NotLogin(c, "登陆超时，请重新登录(-2)", nil)
			c.Abort()
			return
		}
		key := jwt.CheckToken(checkToken[1])
		if key == nil || time.Now().Unix() > key.ExpiresAt {
			response.NotLogin(c, "登陆超时，请重新登录(-3)", nil)
			c.Abort()
			return
		}
		//判断权限
		if key.User.RoleId != 0 {
			data := cache.Remember(cachePrefix.Role+strconv.FormatInt(key.User.RoleId, 10), func() interface{} {
				var info model.Role
				model.Db.Where("id=?", key.User.RoleId).First(&info)
				return info
			}, 60*time.Second)
			role := data.(model.Role)
			if !strings.Contains(role.Auths, c.Request.RequestURI) {
				response.NotAllow(c, "权限不足", nil)
				c.Abort()
				return
			}
		}
		c.Set(cachePrefix.User, key.User)
		c.Next()
	}
}
