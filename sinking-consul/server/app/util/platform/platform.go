package platform

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// CheckPrivilege 按系统检查服务安装和卸载权限。
func CheckPrivilege() error {
	switch runtime.GOOS {
	case "linux":
		if os.Geteuid() != 0 {
			return errors.New("Linux 系统需要 root 权限，请使用 root 用户或 sudo 运行此命令")
		}
	case "windows":
		systemRoot := os.Getenv("SystemRoot")
		if systemRoot == "" {
			return errors.New("无法读取 Windows 系统目录，无法检查管理员权限")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		path := filepath.Join(systemRoot, "System32", "WindowsPowerShell", "v1.0", "powershell.exe")
		script := `$ErrorActionPreference = 'Stop'; $principal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent()); if ($principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) { exit 0 }; exit 5`
		if err := exec.CommandContext(ctx, path, "-NoLogo", "-NoProfile", "-NonInteractive", "-Command", script).Run(); err != nil {
			if ctx.Err() != nil {
				return fmt.Errorf("检查 Windows 管理员权限超时: %w", ctx.Err())
			}
			var exitError *exec.ExitError
			if errors.As(err, &exitError) && exitError.ExitCode() == 5 {
				return errors.New("Windows 系统需要管理员权限，请以管理员身份运行终端后执行此命令")
			}
			return fmt.Errorf("通过系统 PowerShell 检查 Windows 管理员权限失败: %w", err)
		}
	}
	return nil
}
