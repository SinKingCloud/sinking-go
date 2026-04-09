package save

import (
	"server/app/service"
	"time"
)

var (
	saveInterval = time.Minute // 保存间隔时间
)

// Init 定时保存数据
func Init() {
	go func() {
		ticker := time.NewTicker(saveInterval)
		defer ticker.Stop()
		for range ticker.C {
			_ = service.Cluster.Save()
			_ = service.Node.Save()
			_ = service.Config.Save()
		}
	}()
}
