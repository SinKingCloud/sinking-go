package admin

import (
	"server/app/constant"
	"server/app/enum"
	"server/app/enum/cluster_status"
	"server/app/enum/config_status"
	"server/app/enum/log_type"
	"server/app/enum/node_online_status"
	"server/app/service"
	"server/app/service/setting"
	"server/app/util/context"

	"github.com/SinKingCloud/sinking-go/sinking-web"
)

type ControllerSystem struct{}

func (ControllerSystem) Overview(c *context.Context) {
	conf := service.Setting
	clusterNum, _ := service.Cluster.CountAll()
	clusterOnlineNum, _ := service.Cluster.CountByStatus(cluster_status.Online)
	nodeNum, _ := service.Node.CountAll()
	nodeOnlineNum, _ := service.Node.CountByOnlineStatus(node_online_status.Online)
	configNum, _ := service.Config.CountAll()
	configNormalNum, _ := service.Config.CountByStatus(config_status.Normal)
	c.SuccessWithData("获取成功", sinking_web.H{
		"application": sinking_web.H{
			"mode":    c.GetStringWithDefault(conf.GetString(constant.ServerMode), "release"),
			"listen":  conf.GetString(constant.ServerHost) + ":" + c.GetStringWithDefault(conf.GetString(constant.ServerPort), "5678"),
			"address": service.Cluster.GetLocalAddr(),
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

func (ControllerSystem) Enum(c *context.Context) {
	type Form struct {
		Name string `json:"name" default:"" validate:"required" label:"枚举名称"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	if value, ok := enum.Data[form.Name]; ok {
		c.SuccessWithData("获取数据成功", value)
	} else {
		c.Error("枚举类型不存在")
	}
}

func (ControllerSystem) Config(c *context.Context) {
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
		var list []*setting.Config
		for _, v := range form.Configs {
			if v.Key != "" {
				list = append(list, &setting.Config{
					Key:   form.Group + "." + v.Key,
					Value: v.Value,
				})
			}
		}
		err := service.Setting.Sets(list)
		if err == nil {
			service.Log.Create(c.GetRequestIp(), log_type.EventUpdate, "修改系统配置", "修改系统配置["+form.Group+"]数据")
			c.Success("保存成功")
		} else {
			c.Error("保存失败")
		}
	} else {
		service.Log.Create(c.GetRequestIp(), log_type.EventShow, "查看系统配置", "查看系统配置["+form.Group+"]列表")
		c.SuccessWithData("获取成功", service.Setting.GetByGroup(form.Group))
	}
}
