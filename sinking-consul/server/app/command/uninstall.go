package command

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"server/app/util/platform"
)

func readUninstallConfirmation(input io.Reader, output io.Writer) (bool, error) {
	fmt.Fprint(output, "将停止服务并移除系统自启动，程序、配置和数据库会保留，输入 yes 确认卸载: ")
	confirmation, err := bufio.NewReader(input).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, fmt.Errorf("读取卸载确认失败: %w", err)
	}
	if strings.TrimSpace(confirmation) != "yes" {
		fmt.Fprintln(output, "已取消卸载")
		return false, nil
	}
	return true, nil
}

func (s *Server) uninstall() error {
	confirmed, err := readUninstallConfirmation(os.Stdin, os.Stdout)
	if err != nil || !confirmed {
		return err
	}
	if err := platform.CheckPrivilege(); err != nil {
		return err
	}
	log.Println("正在卸载系统自启动...")
	if err := s.daemon.Uninstall(); err != nil {
		return err
	}
	log.Println("系统自启动已卸载，程序、配置和数据库已保留")
	return nil
}
