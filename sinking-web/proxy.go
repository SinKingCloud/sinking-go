package sinking_web

import (
	"net/http/httputil"
	"net/url"
	"strings"
)

// Proxy 反向代理
func (c *Context) Proxy(uri string) {
	target, _ := url.Parse(uri)
	c.Request.Host = c.Request.URL.Host
	c.Request.URL.Path = strings.Replace(c.Request.URL.Path, "/proxy/", "/", 1)
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(c.Writer, c.Request)
}
