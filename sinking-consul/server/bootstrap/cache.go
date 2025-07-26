package bootstrap

import (
	"server/app/util"
	"server/app/util/cache"
	"time"
)

// LoadCache 初始化缓存
func LoadCache() {
	util.Cache = cache.NewCache(3600*time.Second, 60*time.Second)
}
