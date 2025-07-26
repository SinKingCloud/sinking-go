package bootstrap

import (
	"log"
	"server/app/util"
)

func LoadLog() {
	util.Log = log.Default()
}
