package global

import "server/app/constant"

// IsDebug 是否为调试模式。
func (a *Application) IsDebug() bool {
	return a.Config.GetString(constant.ServerMode) == "dev"
}

// ServerAddr 获取服务监听地址。
func (a *Application) ServerAddr() (host string, port int) {
	host = a.Config.GetString(constant.ServerHost)
	port = a.Config.GetInt(constant.ServerPort)
	if host == "" {
		host = "0.0.0.0"
	}
	if port <= 0 {
		port = 5678
	}
	return host, port
}
