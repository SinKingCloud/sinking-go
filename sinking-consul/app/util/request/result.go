package request

import "github.com/SinKingCloud/sinking-go/sinking-consul/app/service"

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type ConfigResult struct {
	Code    int               `json:"code"`
	Data    []*service.Config `json:"data"`
	Message string            `json:"message"`
}

type ServiceResult struct {
	Code    int                `json:"code"`
	Data    []*service.Service `json:"data"`
	Message string             `json:"message"`
}
