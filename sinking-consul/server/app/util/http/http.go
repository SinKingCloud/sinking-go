package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/SinKingCloud/sinking-go/sinking-web"
	"net/http"
	"sync"
	"time"
)

// Server HTTP服务器
type Server struct {
	addr        string
	debug       bool
	engine      *sinking_web.Engine
	server      *http.Server
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
	s.server = &http.Server{
		Addr:    s.addr,
		Handler: s.engine,
	}
	s.running = true
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.mutex.Lock()
			s.running = false
			s.mutex.Unlock()
		}
	}()
	return nil
}

// Stop 停止服务器
func (s *Server) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.running {
		return fmt.Errorf("服务器未运行")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	s.running = false
	s.server = nil
	return nil
}

// Restart 重启服务器
func (s *Server) Restart() error {
	if s.running {
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
