package server

import (
	"fmt"
	"net"

	"garnet/internal/config"
	"garnet/internal/logger"
)

// Server manages the TCP listener and
// incoming client connections.
type Server struct {
	cfg *config.Config
}

// New creates a new server instance.
func New(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

// Start begins listening for TCP connections
// and blocks until the server terminates.
func (s *Server) Start() error {
	addr := fmt.Sprintf(
		"%s:%d",
		s.cfg.Host,
		s.cfg.Port,
	)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	logger.Logger.Printf(
		"server started on %s",
		addr,
	)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Logger.Printf(
				"accept connection failed: %v",
				err,
			)
			continue
		}

		go HandleConnection(conn)
	}
}
