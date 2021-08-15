package sinking_web

import "net/http"

type Route interface {
	NotFound(c *Context)
}

func NotFound(c *Context) {
	c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
}
