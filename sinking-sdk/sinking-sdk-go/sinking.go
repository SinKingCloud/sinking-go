package sinking_sdk_go

import (
	"math/rand"
	"strings"
)

var (
	checkTime = 5 //轮询间隔
)

// Register 注册中心
type Register struct {
	Servers    string            `json:"-"`          //注册中心
	TokenName  string            `json:"-"`          //通信密匙名称
	Token      string            `json:"-"`          //通信密匙
	Name       string            `json:"name"`       //服务名称
	AppName    string            `json:"app_name"`   //所属应用
	EnvName    string            `json:"env_name"`   //环境标识
	GroupName  string            `json:"group_name"` //分组名称
	Addr       string            `json:"addr"`       //服务地址(规则ip:port)
	server     string            //使用节点
	useService map[string]string //使用服务
}

// New 实例化
func New(server string, tokenName string, token string, appName string, envName string) *Register {
	return &Register{
		Servers:   server,
		TokenName: tokenName,
		Token:     token,
		AppName:   appName,
		EnvName:   envName,
	}
}

// Register 注册服务
func (r *Register) Register(groupName string, name string, addr string) *Register {
	r.GroupName = groupName
	r.Name = name
	r.Addr = addr
	return r
}

// UseService 使用服务
func (r *Register) UseService(use map[string]string) *Register {
	r.useService = use
	return r
}

// changeServer 更改注册中心
func (r *Register) changeServer(rand bool) {
	if !rand {
		r.changeServerByHash() //根据hash修改server
	} else {
		r.changeServerByRand() //根据rand修改server
	}
}

// changeServerByHash 根据hash获取server
func (r *Register) changeServerByRand() {
	data := strings.Split(r.Servers, ",")
	if len(data) <= 1 {
		return
	}
	var temp []string
	for _, v := range data {
		if v != r.server {
			temp = append(temp, v)
		}
	}
	r.server = temp[rand.Intn(len(temp))]
}

// changeServerByHash 根据hash获取server
func (r *Register) changeServerByHash() {
	key := Md5Encode(r.AppName + r.EnvName + r.GroupName + r.Name + r.Addr)
	test := NewConsistent()
	data := strings.Split(r.Servers, ",")
	for _, v := range data {
		test.Add(v)
	}
	server, err := test.Get(key)
	if err != nil {
		return
	}
	r.server = server
}

// Listen 监听配置变动及发送服务心跳
func (r *Register) Listen() {
	r.changeServer(false)    //初始化节点根据hash获取
	r.getConfigs(true)       //监听配置列表
	r.registerServices(true) //注册节点并维持心跳
	r.getServices(true)      //监听服务列表
}

// SetOnline 设置服务上线下线
func (r *Register) SetOnline(online bool, imEf bool) {
	OnlineStatusLock.Lock()
	OnlineStatus = online
	OnlineStatusLock.Unlock()
	if !imEf {
		return
	}
	//更改服务状态(即时)
	status := 0
	if !OnlineStatus {
		status = 1
	}
	r.changeServerStatus(Md5Encode(r.AppName+r.EnvName+r.GroupName+r.Addr), status)
}
