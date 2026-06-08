//go:build !linux

package server

import (
	"errors"
	"garnet/internal/config"
)

type Server struct{}

func NewServer(cfg *config.Config) *Server {
	return &Server{}
}

func (s *Server) Start() error {
	return errors.New("epoll server is only supported on Linux. Please run Garnet inside Docker")
}
