package api

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/response"
	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// ServiceRegister 注册服务
func ServiceRegister(s *sinking_web.Context) {
	type register struct {
		Ip   string `form:"ip" json:"ip"`
		Port string `form:"port" json:"port"`
	}
	response.Success(s, "注册服务成功", nil)
}
