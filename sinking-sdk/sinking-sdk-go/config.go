package sinking_sdk_go

import (
	"github.com/spf13/viper"
	"strings"
	"sync"
	"time"
)

// Config 配置结构体
type Config struct {
	AppName   string `json:"app_name"`
	EnvName   string `json:"env_name"`
	GroupName string `json:"group_name"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	Hash      string `json:"hash"`
	Type      string `json:"type"`
	Status    int    `json:"status"`
	viper     *viper.Viper
}

var (
	// configs 配置列表 GroupName.Name
	configs     = make(map[string]map[string]*Config)
	configsLock sync.Mutex
)

// getConfigs 获取配置
func (r *Register) getConfigs(sync bool) {
	//设置注册节点
	fun := func() {
		for {
			test := &RequestServer{
				Server:    r.server,
				TokenName: r.TokenName,
				Token:     r.Token,
			}
			result := test.getConfigs(r.AppName, r.EnvName)
			if result != nil && result.Code == 200 {
				//解析配置
				for _, v := range result.Data {
					configsLock.Lock()
					if configs[v.GroupName] == nil {
						configs[v.GroupName] = map[string]*Config{}
					}
					if configs[v.GroupName][v.Name] == nil || v.Hash != configs[v.GroupName][v.Name].Hash {
						configs[v.GroupName][v.Name] = v
						conf := viper.New()
						conf.SetConfigType(v.Type)
						err := conf.ReadConfig(strings.NewReader(v.Content))
						if err == nil {
							configs[v.GroupName][v.Name].viper = conf
						}
					}
					configsLock.Unlock()
				}
			}
			time.Sleep(time.Duration(checkTime) * time.Second)
		}
	}
	if sync {
		go fun()
	} else {
		fun()
	}
}

type configBuild struct {
	groupName string
	name      string
}

// Config 实例化config
func (r *Register) Config(groupName string) *configBuild {
	c := &configBuild{
		groupName: groupName,
	}
	return c
}

// Name 配置Name
func (c *configBuild) Name(name string) *configBuild {
	c.name = name
	return c
}

// Viper 获取Viper
func (c *configBuild) Viper() *viper.Viper {
	configsLock.Lock()
	defer configsLock.Unlock()
	if configs[c.groupName][c.name] != nil && configs[c.groupName][c.name].viper != nil {
		return configs[c.groupName][c.name].viper
	}
	return viper.New()
}
