package bootstrap

import (
	"log"

	"server/global"
)

func LoadLog() {
	global.App.SetLog(log.Default())
}
