package service

func Run() {
	//检测服务
	checkCluster()
	//同步服务
	synchronize()
	//注册服务
	registerCluster()
}
