package admin

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/model"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/encode"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/jwt"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"time"
)

// UserLogin 账户登陆
func UserLogin(s *sinking_web.Context) {
	type login struct {
		User string `form:"user" json:"user"`
		Pwd  string `form:"pwd" json:"pwd"`
	}
	form := &login{}
	err := s.BindAll(form)
	if err != nil || form.User == "" || form.Pwd == "" {
		response.Error(s, "参数不足", nil)
		return
	}
	//查询用户
	var user model.User
	model.Db.Where("user = ?", form.User).First(&user)
	if user.Id <= 0 {
		response.Error(s, "账户不存在", nil)
		return
	}
	if encode.ComparePasswords(user.Pwd, form.Pwd) {
		token := jwt.SetToken(user)
		model.Db.Where("id = ?", user.Id).Updates(&model.User{
			LoginTime: model.DateTime(time.Now()),
			LoginIp:   s.ClientIP(false),
		})
		response.Success(s, "登陆成功", token)
	} else {
		response.Error(s, "密码错误", nil)
	}
}
