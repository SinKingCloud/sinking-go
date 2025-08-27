package admin

import (
	sinking_web "github.com/SinKingCloud/sinking-go/sinking-web"
	"server/app/constant"
	"server/app/service"
	"server/app/service/cluster"
	"server/app/service/config"
	"server/app/service/node"
	"server/app/util"
	"server/app/util/server"
	"server/app/util/str"
)

type ControllerSystem struct{}

func (ControllerSystem) Password(c *server.Context) {
	type Form struct {
		Password string `json:"password" default:"" validate:"required" label:"密码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	sPwd, _ := str.NewStringTool().BcryptHash(form.Password)
	util.Conf.Set(constant.AuthPassword, sPwd)
	if util.Conf.WriteConfig() != nil {
		c.Error("修改密码失败")
	} else {
		c.Success("修改密码成功")
	}
}

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
