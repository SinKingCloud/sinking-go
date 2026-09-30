package command

import (
	"fmt"
	"log"

	"server/app/util/platform"
)

func (s *Server) install() error {
	if err := s.requireInteractive("install"); err != nil {
		return err
	}
	if err := platform.CheckPrivilege(); err != nil {
		return err
	}
	account, err := s.readAccount("请输入登录账号: ")
	if err != nil {
		return err
	}
	password, err := s.readNewPassword("请输入登录密码: ", "请再次输入登录密码: ")
	if err != nil {
		return err
	}
	defer clear(password)
	if err := s.updateAccount(account, string(password)); err != nil {
		return fmt.Errorf("设置登录账号密码失败: %w", err)
	}
	log.Println("登录账号密码设置成功")
	if err := s.daemon.InstallAutoStart(); err != nil {
		return fmt.Errorf("登录账号密码已设置，但安装系统自启动失败: %w", err)
	}
	log.Println("系统自启动安装成功")
	return s.restart()
}
