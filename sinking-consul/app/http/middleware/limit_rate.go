package middleware

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// LimitRateMiddleware 限流
func LimitRateMiddleware() sinking_web.HandlerFunc {
	return func(c *sinking_web.Context) {
		c.Next()
	}
}
