package sinking_sdk_go

import (
	"fmt"
	"strings"
)

type Register struct {
	Servers   string `json:"-"`          //注册中心
	TokenName string `json:"-"`          //通信密匙名称
	Token     string `json:"-"`          //通信密匙
	Name      string `json:"name"`       //服务名称
	AppName   string `json:"app_name"`   //所属应用
	EnvName   string `json:"env_name"`   //环境标识
	GroupName string `json:"group_name"` //分组名称
	Addr      string `json:"addr"`       //服务地址(规则ip:port)
}

// New 实例化
func New(server string, token string, name string, appName string, envName string, groupName string, addr string) *Register {
	return &Register{
		Servers:   server,
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

}

func (r *Register) getServers() {
	//设置注册节点
	servers := strings.Split(r.Servers, ",")
	for _, v := range servers {
		fmt.Println(v)
	}
}
