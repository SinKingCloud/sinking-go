package global

import (
	"log"

	"server/app/util/cache"
	"server/app/util/daemon"
	"server/app/util/database"

	"github.com/spf13/viper"
)

// Application 应用全局依赖。
type Application struct {
	Database *database.Database
	Config   *viper.Viper
	Cache    cache.Interface
	Log      *log.Logger
	Daemon   *daemon.Daemon
}
