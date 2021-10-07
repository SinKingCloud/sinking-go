package sinking_web

import "time"

type H map[string]interface{}

//route方法重写
var (
	debug        = false
	readTimeOut  = time.Second * 600
	writeTimeout = time.Second * 600
)

// SetDebugMode 设置运行模式为debug
func SetDebugMode(mode bool) {
	debug = mode
}

// SetTimeOut 设置超时时间
func SetTimeOut(read time.Duration, write time.Duration) {
	readTimeOut = read
	writeTimeout = write
}
