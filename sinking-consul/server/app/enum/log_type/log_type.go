package log_type

const (
	EventLogin  = iota
	EventShow   = 1
	EventDelete = 2
	EventUpdate = 3
	EventCreate = 4
)

// Map 类型数据
func Map() map[int]string {
	return map[int]string{
		EventLogin:  "系统登陆",
		EventShow:   "查看数据",
		EventDelete: "删除数据",
		EventUpdate: "修改数据",
		EventCreate: "创建数据",
	}
}
