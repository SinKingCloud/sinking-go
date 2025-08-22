package cluster

type Status int

const (
	Online  Status = iota //在线
	Offline               //离线
)

// Status 状态数据
func (s *Service) Status() map[Status]string {
	return map[Status]string{
		Online:  "在线",
		Offline: "离线",
	}
}
