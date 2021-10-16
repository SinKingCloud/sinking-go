package sinking_sdk_go

import (
	"fmt"
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
			test := &RequestServer{
				Server:    r.server,
				TokenName: r.TokenName,
				Token:     r.Token,
			}
			result := test.getConfigs(r.AppName, r.EnvName)
			fmt.Println(result)
			time.Sleep(time.Duration(checkTime) * time.Second)
		}
	}()
}
