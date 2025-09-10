package constant

import "time"

const (
	LockConfigCreate     = "configCreate_"      //创建配置并发锁
	LockTimeConfigCreate = 1 * 60 * time.Second //创建配置并发锁最大时间

	LockNodeCreate     = "configCreate_"      //创建配置并发锁
	LockTimeNodeCreate = 1 * 60 * time.Second //创建配置并发锁最大时间

	LockSyncData     = "syncData_"          //同步数据并发锁
	LockTimeSyncData = 5 * 60 * time.Second //同步数据并发锁最大时间
)
