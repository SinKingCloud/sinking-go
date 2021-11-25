package sinking_web

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
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

// HttpProxy http反向代理
func (c *Context) HttpProxy(uri string, filter func(r *http.Request) *http.Request) {
	Try(func() {
		target, err := url.Parse(uri)
		if err != nil {
			c.JSON(500, "url format error.")
			return
		}
		c.StatusCode = 200
		filter(c.Request)
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(c.Writer, c.Request)
	}, func(err interface{}) {
		c.StatusCode = 500
	})
}

// WebSocketProxy WebSocketProxy反向代理
func (c *Context) WebSocketProxy(uri string, filter func(r *http.Request) *http.Request) {
	Try(func() {
		u, err := url.Parse(uri)
		if err != nil {
			return
		}
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return
		}
		if u.Scheme != "ws" && u.Scheme != "wss" {
			return
		}
		// 劫持连接
		hijacker, ok := c.Writer.(http.Hijacker)
		if !ok {
			return
		}
		conn, _, err := hijacker.Hijack()
		if err != nil {
			return
		}
		defer func(conn net.Conn) {
			err = conn.Close()
			if err != nil {
				return
			}
		}(conn)
		req := c.Request.Clone(context.TODO())
		req.URL.Path, req.URL.RawPath, req.RequestURI = u.Path, u.Path, u.Path
		req.Host = host
		filter(req)
		var remoteConn net.Conn
		switch u.Scheme {
		case "ws":
			remoteConn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		case "wss":
			remoteConn, err = tls.Dial("tcp", fmt.Sprintf("%s:%s", host, port), &tls.Config{InsecureSkipVerify: true})
		}
		if err != nil {
			_, _ = c.Writer.Write([]byte(err.Error()))
			return
		}
		defer func(remoteConn net.Conn) {
			err = remoteConn.Close()
			if err != nil {
				return
			}
		}(remoteConn)
		b, _ := httputil.DumpRequest(req, false)
		_, err = remoteConn.Write(b)
		if err != nil {
			return
		}
		errChan := make(chan error, 2)
		copyConn := func(a, b net.Conn) {
			_, err := io.Copy(a, b)
			errChan <- err
		}
		go copyConn(conn, remoteConn) // response
		go copyConn(remoteConn, conn) // request
		select {
		case err = <-errChan:
			if err != nil {
				log.Println(err)
			}
		}
	}, func(err interface{}) {
		c.StatusCode = 500
	})
}

// Proxy 通用反向代理
func (c *Context) Proxy(uri string, filter func(r *http.Request) *http.Request) {
	prefix := uri[0:2]
	if prefix == "ws" {
		c.WebSocketProxy(uri, filter)
	} else {
		c.HttpProxy(uri, filter)
	}
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
