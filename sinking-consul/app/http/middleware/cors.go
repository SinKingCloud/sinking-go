package middleware

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"net/http"
)

// Cors 跨域
func Cors() sinking_web.HandlerFunc {
	return func(c *sinking_web.Context) {
		c.SetHeader("Access-Control-Allow-Origin", "*")
		c.SetHeader("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.SetHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.SetHeader("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.SetHeader("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.SetStatus(http.StatusNoContent)
			c.Abort()
		}
		c.Next()
	}
}
