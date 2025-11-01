package context

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"net/http"
)

const (
	SuccessCode    = 200 //请求成功
	ErrorCode      = 500 //请求失败
	TokenErrorCode = 503 //token认证失败
	NotLoginCode   = 403 //未登录
)

func (c *Context) response(code int32, message interface{}, data interface{}) {
	if data == nil {
		data = sinking_web.H{}
	}
	c.JSON(http.StatusOK, sinking_web.H{
		"code":       code,
		"message":    message,
		"data":       data,
		"request_id": c.GetRequestId(),
	})
}

func (c *Context) Return(code int32, message interface{}, data interface{}) {
	c.response(code, message, data)
}

func (c *Context) Success(message interface{}) {
	c.response(SuccessCode, message, nil)
}

func (c *Context) SuccessWithData(message interface{}, data interface{}) {
	c.response(SuccessCode, message, data)
}

func (c *Context) Error(message interface{}) {
	c.response(ErrorCode, message, nil)
}

func (c *Context) ErrorWithData(message interface{}, data interface{}) {
	c.response(ErrorCode, message, data)
}

func (c *Context) TokenError(message interface{}, data interface{}) {
	c.response(TokenErrorCode, message, data)
}

func (c *Context) NotLogin(message interface{}, data interface{}) {
	c.response(NotLoginCode, message, data)
}
