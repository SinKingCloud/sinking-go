package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/service"
	"github.com/SinKingCloud/sinking-go/sinking-consul/app/util/setting"
	"github.com/imroc/req"
	"time"
)

type Request struct {
	Ip      string `form:"ip" json:"ip"`
	Port    string `form:"port" json:"port"`
	Timeout int    `form:"timeout" json:"timeout"`
}

func header() req.Header {
	header := req.Header{
		"Accept": "application/json",
		setting.GetSystemConfig().Servers.TokenName: setting.GetSystemConfig().Servers.Token,
	}
	return header
}

func (request *Request) SetTimeout(time int) *Request {
	if time > 0 {
		request.Timeout = time
	} else {
		request.Timeout = 5
	}
	return request
}

func (request *Request) Register() bool {
	r := req.New()
	r.SetTimeout(time.Duration(request.Timeout) * time.Second)
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

func (request *Request) ServiceList() ([]*service.Service, error) {
	r := req.New()
	r.SetTimeout(time.Duration(request.Timeout) * time.Second)
	res, err := r.Post(fmt.Sprintf("http://%s:%s/api/cluster/services", request.Ip, request.Port), header())
	if err != nil {
		return nil, err
	}
	data := &ServiceResult{}
	err = json.Unmarshal(res.Bytes(), data)
	if err != nil {
		return nil, err
	}
	if data.Code == 200 {
		return data.Data, nil
	}
	return nil, errors.New(data.Message)
}

func (request *Request) ConfigList() ([]*service.Config, error) {
	r := req.New()
	r.SetTimeout(time.Duration(request.Timeout) * time.Second)
	res, err := r.Post(fmt.Sprintf("http://%s:%s/api/cluster/configs", request.Ip, request.Port), header(), req.BodyJSON(request))
	if err != nil {
		return nil, err
	}
	data := &ConfigResult{}
	err = json.Unmarshal(res.Bytes(), data)
	if err != nil {
		return nil, err
	}
	if data.Code == 200 {
		return data.Data, nil
	}
	return nil, errors.New(data.Message)
}
