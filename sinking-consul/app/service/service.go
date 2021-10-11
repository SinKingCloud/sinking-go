package service

type Service struct {
	Name       string `json:"name"`        //服务名称
	GroupName  string `json:"group_name"`  //分组名称
	ClusterNum int    `json:"cluster_num"` //服务数量
	Status     int    `json:"status"`      //服务状态
}
