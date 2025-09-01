package admin

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"server/app/service"
	"server/app/service/cluster"
	"server/app/service/config"
	"server/app/service/node"
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
