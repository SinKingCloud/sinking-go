package sinking_web

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				if !debug {
					c.Fail(http.StatusInternalServerError, MessageInternalServerError)
				} else {
					message := fmt.Sprintf("%s", err)
					c.Fail(http.StatusInternalServerError, trace(message))
				}
			}
		}()
		c.Next()
	}
}
