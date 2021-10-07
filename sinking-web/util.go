package sinking_web

import (
	"net"
	"strings"
)

// ClientIP 获取用户ip(是否使用代理)
func (c *Context) ClientIP(useProxy bool) string {
	if useProxy {
		xForwardedFor := c.Request.Header.Get("X-Forwarded-For")
		ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
		if ip != "" {
			return ip
		}
		ip = strings.TrimSpace(c.Request.Header.Get("X-Real-Ip"))
		if ip != "" {
			return ip
		}
	} else {
		if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
			return ip
		}
	}
	return ""
}

// Try 错误捕获实现
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}
