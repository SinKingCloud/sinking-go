package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"server/app"
	"server/app/util"
	"server/app/util/daemon"
	"server/bootstrap"
)

const (
	pidFileName  = "server.pid"
	logFileName  = "server.log"
	commandUsage = `使用方法:
  server [command]

可用命令:
  start   启动服务
  stop    停止服务
  restart 重启服务
  run     直接运行(非守护进程模式)`
)

// 服务主函数 - 用于守护进程和直接运行模式
func serverMain() {
	bootstrap.Load()
	app.Run()
}

// 检查是否为调试模式 - 仅读取配置不连接数据库
func checkDebugMode() bool {
	bootstrap.LoadConf()
	return util.IsDebug()
}

func main() {
	if runtime.GOOS == "windows" || checkDebugMode() {
		if runtime.GOOS == "windows" {
			log.Println("Windows系统启动...")
		} else {
			log.Println("调试模式启动...")
		}
		serverMain()
		return
	}
	if len(os.Args) < 2 {
		log.Println(commandUsage)
		return
	}
	d, err := daemon.NewUnixDaemon(pidFileName, logFileName, serverMain)
	if err != nil {
		log.Printf("创建守护进程管理器失败: %v\n", err)
		os.Exit(1)
	}
	flag.Parse()
	command := os.Args[1]
	switch command {
	case "start":
		log.Println("正在启动服务...")
		if err := d.Start(); err != nil {
			log.Printf("启动失败: %v\n", err)
			os.Exit(1)
		}
		log.Println("服务已启动")
	case "stop":
		log.Println("正在停止服务...")
		if err := d.Stop(); err != nil {
			log.Printf("停止失败: %v\n", err)
			os.Exit(1)
		}
		log.Println("服务已停止")
	case "restart":
		log.Println("正在重启服务...")
		if err := d.Reload(); err != nil {
			log.Printf("重启失败: %v\n", err)
			os.Exit(1)
		}
		log.Println("服务已重启")
	case "run":
		log.Println("以前台模式运行服务...")
		serverMain()
	default:
		log.Println(commandUsage)
	}
}
