package daemon

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	goDaemon "github.com/sevlyar/go-daemon"
)

func (u *Daemon) runWindows() (result error) {
	signalContext, stopSignal := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stopSignal()
	stopFile := u.PidFileName + ".stop"
	_ = os.Remove(stopFile)
	if err := os.WriteFile(u.PidFileName, []byte(strconv.Itoa(os.Getpid())), 0600); err != nil {
		return fmt.Errorf("写入服务 PID 失败: %w", err)
	}
	stop := make(chan struct{})
	done := make(chan struct{})
	group := sync.WaitGroup{}
	group.Add(1)
	go func() {
		defer group.Done()
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-signalContext.Done():
				close(stop)
				return
			case <-ticker.C:
				if data, err := os.ReadFile(stopFile); err == nil && strings.TrimSpace(string(data)) == strconv.Itoa(os.Getpid()) {
					close(stop)
					return
				}
			}
		}
	}()
	defer func() {
		close(done)
		group.Wait()
		_ = os.Remove(stopFile)
		if err := os.Remove(u.PidFileName); result == nil && err != nil && !errors.Is(err, os.ErrNotExist) {
			result = fmt.Errorf("删除服务 PID 失败: %w", err)
		}
	}()
	u.Service(stop)
	return nil
}

// startWindows 独立启动后台服务；计划任务仅负责开机启动，避免任务尚未退出时 Run 无效。
func (u *Daemon) startWindows() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("当前系统不是 Windows")
	}
	options, err := u.normalizeOptions()
	if err != nil {
		return err
	}
	if pid, err := u.readPidFile(); err == nil && pid > 1 {
		if process, err := os.FindProcess(pid); err == nil {
			_ = process.Release()
			return errors.New("服务已经运行")
		}
	}
	if options.Name != "" {
		if _, err := os.Stat(u.windowsTaskScriptPath(options)); err == nil {
			if err := u.writeWindowsTaskScript(options); err != nil {
				return err
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("检查 Windows 自启动脚本失败: %w", err)
		}
	}
	command := exec.Command(options.Executable, options.Arguments...)
	command.Dir = options.WorkingDirectory
	// 无需约定 run 子命令，再次进入 Start 时直接运行服务。
	command.Env = append(os.Environ(), goDaemon.MARK_NAME+"="+goDaemon.MARK_VALUE)
	logPath := u.LogFileName
	if !filepath.IsAbs(logPath) {
		logPath = filepath.Join(options.WorkingDirectory, logPath)
	}
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		return fmt.Errorf("打开服务日志失败: %w", err)
	}
	defer file.Close()
	command.Stdout, command.Stderr = file, file
	command, err = u.detachRestart(command)
	if err != nil {
		return fmt.Errorf("启动 Windows 服务失败: %w", err)
	}
	_ = command.Process.Release()
	return nil
}

// stopWindows 通知服务进程自行平滑退出。
func (u *Daemon) stopWindows(process *os.Process) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("当前系统不是 Windows")
	}
	stopFile := u.PidFileName + ".stop"
	if err := os.WriteFile(stopFile, []byte(strconv.Itoa(process.Pid)), 0600); err != nil {
		return fmt.Errorf("发送停止通知失败: %w", err)
	}
	deadline := time.NewTimer(u.StopTimeout)
	defer deadline.Stop()
	stopped := make(chan error, 1)
	go func() {
		// Windows 可等待非子进程，确保进程已退出而不只是刚删除 PID 文件。
		_, err := process.Wait()
		stopped <- err
	}()
	select {
	case <-deadline.C:
		return errors.New("等待服务停止超时")
	case err := <-stopped:
		if err != nil {
			return fmt.Errorf("等待服务退出失败: %w", err)
		}
		_ = os.Remove(stopFile)
		return nil
	}
}
