package middleware

import (
	"net/http"
	"server/app/util/server"
	"strings"
)

// Cors 跨域
func Cors(c *server.Context) {
	c.SetHeader("Access-Control-Allow-Origin", "*")
	c.SetHeader("Access-Control-Allow-Headers", "*")
	c.SetHeader("Access-Control-Allow-Methods", "*")
	c.SetHeader("Access-Control-Expose-Headers", "*")
	c.SetHeader("Access-Control-Allow-Credentials", "true")
	if strings.ToLower(c.Request.Method) == "options" {
		c.SetStatus(http.StatusNoContent)
		c.Abort()
	} else {
		c.Next()
	}
}
