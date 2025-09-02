package log

type Type int //日志类型

const (
	EventLogin  Type = iota
	EventShow        = 1
	EventDelete      = 2
	EventUpdate      = 3
	EventCreate      = 4
)

// Types 类型数据
func (s *Service) Types() map[Type]string {
	return map[Type]string{
		EventLogin:  "系统登陆",
		EventShow:   "查看数据",
		EventDelete: "删除数据",
		EventUpdate: "修改数据",
		EventCreate: "创建数据",
	}
}
