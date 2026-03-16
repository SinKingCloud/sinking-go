package client

import (
	"fmt"
	"testing"
	"time"
)

func Test(t *testing.T) {
	go cli1()
	go cli2()
	time.Sleep(time.Hour)
	//cli := NewClient([]string{"集群地址1", "集群地址2"}, "服务组", "服务名", "本机应用地址", "集群密钥")
	//_ = cli.Connect()                             //链接集群
	//_, _ = cli.GetAllService()                    //获取所有服务
	//_, _ = cli.GetAllConfigs()                    //获取所有配置
	//_, _ = cli.GetConfig("配置名")                   //获取配置
	//_, _ = cli.GetService("服务名", Poll)            //获取服务节点-轮询模式
	//_, _ = cli.GetService("服务名", Rand)            //获取服务节点-随机模式
	//_, _ = cli.GetService("服务名", Hash, "user123") //获取服务节点-哈希模式，根据指定的键获取节点
	////rpc相关操作
	//cli.RpcRegister("方法名", func(params json.RawMessage) (interface{}, error) {
	//	return map[string]string{"响应参数1": "响应值1"}, nil
	//}) //注册方法
	//cli.RpcHandle().ServeHTTP(nil, nil)                                     //获取RPC HTTP处理器 后续挂载在http上即可使用
	//_ = cli.RpcCall("服务名", "方法名", map[string]interface{}{"参数1": "值1"}, nil) //调用远程RPC服务
	//
	//_ = cli.Close() //断开集群
}

func cli1() {
	time.Sleep(1400 * time.Millisecond)
	cli := NewClient([]string{"101.201.124.64:5678"}, "sinking-cloud", "gateway", "127.0.0.1:1000", "E10ADC3949BA59ABBE56E057F20F883E")
	_ = cli.Connect() //链接集群
	for {
		fmt.Println(cli.GetAllService())
		time.Sleep(time.Second)
	}
}

func cli2() {
	cli := NewClient([]string{"101.201.124.64:5678"}, "sinking-cloud", "admin", "127.0.0.1:1008", "E10ADC3949BA59ABBE56E057F20F883E")
	_ = cli.Connect() //链接集群
	for {
		fmt.Println(cli.GetAllService())
		time.Sleep(time.Second)
	}
}
