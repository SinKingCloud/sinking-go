package admin

import (
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"server/app/constant"
	"server/app/service"
	"server/app/service/log"
	"server/app/util"
	"server/app/util/page"
	"server/app/util/server"
	"server/app/util/str"
)

type ControllerPerson struct {
}

func (ControllerPerson) Info(c *server.Context) {
	user := c.GetUserInfo()
	c.SuccessWithData("获取成功", sinking_web.H{
		"account":    util.Conf.GetString(constant.AuthAccount),
		"login_ip":   user.LoginIp,
		"login_time": user.LoginTime,
	})
}

func (ControllerPerson) Password(c *server.Context) {
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
		c.Error("修改失败")
	} else {
		service.Log.Create(c.GetRequestIp(), log.EventUpdate, "修改登录信息", "修改登录密码")
		c.Success("修改成功")
	}
}

func (ControllerPerson) Log(c *server.Context) {
	pageInfo := page.ValidatePageDefault(c)
	type Form struct {
		OrderByField    string `json:"order_by_field" default:"id" validate:"oneof=id type ip create_time update_time" label:"排序字段"`
		OrderByType     string `json:"order_by_type" default:"desc" validate:"oneof=desc asc" label:"排序类型"`
		Type            string `json:"type" default:"" validate:"omitempty,numeric" label:"类型"`
		Ip              string `json:"ip" default:"" validate:"omitempty" label:"IP地址"`
		Title           string `json:"title" default:"" validate:"omitempty" label:"标题"`
		Content         string `json:"content" default:"" validate:"omitempty" label:"内容"`
		CreateTimeStart string `json:"create_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建起始时间"`
		CreateTimeEnd   string `json:"create_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"创建结束时间"`
		UpdateTimeStart string `json:"update_time_start" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新起始时间"`
		UpdateTimeEnd   string `json:"update_time_end" default:"" validate:"omitempty,datetime=2006-01-02 15:04:05" label:"更新结束时间"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	where := make(map[string]string)
	if form.Ip != "" {
		where["ip"] = form.Ip
	}
	if form.Type != "" {
		where["type"] = form.Type
	}
	if form.Title != "" {
		where["title"] = form.Title
	}
	if form.Content != "" {
		where["content"] = form.Content
	}
	if form.CreateTimeStart != "" {
		where["create_time_start"] = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where["create_time_end"] = form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where["update_time_start"] = form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where["update_time_end"] = form.UpdateTimeEnd
	}
	data, total, err := service.Log.Select(where, form.OrderByField, form.OrderByType, pageInfo.Page, pageInfo.PageSize)
	if err != nil {
		c.Error("获取失败")
	} else {
		service.Log.Create(c.GetRequestIp(), log.EventShow, "查看系统日志", "查看系统日志列表")
		c.SuccessWithData("获取成功", page.NewPage(total, pageInfo.Page, pageInfo.PageSize, data))
	}
}
