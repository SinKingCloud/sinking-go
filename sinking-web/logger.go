package sinking_web

import (
	"log"
	"os"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v [%d]", c.StatusCode, c.Request.RequestURI, time.Since(t), os.Getppid())
	}
}
