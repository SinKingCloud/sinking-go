package node

type Status int

const (
	Normal Status = iota //正常
	Stop                 //暂停
)

// Status 状态数据
func (s *Service) Status() map[Status]string {
	return map[Status]string{
		Normal: "正常",
		Stop:   "暂停",
	}
}

type OnlineStatus int

const (
	Online  OnlineStatus = iota //在线
	Offline                     //离线
)

// OnlineStatus 在线状态数据
func (s *Service) OnlineStatus() map[OnlineStatus]string {
	return map[OnlineStatus]string{
		Online:  "在线",
		Offline: "离线",
	}
}
