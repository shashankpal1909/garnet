package server_test

import (
	"garnet/internal/config"
	"garnet/internal/server"
	"net"
	"testing"
	"time"
)

func TestServer_New(t *testing.T) {
	cfg := &config.Config{
		Host: "127.0.0.1",
		Port: 6379,
	}

	srv := server.NewSyncServer(cfg)
	if srv == nil {
		t.Fatal("Expected New to return a Server instance, got nil")
	}
}

func TestServer_StartError(t *testing.T) {
	// Bind to port 0 to get a random free port
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to bind random port: %v", err)
	}
	defer l.Close()

	port := l.Addr().(*net.TCPAddr).Port

	cfg := &config.Config{
		Host: "127.0.0.1",
		Port: port,
	}

	srv := server.NewSyncServer(cfg)

	// Start should immediately return an error because the port is already bound by our listener
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("Expected Start to return an error when port is in use, got nil")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Start blocked instead of returning error")
	}
}
