package sinking_web

import (
	"github.com/SinKingCloud/sinking-go/sinking-web/constant/message"
	"net/http"
)

//route方法重写
var route routeStruct

type routeStruct struct {
	NotFound func(*Context)
}

func SetNotFoundHandle(fun func(*Context)) {
	route.NotFound = fun
}

func NotFoundHandle(c *Context) {
	if route.NotFound != nil {
		route.NotFound(c)
	} else {
		c.JSON(http.StatusNotFound, H{"code": http.StatusNotFound, "message": message.NotFound})
	}
}

//context方法重写
var context contextStruct

type contextStruct struct {
	Fail       func(c *Context, code int, message string)
	JsonFormat func()
}

func SetFailHandle(fun func(c *Context, code int, message string)) {
	context.Fail = fun
}

func FailHandle(c *Context, code int, message string) {
	if context.Fail != nil {
		context.Fail(c, code, message)
	} else {
		c.JSON(code, H{"code": code, "message": message})
	}
}
