package sync

import (
	"server/app/service"
)

func exec(task *Task) {
	switch task.Type {
	case RegisterService:
		_ = service.Cluster.RegisterService(task.RemoteAddress)
		break
	case SynchronizeData:
		_ = service.Cluster.SynchronizeData(task.RemoteAddress)
		break
	}
}
