package daemon

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// AutoStartOptions 描述程序启动参数和系统自启动项。
// Executable、WorkingDirectory 默认使用当前程序和工作目录。
// Arguments 为 nil 时沿用当前进程参数，空切片表示不传参数；Name 仅管理系统自启动时必填。
type AutoStartOptions struct {
	Name               string
	Executable         string
	WorkingDirectory   string
	Arguments          []string
	KeepChildProcesses bool // Linux 服务退出时保留子进程，由应用决定如何清理。
}

// installAutoStart 安装并启用系统自启动。
func (u *Daemon) installAutoStart() error {
	options, err := u.normalizeOptions()
	if err != nil {
		return err
	}
	if options.Name == "" {
		return errors.New("自启动名称不能为空")
	}
	switch runtime.GOOS {
	case "linux":
		return u.installSystemd(options)
	case "windows":
		return u.installWindowsTask(options)
	case "darwin":
		return u.installLaunchAgent(options)
	default:
		return fmt.Errorf("当前系统不支持安装自启动: %s", runtime.GOOS)
	}
}

// uninstallAutoStart 删除系统自启动。
func (u *Daemon) uninstallAutoStart() error {
	options, err := u.normalizeOptions()
	if err != nil {
		return err
	}
	if options.Name == "" {
		return errors.New("自启动名称不能为空")
	}
	switch runtime.GOOS {
	case "linux":
		return u.uninstallSystemd(options)
	case "windows":
		return u.uninstallWindowsTask(options)
	case "darwin":
		return u.uninstallLaunchAgent(options)
	default:
		return fmt.Errorf("当前系统不支持删除自启动: %s", runtime.GOOS)
	}
}

func (u *Daemon) normalizeOptions() (AutoStartOptions, error) {
	options := u.autoStart
	var err error
	options.Name = strings.TrimSpace(options.Name)
	for index := range options.Name {
		value := options.Name[index]
		if (value >= 'a' && value <= 'z') || (value >= 'A' && value <= 'Z') ||
			(value >= '0' && value <= '9') || value == '-' || value == '_' || value == '.' {
			continue
		}
		return AutoStartOptions{}, errors.New("自启动名称只能包含字母、数字、点、横线和下划线")
	}
	if options.Executable == "" {
		options.Executable, err = os.Executable()
		if err != nil {
			return AutoStartOptions{}, fmt.Errorf("获取程序路径失败: %w", err)
		}
	}
	options.Executable, err = filepath.Abs(options.Executable)
	if err != nil {
		return AutoStartOptions{}, fmt.Errorf("解析程序路径失败: %w", err)
	}
	if options.WorkingDirectory == "" {
		options.WorkingDirectory, err = os.Getwd()
	} else {
		options.WorkingDirectory, err = filepath.Abs(options.WorkingDirectory)
	}
	if err != nil {
		return AutoStartOptions{}, fmt.Errorf("解析工作目录失败: %w", err)
	}
	if info, err := os.Stat(options.Executable); err != nil {
		return AutoStartOptions{}, fmt.Errorf("读取程序路径失败: %w", err)
	} else if info.IsDir() {
		return AutoStartOptions{}, errors.New("程序路径不能是目录")
	} else if runtime.GOOS != "windows" && info.Mode().Perm()&0111 == 0 {
		return AutoStartOptions{}, errors.New("程序文件没有执行权限")
	}
	if info, err := os.Stat(options.WorkingDirectory); err != nil {
		return AutoStartOptions{}, fmt.Errorf("读取工作目录失败: %w", err)
	} else if !info.IsDir() {
		return AutoStartOptions{}, errors.New("工作目录不是目录")
	}
	if options.Arguments == nil {
		options.Arguments = append([]string{}, os.Args[1:]...)
	} else {
		options.Arguments = append([]string{}, options.Arguments...)
	}
	return options, nil
}

func (u *Daemon) commandOutput(name string, args ...string) error {
	output, err := u.commandText(name, args...)
	if err != nil {
		message := strings.TrimSpace(string(output))
		if message != "" {
			return fmt.Errorf("执行 %s 失败: %w: %s", name, err, message)
		}
		return fmt.Errorf("执行 %s 失败: %w", name, err)
	}
	return nil
}

func (u *Daemon) commandText(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

func (u *Daemon) quoteSystemd(value string) string {
	value = strings.ReplaceAll(value, "%", "%%")
	value = strings.ReplaceAll(value, "$", "$$")
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `"`, `\"`)
	value = strings.ReplaceAll(value, "\n", `\n`)
	return `"` + value + `"`
}

func (u *Daemon) systemdUnitPath(name string) string {
	return filepath.Join("/etc/systemd/system", name+".service")
}

func (u *Daemon) installSystemd(options AutoStartOptions) error {
	if err := u.writeSystemd(options); err != nil {
		return err
	}
	return u.commandOutput("systemctl", "enable", options.Name+".service")
}

// writeSystemd 同步启动参数，不改变用户的开机自启开关。
func (u *Daemon) writeSystemd(options AutoStartOptions) error {
	if strings.ContainsAny(options.WorkingDirectory, "\r\n") {
		return errors.New("自启动工作目录不能包含换行符")
	}
	arguments := make([]string, 0, len(options.Arguments)+1)
	arguments = append(arguments, u.quoteSystemd(options.Executable))
	for _, argument := range options.Arguments {
		arguments = append(arguments, u.quoteSystemd(argument))
	}
	killMode := "mixed"
	if options.KeepChildProcesses {
		killMode = "process"
	}
	unit := fmt.Sprintf(`[Unit]
Description=%s
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
WorkingDirectory=%s
ExecStart=%s
Restart=on-failure
RestartSec=5
TimeoutStopSec=%d
KillMode=%s

[Install]
WantedBy=multi-user.target
`, options.Name, strings.ReplaceAll(options.WorkingDirectory, "%", "%%"), strings.Join(arguments, " "), int(u.StopTimeout.Seconds())+10, killMode)
	path := u.systemdUnitPath(options.Name)
	if err := os.WriteFile(path, []byte(unit), 0644); err != nil {
		return fmt.Errorf("写入 systemd 服务文件失败: %w", err)
	}
	return u.commandOutput("systemctl", "daemon-reload")
}

// syncSystemdKillMode 升级已安装服务的进程清理策略，保留原有启动参数与自启开关。
func (u *Daemon) syncSystemdKillMode() error {
	if !u.systemdInstalled() {
		return nil
	}
	path := u.systemdUnitPath(u.autoStart.Name)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取 systemd 服务文件失败: %w", err)
	}
	killMode := "mixed"
	if u.autoStart.KeepChildProcesses {
		killMode = "process"
	}
	lines := strings.Split(string(data), "\n")
	serviceLine := -1
	found := false
	for index, line := range lines {
		if strings.TrimSpace(line) == "[Service]" {
			serviceLine = index
		}
		if !strings.HasPrefix(strings.TrimSpace(line), "KillMode=") {
			continue
		}
		lines[index] = "KillMode=" + killMode
		found = true
		break
	}
	if !found {
		if serviceLine < 0 {
			return errors.New("systemd 服务文件缺少 Service 配置")
		}
		lines[serviceLine] += "\nKillMode=" + killMode
	}
	updated := []byte(strings.Join(lines, "\n"))
	if bytes.Equal(updated, data) {
		return nil
	}
	if err := os.WriteFile(path, updated, 0644); err != nil {
		return fmt.Errorf("更新 systemd 进程清理策略失败: %w", err)
	}
	if err := u.commandOutput("systemctl", "daemon-reload"); err != nil {
		// 未生效时还原文件，下一次操作仍会尝试同步。
		return errors.Join(err, os.WriteFile(path, data, 0644))
	}
	return nil
}

// systemdInstalled 只接管本程序安装的服务，未安装时继续使用普通守护进程。
func (u *Daemon) systemdInstalled() bool {
	if runtime.GOOS != "linux" || u.autoStart.Name == "" {
		return false
	}
	// 残留的 unit 文件不代表当前系统由 systemd 管理。
	if _, err := os.Stat("/run/systemd/system"); err != nil {
		return false
	}
	if _, err := exec.LookPath("systemctl"); err != nil {
		return false
	}
	_, err := os.Stat(u.systemdUnitPath(u.autoStart.Name))
	return err == nil
}

// SyncAutoStart 更新已安装的启动项，使手动启动与开机启动使用相同的参数。
func (u *Daemon) SyncAutoStart() error {
	if !u.systemdInstalled() {
		return nil
	}
	options, err := u.normalizeOptions()
	if err != nil {
		return err
	}
	return u.writeSystemd(options)
}

func (u *Daemon) uninstallSystemd(options AutoStartOptions) error {
	serviceName := options.Name + ".service"
	if _, err := os.Stat(u.systemdUnitPath(options.Name)); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err := u.commandOutput("systemctl", "disable", "--now", serviceName); err != nil {
		return err
	}
	if err := os.Remove(u.systemdUnitPath(options.Name)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("删除 systemd 服务文件失败: %w", err)
	}
	return u.commandOutput("systemctl", "daemon-reload")
}

func (u *Daemon) windowsTaskScriptPath(options AutoStartOptions) string {
	return filepath.Join(options.WorkingDirectory, "."+options.Name+"-startup.cmd")
}

func (u *Daemon) installWindowsTask(options AutoStartOptions) error {
	if err := u.writeWindowsTaskScript(options); err != nil {
		return err
	}
	scriptPath := u.windowsTaskScriptPath(options)
	// 计划任务默认运行 72 小时后停止；常驻服务必须禁用这一时限。
	task := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <Triggers><BootTrigger><Enabled>true</Enabled></BootTrigger></Triggers>
  <Principals><Principal id="Service"><UserId>S-1-5-18</UserId><RunLevel>HighestAvailable</RunLevel></Principal></Principals>
  <Settings><MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy><DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries><StopIfGoingOnBatteries>false</StopIfGoingOnBatteries><StartWhenAvailable>true</StartWhenAvailable><ExecutionTimeLimit>PT0S</ExecutionTimeLimit></Settings>
  <Actions Context="Service"><Exec><Command>cmd.exe</Command><Arguments>%s</Arguments><WorkingDirectory>%s</WorkingDirectory></Exec></Actions>
</Task>`, u.xmlText(`/d /s /c ""`+scriptPath+`""`), u.xmlText(options.WorkingDirectory))
	file, err := os.CreateTemp("", "daemon-startup-*.xml")
	if err != nil {
		return fmt.Errorf("创建 Windows 自启动任务配置失败: %w", err)
	}
	defer os.Remove(file.Name())
	_, err = file.WriteString(task)
	closeErr := file.Close()
	if err != nil {
		return fmt.Errorf("写入 Windows 自启动任务配置失败: %w", err)
	}
	if closeErr != nil {
		return fmt.Errorf("保存 Windows 自启动任务配置失败: %w", closeErr)
	}
	if err := u.commandOutput("schtasks.exe", "/Create", "/TN", options.Name, "/XML", file.Name(), "/F"); err != nil {
		_ = os.Remove(scriptPath)
		return err
	}
	return nil
}

// writeWindowsTaskScript 保存启动参数，供安装和带参数启动时共用。
func (u *Daemon) writeWindowsTaskScript(options AutoStartOptions) error {
	if strings.ContainsAny(options.WorkingDirectory, "\r\n\"") {
		return errors.New("Windows自启动路径包含不支持的字符")
	}
	arguments := make([]string, 0, len(options.Arguments)+1)
	for _, argument := range append([]string{options.Executable}, options.Arguments...) {
		if strings.ContainsAny(argument, "\r\n\"") {
			return errors.New("Windows自启动参数包含不支持的字符")
		}
		argument = strings.ReplaceAll(argument, "%", "%%")
		argument += strings.Repeat(`\`, len(argument)-len(strings.TrimRight(argument, `\`)))
		arguments = append(arguments, `"`+argument+`"`)
	}
	logPath := u.LogFileName
	if !filepath.IsAbs(logPath) {
		logPath = filepath.Join(options.WorkingDirectory, logPath)
	}
	if strings.ContainsAny(logPath, "\r\n\"") {
		return errors.New("Windows 日志路径包含不支持的字符")
	}
	script := fmt.Sprintf("@echo off\r\nsetlocal DisableDelayedExpansion\r\ncd /d \"%s\" || exit /b 1\r\n%s >>\"%s\" 2>&1\r\n",
		strings.ReplaceAll(options.WorkingDirectory, "%", "%%"), strings.Join(arguments, " "), strings.ReplaceAll(logPath, "%", "%%"))
	if err := os.WriteFile(u.windowsTaskScriptPath(options), []byte(script), 0644); err != nil {
		return fmt.Errorf("写入 Windows 自启动脚本失败: %w", err)
	}
	return nil
}

func (u *Daemon) uninstallWindowsTask(options AutoStartOptions) error {
	if err := u.commandOutput("schtasks.exe", "/Query", "/TN", options.Name); err != nil {
		return nil
	}
	if err := u.commandOutput("schtasks.exe", "/Delete", "/TN", options.Name, "/F"); err != nil {
		return err
	}
	if err := os.Remove(u.windowsTaskScriptPath(options)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("删除 Windows 自启动脚本失败: %w", err)
	}
	return nil
}

func (u *Daemon) launchAgentPath(options AutoStartOptions) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户目录失败: %w", err)
	}
	return filepath.Join(home, "Library", "LaunchAgents", options.Name+".plist"), nil
}

func (u *Daemon) launchAgentLogPath(options AutoStartOptions) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户目录失败: %w", err)
	}
	return filepath.Join(home, "Library", "Logs", options.Name+".log"), nil
}

func (u *Daemon) xmlText(value string) string {
	var output bytes.Buffer
	_ = xml.EscapeText(&output, []byte(value))
	return output.String()
}

func (u *Daemon) installLaunchAgent(options AutoStartOptions) error {
	plistPath, err := u.launchAgentPath(options)
	if err != nil {
		return err
	}
	logPath, err := u.launchAgentLogPath(options)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(plistPath), 0755); err != nil {
		return fmt.Errorf("创建 LaunchAgents 目录失败: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return fmt.Errorf("创建 LaunchAgent 日志目录失败: %w", err)
	}
	var content strings.Builder
	content.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>`)
	content.WriteString("<key>Label</key><string>" + u.xmlText(options.Name) + "</string>")
	content.WriteString("<key>ProgramArguments</key><array>")
	content.WriteString("<string>" + u.xmlText(options.Executable) + "</string>")
	for _, argument := range options.Arguments {
		content.WriteString("<string>" + u.xmlText(argument) + "</string>")
	}
	content.WriteString("</array>")
	content.WriteString("<key>WorkingDirectory</key><string>" + u.xmlText(options.WorkingDirectory) + "</string>")
	content.WriteString("<key>RunAtLoad</key><true/>")
	content.WriteString("<key>KeepAlive</key><true/>")
	content.WriteString("<key>StandardOutPath</key><string>" + u.xmlText(logPath) + "</string>")
	content.WriteString("<key>StandardErrorPath</key><string>" + u.xmlText(logPath) + "</string>")
	content.WriteString("</dict></plist>\n")
	if err := os.WriteFile(plistPath, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("写入 LaunchAgent 文件失败: %w", err)
	}
	return nil
}

func (u *Daemon) uninstallLaunchAgent(options AutoStartOptions) error {
	path, err := u.launchAgentPath(options)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("删除 LaunchAgent 文件失败: %w", err)
	}
	return nil
}
