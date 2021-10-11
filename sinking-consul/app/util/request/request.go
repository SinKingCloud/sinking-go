package request

import (
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"github.com/imroc/req"
	"log"
)

type Request struct {
	Ip   string
	Port string
}

func (request *Request) Register() {
	r := req.New()
	header := req.Header{
		"Accept": "application/json",
		setting.GetSystemConfig().Servers.TokenName: setting.GetSystemConfig().Servers.Token,
	}
	param := req.Param{
		"name": "imroc",
		"cmd":  "add",
	}
	// 只有url必选，其它参数都是可选
	_, err := r.Post("http://127.0.0.1", header, param)
	if err != nil {
		log.Fatal(err)
	}
}
