package cluster_status

const (
	Online  = iota //在线
	Offline        //离线
)

// Map 状态数据
func Map() map[int]string {
	return map[int]string{
		Online:  "在线",
		Offline: "离线",
	}
}
