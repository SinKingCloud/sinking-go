package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/SinKingCloud/sinking-go/sinking-web"
)

// Server HTTP服务器
type Server struct {
	addr        string
	debug       bool
	engine      *sinking_web.Engine
	server      *http.Server
	done        chan error
	running     bool
	mutex       sync.RWMutex
	handlers    []func(engine *sinking_web.Engine)
	errorHandle *sinking_web.ErrorHandel
}

// NewServer 创建新的HTTP服务器实例
func NewServer(addr string, debug bool) *Server {
	return &Server{
		addr:     addr,
		debug:    debug,
		running:  false,
		handlers: make([]func(engine *sinking_web.Engine), 0),
	}
}

// init 初始化服务器
func (s *Server) init() {
	if s.engine != nil {
		return
	}
	s.engine = sinking_web.Default()
	s.engine.SetDebugMode(s.debug)
	if s.errorHandle != nil {
		s.engine.SetErrorHandle(s.errorHandle)
	}
	for _, handler := range s.handlers {
		handler(s.engine)
	}
}

// Handle 添加处理器函数
func (s *Server) Handle(handler func(engine *sinking_web.Engine)) *Server {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.handlers = append(s.handlers, handler)
	if s.engine != nil {
		handler(s.engine)
	}
	return s
}

// ErrorHandle 设置错误处理回调
func (s *Server) ErrorHandle(handle *sinking_web.ErrorHandel) {
	s.errorHandle = handle
}

// Start 启动服务器
func (s *Server) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.running {
		return fmt.Errorf("服务器已经在运行中")
	}
	s.init()
	if s.addr == "" {
		s.addr = ":5678"
	}
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("监听HTTP地址失败: %w", err)
	}
	server := &http.Server{
		Addr:    s.addr,
		Handler: s.engine,
	}
	done := make(chan error, 1)
	s.server = server
	s.done = done
	s.running = true
	go func() {
		err := server.Serve(listener)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		s.mutex.Lock()
		if s.server == server {
			s.running = false
		}
		s.mutex.Unlock()
		done <- err
		close(done)
	}()
	return nil
}

// Done 返回当前HTTP服务的退出结果。
func (s *Server) Done() <-chan error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.done
}

// Stop 停止服务器
func (s *Server) Stop() error {
	s.mutex.Lock()
	if !s.running {
		s.mutex.Unlock()
		return fmt.Errorf("服务器未运行")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.mutex.Unlock()
		return err
	}
	done := s.done
	s.running = false
	s.server = nil
	s.mutex.Unlock()
	return <-done
}

// Restart 重启服务器
func (s *Server) Restart() error {
	s.mutex.RLock()
	running := s.running
	s.mutex.RUnlock()
	if running {
		if err := s.Stop(); err != nil {
			return fmt.Errorf("停止服务器失败: %v", err)
		}
	}
	time.Sleep(500 * time.Millisecond)
	if err := s.Start(); err != nil {
		return fmt.Errorf("启动服务器失败: %v", err)
	}
	return nil
}
