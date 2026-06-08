package server

import (
	"fmt"
	"net"

	"garnet/internal/config"
	"garnet/internal/logger"
)

// SyncServer manages the TCP listener and
// incoming client connections.
type SyncServer struct {
	cfg *config.Config
}

// NewSyncServer creates a new server instance.
func NewSyncServer(cfg *config.Config) *SyncServer {
	return &SyncServer{
		cfg: cfg,
	}
}

// Start begins listening for TCP connections
// and blocks until the server terminates.
func (s *SyncServer) Start() error {

	addr := fmt.Sprintf(
		"%s:%d",
		s.cfg.Host,
		s.cfg.Port,
	)

	logger.Logger.Print(Banner)
	logger.Logger.Printf(
		"listening on %s",
		addr,
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

		HandleConnection(conn)
	}
}
