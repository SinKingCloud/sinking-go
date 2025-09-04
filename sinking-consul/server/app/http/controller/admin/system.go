package admin

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"server/app/constant"
	"server/app/service"
	"server/app/service/cluster"
	"server/app/service/config"
	"server/app/service/log"
	"server/app/service/node"
	"server/app/util"
	"server/app/util/server"
)

type ControllerSystem struct{}

func (ControllerSystem) Overview(c *server.Context) {
	clusterNum, _ := service.Cluster.CountAll()
	clusterOnlineNum, _ := service.Cluster.CountByStatus(cluster.Online)
	nodeNum, _ := service.Node.CountAll()
	nodeOnlineNum, _ := service.Node.CountByOnlineStatus(node.Online)
	configNum, _ := service.Config.CountAll()
	configNormalNum, _ := service.Config.CountByStatus(config.Normal)
	c.SuccessWithData("获取成功", sinking_web.H{
		"application": sinking_web.H{
			"mode":    c.GetStringWithDefault(util.Conf.GetString(constant.ServerMode), "release"),
			"address": util.Conf.GetString(constant.ServerHost) + ":" + c.GetStringWithDefault(util.Conf.GetString(constant.ServerPort), "5678"),
		},
		"cluster": sinking_web.H{
			"total":   clusterNum,
			"online":  clusterOnlineNum,
			"offline": clusterNum - clusterOnlineNum,
		},
		"node": sinking_web.H{
			"total":   nodeNum,
			"online":  nodeOnlineNum,
			"offline": nodeNum - nodeOnlineNum,
		},
		"config": sinking_web.H{
			"total":    configNum,
			"normal":   configNormalNum,
			"abnormal": configNum - configNormalNum,
		},
	})
}

func (ControllerSystem) Enum(c *server.Context) {
	type Form struct {
		Name string `json:"name" default:"" validate:"required,name" label:"枚举名称"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if value, ok := service.Enum[form.Name]; ok {
		c.SuccessWithData("获取数据成功", value)
	} else {
		c.Error("枚举类型不存在")
	}
}

func (ControllerSystem) Config(c *server.Context) {
	type Config struct {
		Key   string `json:"key" default:"" validate:"required,max=50" label:"配置标识"`
		Value string `json:"value" default:"" validate:"omitempty" label:"配置内容"`
	}
	type Form struct {
		Action  string    `json:"action" default:"get" validate:"oneof=get set" label:"操作类型"`
		Group   string    `json:"group" default:"web" validate:"oneof=web ui" label:"配置分组"`
		Configs []*Config `json:"configs" default:"" validate:"omitempty,gte=1" label:"配置标识"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if form.Action == "set" {
		for _, v := range form.Configs {
			if v.Key != "" {
				util.Conf.Set(form.Group+"."+v.Key, v.Value)
			}
		}
		err := util.Conf.WriteConfig()
		if err == nil {
			service.Log.Create(c.GetRequestIp(), log.EventUpdate, "修改系统配置", "修改系统配置["+form.Group+"]数据")
			c.Success("保存成功")
		} else {
			c.Error("保存失败")
		}
	} else {
		service.Log.Create(c.GetRequestIp(), log.EventShow, "查看系统配置", "查看系统配置["+form.Group+"]列表")
		c.SuccessWithData("获取成功", util.Conf.AllSettings()[form.Group])
	}
}
