package sinking_sdk_go

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	//实例化一个server
	server := New("127.0.0.1:8888", "sinking-token", "test_token", "default", "dev")
	//注册并监听服务
	server.Register("sinking-go-api", "sinking-go-api-order", "106.52.89.187"). //服务信息
											UseService(map[string]string{"sinking-go-api": "sinking-go-api-order"}). //需要使用的服务
											Listen()                                                                 //监听

	//延迟5s等待初始化完毕
	time.Sleep(5 * time.Second)

	//rpc调用服务
	body, err := server.Rpc("sinking-go-api", "sinking-go-api-order").
		Timeout(5).                                                    //超时熔断
		Method(http.MethodPost).                                       //请求类型
		ReTry(5).                                                      //重试次数
		Call("/index/login", &Param{"user": "admin", "pwd": "123456"}) //携带参数
	fmt.Println(body, err)

	//获取配置
	serverConf := server.Config("sinking-go-api").Name("database").Viper()
	fmt.Println(serverConf.Get("host"))
}
