package sync

type Type int //推送类型

const (
	RegisterService Type = iota //注册服务任务
	SynchronizeData             //同步数据任务
)
