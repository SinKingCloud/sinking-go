package sinking_web

type H map[string]interface{}

//route方法重写
var (
	debug = false
)

// SetDebugMode 设置运行模式为debug
func SetDebugMode(mode bool) {
	debug = mode
}
