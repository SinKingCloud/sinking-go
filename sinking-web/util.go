package sinking_web

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
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

// 代理poll
var proxyMap = &sync.Map{}

// HttpProxy http反向代理
func (c *Context) HttpProxy(uri string, logger *log.Logger, filter func(r *http.Request, w http.ResponseWriter, proxy *httputil.ReverseProxy), errorHandle func(http.ResponseWriter, *http.Request, error)) (err error) {
	Try(func() {
		var proxy *httputil.ReverseProxy
		if value, ok := proxyMap.Load(uri); ok {
			proxy = value.(*httputil.ReverseProxy)
		} else {
			target, e := url.Parse(uri)
			if e != nil {
				err = e
				return
			}
			tmp := httputil.NewSingleHostReverseProxy(target)
			dialer := &net.Dialer{
				Timeout:   time.Minute,
				KeepAlive: time.Minute,
			}
			tmp.Transport = &http.Transport{
				Proxy:             http.ProxyFromEnvironment,
				DialContext:       dialer.DialContext,
				ForceAttemptHTTP2: true,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}
			if logger != nil {
				tmp.ErrorLog = logger
			}
			if errorHandle != nil {
				tmp.ErrorHandler = errorHandle
			}
			proxyMap.Store(uri, tmp)
			proxy = tmp
		}
		if filter != nil {
			filter(c.Request, c.Writer, proxy)
		}
		proxy.ServeHTTP(c.Writer, c.Request)
	}, func(e interface{}) {
		c.StatusCode = 500
		if errMsg, ok := e.(string); ok {
			err = errors.New(errMsg)
		} else {
			err = errors.New("http proxy error")
		}
	})
	return err
}

// WebSocketProxy WebSocketProxy反向代理
func (c *Context) WebSocketProxy(uri string, logger *log.Logger, filter func(r *http.Request, w http.ResponseWriter), errorHandle func(http.ResponseWriter, *http.Request, error)) (err error) {
	Try(func() {
		u, e := url.Parse(uri)
		if e != nil {
			err = errors.New("url.Parse error")
			return
		}
		host, port, e := net.SplitHostPort(u.Host)
		if e != nil {
			err = errors.New("host and port must be valid")
			return
		}
		if u.Scheme != "ws" && u.Scheme != "wss" {
			err = errors.New("url scheme error")
			return
		}
		// 劫持连接
		hijacker, ok := c.Writer.(http.Hijacker)
		if !ok {
			err = errors.New("hijacker format error")
			return
		}
		conn, _, e := hijacker.Hijack()
		if e != nil {
			err = errors.New("hijacker error" + e.Error())
			return
		}
		defer func(conn net.Conn) {
			e2 := conn.Close()
			if e2 != nil {
				return
			}
		}(conn)
		req := c.Request.Clone(context.TODO())
		req.URL.Path, req.URL.RawPath, req.RequestURI = u.Path, u.Path, u.Path
		req.Host = host
		if filter != nil {
			filter(req, c.Writer)
		}
		var remoteConn net.Conn
		switch u.Scheme {
		case "ws":
			remoteConn, e = net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		case "wss":
			remoteConn, e = tls.Dial("tcp", fmt.Sprintf("%s:%s", host, port), &tls.Config{InsecureSkipVerify: true})
		}
		if e != nil {
			_, _ = c.Writer.Write([]byte(e.Error()))
			err = errors.New("remote connection failed")
			return
		}
		defer func(remoteConn net.Conn) {
			e3 := remoteConn.Close()
			if e3 != nil {
				return
			}
		}(remoteConn)
		b, e := httputil.DumpRequest(req, false)
		if e != nil {
			err = errors.New("http request failed")
			return
		}
		_, e = remoteConn.Write(b)
		if e != nil {
			err = errors.New("conn write failed")
			return
		}
		errChan := make(chan error, 2)
		copyConn := func(a, b net.Conn) {
			_, err2 := io.Copy(a, b)
			errChan <- err2
		}
		go copyConn(conn, remoteConn) // response
		go copyConn(remoteConn, conn) // request
		select {
		case e = <-errChan:
			if e != nil {
				err = e
			}
		}
	}, func(e interface{}) {
		c.StatusCode = 500
		if errMsg, ok := e.(string); ok {
			err = errors.New(errMsg)
		} else {
			err = errors.New("http proxy error")
		}
		if errorHandle != nil {
			errorHandle(c.Writer, c.Request, err)
		}
	})
	if logger != nil {
		logger.Println("websocket proxy error:", err)
	}
	return err
}

// Proxy 通用反向代理
func (c *Context) Proxy(uri string, logger *log.Logger, filter func(r *http.Request, w http.ResponseWriter, proxy *httputil.ReverseProxy), errorHandle func(http.ResponseWriter, *http.Request, error)) error {
	prefix := uri[0:2]
	if prefix == "ws" {
		fun := func(r *http.Request, w http.ResponseWriter) {
			filter(r, w, nil)
		}
		return c.WebSocketProxy(uri, logger, fun, errorHandle)
	} else {
		return c.HttpProxy(uri, logger, filter, errorHandle)
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
