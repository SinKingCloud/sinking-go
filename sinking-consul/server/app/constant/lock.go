package constant

import "time"

const (
	LockConfigSet     = "ConfigSet_"     //修改配置并发锁
	LockTimeConfigSet = 60 * time.Second //修改配置并发锁最大时间
)
