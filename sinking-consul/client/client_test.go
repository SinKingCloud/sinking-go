package client

import (
	"testing"
)

func Test(t *testing.T) {
	cli := NewClient([]string{"集群地址1", "集群地址2"}, "服务组", "服务名", "本机应用地址", "集群密钥")
	_ = cli.Connect()                  //链接集群
	_, _ = cli.GetAllService()         //获取所有服务
	_, _ = cli.GetAllConfigs()         //获取所有配置
	_, _ = cli.GetConfig("配置名")        //获取配置
	_, _ = cli.GetService("服务名", Poll) //获取服务节点,支持轮询和随机
	_ = cli.Close()                    //断开集群
}
