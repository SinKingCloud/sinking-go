package auth

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"server/app/constant"
	"server/app/util"
	"server/app/util/server"
)

func Info(c *server.Context) {
	config := util.Conf
	c.SuccessWithData("获取成功", sinking_web.H{
		"title":    c.GetStringWithDefault(config.GetString(constant.WebTitle), "豁者云服务"),
		"name":     c.GetStringWithDefault(config.GetString(constant.WebName), "云上豁者"),
		"keywords": c.GetStringWithDefault(config.GetString(constant.WebKeyWords), "一站式云服务"),
		"describe": c.GetStringWithDefault(config.GetString(constant.WebDescribe), "提供一站式云服务解决方案"),
		"ui": sinking_web.H{
			"layout":    c.GetStringWithDefault(config.GetString(constant.UiLayout), "left"),
			"watermark": c.GetBoolWithDefault(config.GetString(constant.UiWaterMark), false),
			"theme":     c.GetStringWithDefault(config.GetString(constant.UiTheme), "light"),
			"compact":   c.GetBoolWithDefault(config.GetString(constant.UiCompact), false),
			"color":     c.GetStringWithDefault(config.GetString(constant.UiColor), "rgb(0,81,235)"),
		},
	})
}
