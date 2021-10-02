package sinking_web

import (
	"fmt"
	message2 "github.com/SinKingCloud/sinking-go/sinking-web/constant/message"
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
				message := fmt.Sprintf("%s", err)
				if debug {
					c.Fail(http.StatusInternalServerError, message2.InternalServerError)
				} else {
					c.Fail(http.StatusInternalServerError, trace(message))
				}
			}
		}()
		c.Next()
	}
}
