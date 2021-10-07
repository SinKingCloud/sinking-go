package sinking_web

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Proxy 反向代理
func (c *Context) Proxy(uri string) {
	Try(func() {
		test := c.Writer
		target, _ := url.Parse(uri)
		c.Request.Host = c.Request.URL.Host
		c.Request.URL.Path = strings.Replace(c.Request.URL.Path, "/proxy/", "/", 1)
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(test, c.Request)
	}, func(err interface{}) {
		//c.StatusCode = 200
		fmt.Println(err)
	})
}
