package sync

import (
	"server/app/service"
)

func exec(task *Task) {
	switch task.Type {
	case RegisterService:
		_ = service.Cluster.RegisterRemoteService(task.RemoteAddress)
		break
	case SynchronizeData:
		_ = service.Cluster.SynchronizeRemoteData(task.RemoteAddress)
		break
	}
}
