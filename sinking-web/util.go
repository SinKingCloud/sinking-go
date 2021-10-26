package sinking_web

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
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

// Proxy 反向代理
func (c *Context) Proxy(pattern string, uri string, filter func(r *http.Request) *http.Request) {
	Try(func() {
		target, err := url.Parse(uri)
		if err != nil {
			c.JSON(500, "url format error.")
			return
		}
		c.StatusCode = 200
		c.Request.Host = c.Request.URL.Host
		c.Request.URL.Path = strings.Replace(c.Request.URL.Path, pattern[0:strings.Index(pattern, "*")], "/", 1)
		filter(c.Request)
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(c.Writer, c.Request)
	}, func(err interface{}) {
		c.StatusCode = 500
	})
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
