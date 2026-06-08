//go:build linux

package server

import (
	"fmt"
	"net"
	"syscall"

	"garnet/internal/config"
	"garnet/internal/logger"
)

type Server struct {
	cfg         *config.Config
	fd          int
	epfd        int
	connections map[int]*Connection
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:         cfg,
		connections: make(map[int]*Connection),
	}
}

func (s *Server) Start() error {
	logger.Logger.Print(Banner)
	logger.Logger.Printf("starting epoll server on %s:%d", s.cfg.Host, s.cfg.Port)

	// 1. Create socket
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("socket error: %v", err)
	}
	s.fd = fd
	defer syscall.Close(fd)

	// Set SO_REUSEADDR
	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return fmt.Errorf("setsockopt error: %v", err)
	}

	// 2. Bind
	ip := net.ParseIP(s.cfg.Host)
	var addr [4]byte
	copy(addr[:], ip.To4())
	sa := &syscall.SockaddrInet4{
		Port: s.cfg.Port,
		Addr: addr,
	}

	if err := syscall.Bind(fd, sa); err != nil {
		return fmt.Errorf("bind error: %v", err)
	}

	// 3. Listen
	if err := syscall.Listen(fd, 1024); err != nil {
		return fmt.Errorf("listen error: %v", err)
	}

	// 4. Set non-blocking
	if err := syscall.SetNonblock(fd, true); err != nil {
		return fmt.Errorf("setnonblock error: %v", err)
	}

	return s.startEventLoop()
}
