package admin

import (
	"server/app/constant"
	"server/app/enum/log_type"
	"server/app/repository/log"
	"server/app/service"
	"server/app/util/context"
	"server/app/util/page"

	"github.com/SinKingCloud/sinking-go/sinking-web"
)

type ControllerPerson struct {
}

func (ControllerPerson) Info(c *context.Context) {
	user := c.GetUserInfo()
	c.SuccessWithData("获取成功", sinking_web.H{
		"account":    service.Setting.GetString(constant.AuthAccount),
		"login_ip":   user.LoginIp,
		"login_time": user.LoginTime,
	})
}

func (ControllerPerson) Password(c *context.Context) {
	type Form struct {
		Password string `json:"password" default:"" validate:"required" label:"密码"`
	}
	form := &Form{}
	if ok, msg := c.ValidatorAll(form); !ok {
		c.Error(msg)
		return
	}
	err := service.Auth.ChangePassword(form.Password)
	if err != nil {
		c.Error(err.Error())
	} else {
		service.Log.Create(c.GetRequestIp(), log_type.EventUpdate, "修改登录信息", "修改登录密码")
		c.Success("修改成功")
	}
}

func (ControllerPerson) Log(c *context.Context) {
	pageNum, pageSize := c.ValidatePage()
	orderByField, orderByType := c.ValidateOrderBy("id", "desc", "id,type,ip,create_time,update_time")
	type Form struct {
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
	where := &log.SelectLog{}
	if form.Ip != "" {
		where.Ip = form.Ip
	}
	if form.Type != "" {
		where.Type = form.Type
	}
	if form.Title != "" {
		where.Title = form.Title
	}
	if form.Content != "" {
		where.Content = form.Content
	}
	if form.CreateTimeStart != "" {
		where.CreateTimeStart = form.CreateTimeStart
	}
	if form.CreateTimeEnd != "" {
		where.CreateTimeEnd = form.CreateTimeEnd
	}
	if form.UpdateTimeStart != "" {
		where.UpdateTimeStart = form.UpdateTimeStart
	}
	if form.UpdateTimeEnd != "" {
		where.UpdateTimeEnd = form.UpdateTimeEnd
	}
	data, total, err := service.Log.Select(where, orderByField, orderByType, pageNum, pageSize)
	if err != nil {
		c.Error("获取失败")
	} else {
		service.Log.Create(c.GetRequestIp(), log_type.EventShow, "查看系统日志", "查看系统日志列表")
		c.SuccessWithData("获取成功", page.NewPage(total, pageNum, pageSize, data))
	}
}
