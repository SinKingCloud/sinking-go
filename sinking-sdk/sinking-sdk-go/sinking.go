package sinking_sdk_go

var (
	checkTime = 5 //轮询间隔
)

// Register 注册中心
type Register struct {
	Servers   string `json:"-"`          //注册中心
	TokenName string `json:"-"`          //通信密匙名称
	Token     string `json:"-"`          //通信密匙
	Name      string `json:"name"`       //服务名称
	AppName   string `json:"app_name"`   //所属应用
	EnvName   string `json:"env_name"`   //环境标识
	GroupName string `json:"group_name"` //分组名称
	Addr      string `json:"service"`    //服务地址(规则ip:port)
}

// New 实例化
func New(server string, tokenName string, token string, name string, appName string, envName string, groupName string, addr string) *Register {
	return &Register{
		Servers:   server,
		TokenName: tokenName,
		Token:     token,
		Name:      name,
		AppName:   appName,
		EnvName:   envName,
		GroupName: groupName,
		Addr:      addr,
	}
}

// Listen 监听配置变动及发送服务心跳
func (r *Register) Listen() {
	r.registerServices() //注册节点并维持心跳
	r.getServices()      //监听服务列表
}
