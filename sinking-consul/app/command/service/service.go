package service

func Run() {
	//检测服务
	checkCluster()
	//注册服务
	registerCluster()
	//同步服务
	synchronize()
}
