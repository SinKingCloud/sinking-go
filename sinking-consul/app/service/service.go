package service

// Services 服务列表
var Services = make(map[string]*Service)

// Service 服务列表
type Service struct {
	Name          string `json:"name"`            //服务名称
	AppName       string `json:"app_name"`        //所属应用
	EnvName       string `json:"env_name"`        //环境标识
	GroupName     string `json:"group_name"`      //分组名称
	Addr          string `json:"addr"`            //服务地址(规则ip:port)
	ServiceHash   string `json:"service_hash"`    //标识hash(规则md5(Addr))
	LastHeartTime int64  `json:"last_heart_time"` //上次心跳时间
	Status        int    `json:"status"`          //服务状态(0:正常/1:异常)
}
