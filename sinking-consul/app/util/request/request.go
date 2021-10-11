package request

import (
	"encoding/json"
	"fmt"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"github.com/imroc/req"
	"time"
)

type Request struct {
	Ip   string `form:"ip" json:"ip"`
	Port string `form:"port" json:"port"`
}

func header() req.Header {
	header := req.Header{
		"Accept": "application/json",
		setting.GetSystemConfig().Servers.TokenName: setting.GetSystemConfig().Servers.Token,
	}
	return header
}

func (request *Request) Register() bool {
	r := req.New()
	r.SetTimeout(5 * time.Second)
	res, err := r.Post(fmt.Sprintf("http://%s:%s/api/cluster/register", request.Ip, request.Port), header(), req.BodyJSON(request))
	if err != nil {
		return false
	}
	data := &Result{}
	err = json.Unmarshal(res.Bytes(), data)
	if err != nil {
		return false
	}
	if data.Code == 200 {
		return true
	}
	return false
}
