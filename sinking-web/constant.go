package sinking_web

import "math"

const (
	FrameWorkVersion           = "1.0.0"                 //框架版本
	MessageNotFound            = "404 Not Found"         //默认资源不存在消息
	MessageInternalServerError = "Internal Server Error" //默认错误消息
)

const (
	ContentType = "Content-Type" //主体类型
)

const (
	ContentTypeJson = "application/json;charset=utf-8;" //返回json的头
	ContentTypeText = "text/plain;charset=utf-8;"       //返回字符的头
	ContentTypeHtml = "text/html;charset=utf-8;"        //返回html的头
	HeaderLocation  = "Location"                        //重定向跳转
)

const (
	BindFormTagName         = "form"    //表单绑定参数tag名称
	BindDefaultValueTagName = "default" //表单绑定默认值tag名称
)

const (
	abortIndex             int = math.MaxInt >> 1 // 中间件锁
	defaultMultipartMemory     = 32 << 20         // 默认上传buf大小
)
