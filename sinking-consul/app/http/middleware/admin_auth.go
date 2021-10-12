package middleware

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// AdminAuth 用户鉴权
func AdminAuth() sinking_web.HandlerFunc {
	return func(c *sinking_web.Context) {
		c.Next()
	}
}
