package bootstrap

import (
	"time"

	"server/app/util/cache"
	"server/global"
)

// LoadCache 初始化缓存
func LoadCache() {
	if global.App.Cache != nil {
		return
	}
	global.App.SetCache(cache.NewCache(3600*time.Second, 60*time.Second))
}
