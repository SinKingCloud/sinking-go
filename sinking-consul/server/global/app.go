package global

import (
	"log"

	"server/app/util/cache"
	"server/app/util/daemon"
	"server/app/util/database"

	"github.com/spf13/viper"
)

var App = &Application{}

// SetDataBase 设置数据库实例，支持链式调用。
func (a *Application) SetDataBase(database *database.Database) *Application {
	a.Database = database
	return a
}

// SetConfig 设置配置实例，支持链式调用。
func (a *Application) SetConfig(config *viper.Viper) *Application {
	a.Config = config
	return a
}

// SetCache 设置缓存实例，支持链式调用。
func (a *Application) SetCache(cache cache.Interface) *Application {
	a.Cache = cache
	return a
}

// SetLog 设置日志实例，支持链式调用。
func (a *Application) SetLog(logger *log.Logger) *Application {
	a.Log = logger
	return a
}

// SetDaemon 设置进程管理器，支持链式调用。
func (a *Application) SetDaemon(manager *daemon.Daemon) *Application {
	a.Daemon = manager
	return a
}
