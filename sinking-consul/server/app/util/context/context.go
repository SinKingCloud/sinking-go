package context

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Context 接管上下文方便扩展方法
type Context struct {
	*sinking_web.Context
}

// HandleFunc 方法转换
func HandleFunc(handler func(c *Context)) func(ctx *sinking_web.Context) {
	return func(c *sinking_web.Context) {
		handler(&Context{
			Context: c,
		})
	}
}
