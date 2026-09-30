package command

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"server/app"
	"server/app/util/daemon"
	"server/bootstrap"
	"server/global"
)

const (
	pidFileName = "server.pid"
	logFileName = "server.log"
	serviceName = "sinking-consul"
	usage       = `使用方法:
  server [command] [options]

可用命令:
  start   启动服务
  stop    停止服务
  restart 重启服务
  run     直接运行（非守护进程模式）
  install 设置登录账号密码、安装系统自启动并启动服务
  uninstall 停止服务并移除系统自启动，保留程序、配置和数据库
  user    重置登录账号
  pwd     重置登录密码

  -h, --help 显示帮助

登录信息保存在当前工作目录的 config/application.yml。`
)

// Server 管理服务运行命令。
type Server struct {
	daemon *daemon.Daemon
}

// Run 是服务命令的统一入口。
func Run(args []string) error {
	server, err := NewServer()
	if err != nil {
		return err
	}
	return server.Execute(args)
}

// NewServer 创建服务命令。
func NewServer() (*Server, error) {
	server := &Server{}
	manager, err := daemon.NewDaemon(pidFileName, logFileName, server.run)
	if err != nil {
		return nil, fmt.Errorf("创建守护进程管理器失败: %w", err)
	}
	manager.SetChildArgs(os.Args[0], "start")
	manager.SetAutoStartOptions(daemon.AutoStartOptions{
		Name:      serviceName,
		Arguments: []string{"run"},
	})
	server.daemon = manager
	global.App.SetDaemon(manager)
	return server, nil
}

// Execute 执行服务命令。
func (s *Server) Execute(args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		log.Println(usage)
		return nil
	}
	command := args[0]
	switch command {
	case "run", "start", "stop", "restart", "install", "uninstall", "user", "pwd":
	default:
		return fmt.Errorf("未知命令: %s\n%s", command, usage)
	}
	options := flag.NewFlagSet(command, flag.ContinueOnError)
	options.SetOutput(io.Discard)
	if err := options.Parse(args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			log.Println(usage)
			return nil
		}
		return fmt.Errorf("命令参数错误: %w", err)
	}
	if options.NArg() > 0 {
		return fmt.Errorf("无法识别的命令参数: %s", strings.Join(options.Args(), " "))
	}
	switch command {
	case "run":
		log.Println("以前台模式运行服务...")
		return s.daemon.Run()
	case "install":
		return s.install()
	case "uninstall":
		return s.uninstall()
	case "user":
		return s.user()
	case "pwd":
		return s.pwd()
	case "start":
		return s.start()
	case "stop":
		return s.stop()
	case "restart":
		return s.restart()
	}
	return nil
}

func (s *Server) run(stop <-chan struct{}) {
	bootstrap.Load()
	app.Run(stop)
}

func (s *Server) start() error {
	log.Println("正在启动服务...")
	if err := s.daemon.Start(); err != nil {
		return fmt.Errorf("启动失败: %w", err)
	}
	if !s.daemon.IsChildProcess() {
		log.Println("服务已启动")
	}
	return nil
}

func (s *Server) stop() error {
	log.Println("正在停止服务...")
	if err := s.daemon.Stop(); err != nil {
		return fmt.Errorf("停止失败: %w", err)
	}
	log.Println("服务已停止")
	return nil
}

func (s *Server) restart() error {
	log.Println("正在重启服务...")
	if err := s.daemon.Reload(); err != nil {
		return fmt.Errorf("重启失败: %w", err)
	}
	if !s.daemon.IsChildProcess() {
		log.Println("服务已重启")
	}
	return nil
}
