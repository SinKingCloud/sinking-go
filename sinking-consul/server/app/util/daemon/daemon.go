package daemon

import (
	"errors"
	"fmt"
	"github.com/sevlyar/go-daemon"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type UnixDaemon struct {
	PidFileName string
	LogFileName string
	Service     func()
}

// NewUnixDaemon 实例化进程守护
func NewUnixDaemon(pidFileName string, logFileName string, service func()) (*UnixDaemon, error) {
	if pidFileName == "" || logFileName == "" || service == nil {
		return nil, errors.New("参数不能为空")
	}
	return &UnixDaemon{
		PidFileName: pidFileName,
		LogFileName: logFileName,
		Service:     service,
	}, nil
}

// Start 启动
func (u *UnixDaemon) Start() error {
	// 创建守护进程上下文
	daemonCtx := &daemon.Context{
		PidFileName: u.PidFileName,
		LogFileName: u.LogFileName,
		WorkDir:     "./",
		Umask:       027,
	}
	// 启动守护进程并获取新的进程上下文和PID
	d, err := daemonCtx.Reborn()
	if err != nil {
		return fmt.Errorf("启动失败: %v", err)
	}
	// 父进程，已经启动了守护进程，直接返回
	if d != nil {
		return nil
	}
	// 子进程，这里是守护进程的实际执行逻辑
	defer func() {
		_ = daemonCtx.Release()
	}()
	go u.Service()
	// 等待信号，以便能够正确处理停止等操作
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	return nil
}

// Stop 停止
func (u *UnixDaemon) Stop() error {
	// 读取PID文件获取守护进程的PID
	pid, err := u.readPidFile()
	if err != nil {
		return fmt.Errorf("无法读取PID文件: %v", err)
	}
	// 发送终止信号给守护进程
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		return fmt.Errorf("无法停止守护进程: %v", err)
	}
	// 删除PID文件
	if err := os.Remove("server.pid"); err != nil {
		return fmt.Errorf("无法删除PID文件: %v", err)
	}
	return nil
}

// Reload 重启
func (u *UnixDaemon) Reload() error {
	// 先停止守护进程
	_ = u.Stop()
	time.Sleep(time.Second)
	// 再启动守护进程
	if err := u.Start(); err != nil {
		return fmt.Errorf("无法重新启动守护进程: %v", err)
	}
	return nil
}

// readPidFile 读取pid
func (u *UnixDaemon) readPidFile() (int, error) {
	data, err := os.ReadFile(u.PidFileName)
	if err != nil {
		return 0, err
	}
	var pid int
	_, err = fmt.Sscanf(string(data), "%d", &pid)
	if err != nil {
		return 0, err
	}
	return pid, nil
}
