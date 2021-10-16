package sinking_sdk_go

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

var server *Register

func Test_main(t *testing.T) {
	//实例化一个server
	server = New("106.52.89.187:80", "sinking-token", "test_token", "default", "dev")
	//注册并监听服务
	server.Register("sinking-go-api", "sinking-go-api-order", "106.52.89.187").UseService(map[string]string{
		"sinking-go-api": "sinking-go-api-order", //需要使用的服务
	}).Listen()
	time.Sleep(5 * time.Second) //延迟5s等待初始化完毕
	//rpc调用服务
	body, err := server.Rpc("sinking-go-api", "sinking-go-api-order").Timeout(5).Method(http.MethodPost).ReTry(5).Call("/index/login", &Param{
		"user": "admin",
		"pwd":  "123456",
	})
	fmt.Println(body, err)
	//获取配置
	serverConf := server.Config("sinking-go-api").Name("database").Viper()
	fmt.Println(serverConf.Get("host"))
}
