package config_status

const (
	Normal = iota //正常
	Stop          //暂停
)

// Map 状态数据
func Map() map[int]string {
	return map[int]string{
		Normal: "正常",
		Stop:   "暂停",
	}
}
