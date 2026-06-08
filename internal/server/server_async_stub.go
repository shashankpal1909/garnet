//go:build !linux

package server

import (
	"errors"
	"garnet/internal/config"
)

type AsyncServer struct{}

func NewAsyncServer(cfg *config.Config) *AsyncServer {
	return &AsyncServer{}
}

func (s *AsyncServer) Start() error {
	return errors.New("async epoll server is only supported on Linux. Please change config mode to 'sync' or run Garnet inside Docker")
}
