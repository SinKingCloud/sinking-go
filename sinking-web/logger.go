package sinking_web

import (
	"log"
	"os"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		c.Next()
		log.Printf("[STATUS:%d] %s in %v [%d]", c.StatusCode, c.Request.RequestURI, time.Since(t), os.Getppid())
	}
}
