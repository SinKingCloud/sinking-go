package sinking_sdk_go

import (
	"fmt"
	"net/http"
	"testing"
)

func Test_main(t *testing.T) {
	//实例化一个server
	server := New("42.157.128.40:1817", "sinking-token", "97cb316ec237b3937307d94c38d21785", "cloud-server", "dev")
	//注册并监听服务
	server.Register("default", "cloud-gateway", "127.0.0.1:1000"). //服务信息
									UseService(map[string][]string{"default": {"cloud-gateway"}}). //需要使用的服务
									Listen()                                                       //监听
	//rpc调用服务
	body, err := server.Rpc("default", "cloud-gateway").
		Timeout(1).                                                          //超时熔断
		Method(http.MethodPost).                                             //请求类型
		ReTry(5).                                                            //重试次数
		Call("/index/login", &Param{"user": "admin", "pwd": "123456"}, Poll) //携带参数
	fmt.Println(body, err)
	//获取配置
	serverConf := server.Config("default").Name("test").Viper()
	fmt.Println(serverConf.Get("host"))
}
