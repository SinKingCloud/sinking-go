package service

import (
	"server/app/service/auth"
)

// service实例
var (
	Auth = auth.GetIns()
)

// Enum 枚举信息
var Enum = map[string]interface{}{}
