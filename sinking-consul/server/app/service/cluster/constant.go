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

type IsDelete int

const (
	False IsDelete = iota
	True
)

// IsDelete 是否删除
func (s *Service) IsDelete() map[IsDelete]string {
	return map[IsDelete]string{
		False: "未删除",
		True:  "已删除",
	}
}
