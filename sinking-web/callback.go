package sinking_web

import (
	"log"
	"net/http"
)

type Route struct {
	NotFound func(*Context)
}

var route Route

func NotFoundHandle(c *Context) {
	if route.NotFound != nil {
		route.NotFound(c)
	} else {
		log.Println("未实现NotFound方法")
	}
	c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
}
