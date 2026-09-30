package daemon

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	goDaemon "github.com/sevlyar/go-daemon"
)

// Restart 请求当前服务自行重启，调用成功表示已提交重启。
// 旧进程退出后按配置或原始参数启动程序，不要求程序实现特定子命令。
func (u *Daemon) Restart() (result error) {
	if !u.restarting.CompareAndSwap(false, true) {
		return errors.New("程序正在重启，请勿重复操作")
	}
	defer func() {
		if result != nil {
			u.restarting.Store(false)
		}
	}()
	options, err := u.normalizeOptions()
	if err != nil {
		return err
	}
	pid, err := u.readPidFile()
	if err != nil {
		return fmt.Errorf("读取当前服务 PID 失败: %w", err)
	}
	if pid != os.Getpid() {
		return errors.New("Restart 只能在当前运行的服务中调用，外部重启请使用 Reload")
	}
	if u.systemdInstalled() {
		if err := u.syncSystemdKillMode(); err != nil {
			return err
		}
		output, err := u.commandText("systemctl", "show", "--property=MainPID", "--value", options.Name+".service")
		if err != nil {
			return fmt.Errorf("查询系统服务失败: %w", err)
		}
		pid, err := strconv.Atoi(strings.TrimSpace(string(output)))
		if err != nil {
			return fmt.Errorf("读取系统服务 PID 失败: %w", err)
		}
		if pid == os.Getpid() {
			return u.commandOutput("systemctl", "--no-block", "restart", options.Name+".service")
		}
		if pid > 0 {
			return fmt.Errorf("同名系统服务正在运行其他进程: %d", pid)
		}
	}
	var command *exec.Cmd
	if runtime.GOOS == "darwin" && options.Name != "" {
		uid := strconv.Itoa(os.Getuid())
		for _, domain := range []string{"gui/" + uid, "user/" + uid, "system"} {
			service := domain + "/" + options.Name
			output, err := u.commandText("launchctl", "print", service)
			if err != nil {
				continue
			}
			for _, line := range strings.Split(string(output), "\n") {
				// 只操作当前进程所属服务，不影响其他用户或目录中的同名实例。
				if strings.TrimSpace(line) != "pid = "+strconv.Itoa(os.Getpid()) {
					continue
				}
				// 先平滑退出，再启动；兼容未配置 KeepAlive 的 LaunchDaemon。
				timeout := u.StopTimeout
				if timeout <= 0 {
					timeout = defaultStopTimeout
				}
				command = exec.Command("/bin/sh", "-c", `launchctl kill SIGTERM "$1" || exit
attempt=0
while kill -0 "$2" 2>/dev/null; do
    if [ "$attempt" -ge "$3" ]; then echo "等待服务停止超时" >&2; exit 1; fi
    sleep 1
    attempt=$((attempt + 1))
done
exec launchctl kickstart "$1"`, "daemon-restart", service, strconv.Itoa(os.Getpid()), strconv.FormatInt(int64((timeout+time.Second-1)/time.Second), 10))
				break
			}
			if command != nil {
				break
			}
		}
	}
	if command == nil {
		command, err = u.restartCommand(options)
		if err != nil {
			return err
		}
	}
	command.Dir = options.WorkingDirectory
	for _, value := range os.Environ() {
		// 控制进程不是 go-daemon 子进程，不能继承它的重生标记。
		if !strings.HasPrefix(value, goDaemon.MARK_NAME+"=") {
			command.Env = append(command.Env, value)
		}
	}
	logPath := u.LogFileName
	if !filepath.IsAbs(logPath) {
		logPath = filepath.Join(options.WorkingDirectory, logPath)
	}
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		return fmt.Errorf("打开重启日志失败: %w", err)
	}
	defer file.Close()
	command.Stdout, command.Stderr = file, file
	command, err = u.detachRestart(command)
	if err != nil {
		return fmt.Errorf("启动重启控制进程失败: %w", err)
	}
	// 控制进程需要等待当前进程退出，不能在这里同步 Wait。
	go func() {
		if err := command.Wait(); err != nil {
			u.restarting.Store(false)
			log.Printf("重启控制进程执行失败: %v", err)
		}
	}()
	return nil
}

// restartCommand 先平滑停止当前进程，确认退出后按原始参数启动程序。
func (u *Daemon) restartCommand(options AutoStartOptions) (*exec.Cmd, error) {
	timeout := u.StopTimeout
	if timeout <= 0 {
		timeout = defaultStopTimeout
	}
	if runtime.GOOS != "windows" {
		arguments := []string{"-c", `previous="$1"
timeout="$2"
shift 2
kill -TERM "$previous" || exit 1
attempt=0
while kill -0 "$previous" 2>/dev/null; do
    if [ "$attempt" -ge "$timeout" ]; then echo "等待服务停止超时" >&2; exit 1; fi
    sleep 1
    attempt=$((attempt + 1))
done
exec "$@"`, "daemon-restart", strconv.Itoa(os.Getpid()), strconv.FormatInt(int64((timeout+time.Second-1)/time.Second), 10), options.Executable}
		arguments = append(arguments, options.Arguments...)
		return exec.Command("/bin/sh", arguments...), nil
	}
	systemRoot := os.Getenv("SystemRoot")
	if !filepath.IsAbs(systemRoot) {
		return nil, errors.New("无法获取 Windows 系统目录")
	}
	stopFile, err := filepath.Abs(u.PidFileName + ".stop")
	if err != nil {
		return nil, fmt.Errorf("解析服务停止通知路径失败: %w", err)
	}
	// 按 Windows 命令行规则转义，保留空参数、引号和末尾反斜杠。
	arguments := make([]string, len(options.Arguments))
	for index, argument := range options.Arguments {
		var escaped strings.Builder
		escaped.WriteByte('"')
		slashes := 0
		for _, value := range argument {
			if value == '\\' {
				slashes++
				continue
			}
			if value == '"' {
				escaped.WriteString(strings.Repeat(`\`, slashes*2+1))
			} else {
				escaped.WriteString(strings.Repeat(`\`, slashes))
			}
			escaped.WriteRune(value)
			slashes = 0
		}
		escaped.WriteString(strings.Repeat(`\`, slashes*2))
		escaped.WriteByte('"')
		arguments[index] = escaped.String()
	}
	quote := func(value string) string {
		return "'" + strings.ReplaceAll(value, "'", "''") + "'"
	}
	// 先持有旧进程句柄再通知退出，避免 PID 被复用；超时不启动第二个实例。
	// 仅重定向标准输入，使 .NET 明确继承控制进程的两个日志句柄。
	script := fmt.Sprintf(`$ErrorActionPreference = 'Stop'
try {
    $previous = [System.Diagnostics.Process]::GetProcessById(%d)
    try {
        $null = $previous.Handle
        [System.IO.File]::WriteAllText(%s, '%d')
        if (-not $previous.WaitForExit(%d)) { throw '等待服务停止超时' }
    } finally {
        $previous.Dispose()
    }
    $start = New-Object System.Diagnostics.ProcessStartInfo
    $start.FileName = %s
    $start.Arguments = %s
    $start.WorkingDirectory = %s
    $start.UseShellExecute = $false
    $start.CreateNoWindow = $true
    $start.RedirectStandardInput = $true
    $next = [System.Diagnostics.Process]::Start($start)
    try {
        $next.StandardInput.Close()
    } finally {
        $next.Dispose()
    }
} catch {
    [Console]::Error.WriteLine($_.Exception.Message)
    exit 1
}`, os.Getpid(), quote(stopFile), os.Getpid(), min(timeout.Milliseconds(), int64(1<<31-1)),
		quote(options.Executable), quote(strings.Join(arguments, " ")), quote(options.WorkingDirectory))
	return exec.Command(filepath.Join(systemRoot, "System32", "WindowsPowerShell", "v1.0", "powershell.exe"),
		"-NoProfile", "-NonInteractive", "-Command", script), nil
}

// detachRestart 根据系统独立启动进程；反射处理不同平台的 SysProcAttr 字段。
func (u *Daemon) detachRestart(command *exec.Cmd) (*exec.Cmd, error) {
	if command.SysProcAttr == nil {
		command.SysProcAttr = &syscall.SysProcAttr{}
	}
	attributes := reflect.ValueOf(command.SysProcAttr).Elem()
	if runtime.GOOS != "windows" {
		field := attributes.FieldByName("Setsid")
		if !field.IsValid() || !field.CanSet() || field.Kind() != reflect.Bool {
			return nil, fmt.Errorf("当前系统不支持独立启动进程: %s", runtime.GOOS)
		}
		field.SetBool(true)
		return command, command.Start()
	}
	const (
		detachedProcess        = 0x00000008
		createNewProcessGroup  = 0x00000200
		createBreakawayFromJob = 0x01000000
	)
	flags := attributes.FieldByName("CreationFlags")
	if !flags.IsValid() || !flags.CanSet() || flags.Kind() != reflect.Uint32 {
		return nil, errors.New("当前系统不支持设置 Windows 进程属性")
	}
	flags.SetUint(flags.Uint() | detachedProcess | createNewProcessGroup | createBreakawayFromJob)
	err := command.Start()
	if errors.Is(err, syscall.Errno(5)) { // ERROR_ACCESS_DENIED：所属 Job 可能禁止脱离。
		flags.SetUint(flags.Uint() &^ createBreakawayFromJob)
		// Start 失败的 Cmd 也不能再次使用；重新创建，保留参数、环境和日志句柄。
		command = &exec.Cmd{
			Path: command.Path, Args: command.Args, Dir: command.Dir, Env: command.Env,
			Stdin: command.Stdin, Stdout: command.Stdout, Stderr: command.Stderr,
			SysProcAttr: command.SysProcAttr,
		}
		err = command.Start()
	}
	return command, err
}
