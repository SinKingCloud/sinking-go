package sinking_sdk_go

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	AppName   string `json:"app_name"`
	EnvName   string `json:"env_name"`
	GroupName string `json:"group_name"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	Hash      string `json:"hash"`
	Type      string `json:"type"`
	Status    int    `json:"status"`
}

// getConfigs 获取配置
func (r *Register) getConfigs() {
	//设置注册节点
	go func() {
		for {
			servers := strings.Split(r.Servers, ",")
			for _, v := range servers {
				test := &RequestServer{
					Server:    v,
					TokenName: r.TokenName,
					Token:     r.Token,
				}
				result := test.getConfigs(r.AppName, r.EnvName)
				fmt.Println(result.Data[0])
			}
			time.Sleep(time.Duration(checkTime) * time.Second)
		}
	}()
}
