package response

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"net/http"
)

const (
	SuccessCode    = 200 //请求成功
	ErrorCode      = 500 //请求失败
	TokenErrorCode = 503 //token认证失败
)

func response(c *sinking_web.Context, code int32, message interface{}, data interface{}) {
	c.JSON(http.StatusOK, sinking_web.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

func Return(c *sinking_web.Context, code int32, message interface{}, data interface{}) {
	response(c, code, message, data)
}

func Success(c *sinking_web.Context, message interface{}, data interface{}) {
	response(c, SuccessCode, message, data)
}

func Error(c *sinking_web.Context, message interface{}, data interface{}) {
	response(c, ErrorCode, message, data)
}

func TokenError(c *sinking_web.Context, message interface{}, data interface{}) {
	response(c, TokenErrorCode, message, data)
}
