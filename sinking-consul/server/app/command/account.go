package command

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"server/app/constant"
	"server/app/service/auth"
	"server/app/service/setting"
	"server/app/util/validator"
	"server/bootstrap"
	"server/global"

	"golang.org/x/term"
)

func (s *Server) user() error {
	if err := s.requireInteractive("user"); err != nil {
		return err
	}
	account, err := s.readAccount("请输入新账号: ")
	if err != nil {
		return err
	}
	if err := s.updateAccount(account, ""); err != nil {
		return fmt.Errorf("重置登录账号失败: %w", err)
	}
	log.Println("登录账号重置成功，已清除保存的登录会话；服务运行中时，请执行 restart 使修改立即生效")
	if global.App.Config.GetString(constant.AuthPassword) == "" {
		log.Printf("登录密码尚未设置，请继续执行 %s pwd", os.Args[0])
	}
	return nil
}

func (s *Server) pwd() error {
	if err := s.requireInteractive("pwd"); err != nil {
		return err
	}
	password, err := s.readNewPassword("请输入新密码: ", "请再次输入新密码: ")
	if err != nil {
		return err
	}
	defer clear(password)
	if err := s.updateAccount("", string(password)); err != nil {
		return fmt.Errorf("重置登录密码失败: %w", err)
	}
	log.Println("登录密码重置成功，已清除保存的登录会话；服务运行中时，请执行 restart 使修改立即生效")
	if global.App.Config.GetString(constant.AuthAccount) == "" {
		log.Printf("登录账号尚未设置，请继续执行 %s user", os.Args[0])
	}
	return nil
}

func (s *Server) updateAccount(account string, password string) error {
	bootstrap.LoadConf()
	config := setting.NewService(global.App.Config)
	return auth.NewService(config, nil).UpdateAccount(account, password)
}

func (s *Server) requireInteractive(command string) error {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("请在交互式终端中执行 %s 命令", command)
	}
	return nil
}

func (s *Server) readAccount(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	account, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("读取账号失败: %w", err)
	}
	if errors.Is(err, io.EOF) && account == "" {
		return "", errors.New("已取消输入")
	}
	form := struct {
		Account string `validate:"required,alphanum,min=3,max=20" label:"登录账号"`
	}{Account: strings.TrimSpace(account)}
	if ok, message := validator.Check(&form); !ok {
		return "", errors.New(message)
	}
	return form.Account, nil
}

func (s *Server) readNewPassword(prompt string, confirmPrompt string) ([]byte, error) {
	password, err := s.readPassword(prompt)
	if err != nil {
		return nil, err
	}
	confirm, err := s.readPassword(confirmPrompt)
	if err != nil {
		clear(password)
		return nil, err
	}
	defer clear(confirm)
	if !bytes.Equal(password, confirm) {
		clear(password)
		return nil, errors.New("两次输入的密码不一致")
	}
	if !utf8.Valid(password) {
		clear(password)
		return nil, errors.New("登录密码格式不合法")
	}
	length := utf8.RuneCount(password)
	if length < 6 || length > 20 {
		clear(password)
		return nil, errors.New("登录密码长度必须在6到20个字符之间")
	}
	if len(password) > 72 {
		clear(password)
		return nil, errors.New("登录密码不能超过72个字节")
	}
	return password, nil
}

func (s *Server) readPassword(prompt string) ([]byte, error) {
	fmt.Fprint(os.Stderr, prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		clear(password)
		return nil, fmt.Errorf("读取密码失败: %w", err)
	}
	return password, nil
}
