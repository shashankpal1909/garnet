//go:build linux

package server_test

import (
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"

	"garnet/internal/config"
	"garnet/internal/server"
)

func TestAsyncServer_New(t *testing.T) {
	cfg := &config.Config{
		Host:       "127.0.0.1",
		Port:       6379,
		MaxClients: 10,
	}

	srv := server.NewAsyncServer(cfg)
	if srv == nil {
		t.Fatal("Expected NewAsyncServer to return an AsyncServer instance, got nil")
	}
}

func TestAsyncServer_StartError(t *testing.T) {
	// Bind to port 0 to get a random free port
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to bind random port: %v", err)
	}
	defer l.Close()

	port := l.Addr().(*net.TCPAddr).Port

	cfg := &config.Config{
		Host:       "127.0.0.1",
		Port:       port,
		MaxClients: 10,
	}

	srv := server.NewAsyncServer(cfg)

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

func TestAsyncServer_EndToEnd(t *testing.T) {
	// Get a random free port for the async server
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to bind random port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close() // Free the port so the server can bind to it

	cfg := &config.Config{
		Host:       "127.0.0.1",
		Port:       port,
		MaxClients: 10,
	}

	srv := server.NewAsyncServer(cfg)
	go func() {
		_ = srv.Start()
	}()

	// Wait a moment for server to start
	time.Sleep(100 * time.Millisecond)

	// Connect a client
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("Failed to connect to async server: %v", err)
	}

	// Send a valid PING command
	_, err = conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
	if err != nil {
		t.Fatalf("Failed to write to async server: %v", err)
	}

	// Wait for the response
	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read from async server: %v", err)
	}

	got := buf[:n]
	want := []byte("+PONG\r\n")

	if !bytes.Equal(got, want) {
		t.Errorf("Expected response %q, got %q", want, got)
	}

	// Test fragmentation: write in two chunks
	_, err = conn.Write([]byte("*2\r\n$4\r\nE"))
	if err != nil {
		t.Fatalf("Failed to write partial to async server: %v", err)
	}
	time.Sleep(50 * time.Millisecond)
	_, err = conn.Write([]byte("CHO\r\n$5\r\nHELLO\r\n"))
	if err != nil {
		t.Fatalf("Failed to write partial to async server: %v", err)
	}

	// Read response
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err = conn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read from async server after partial writes: %v", err)
	}

	got = buf[:n]
	want = []byte("$5\r\nHELLO\r\n")
	if !bytes.Equal(got, want) {
		t.Errorf("Expected fragmented response %q, got %q", want, got)
	}

	// Test empty read/disconnection
	conn.Close()

	// Wait a moment for server to process disconnection
	time.Sleep(50 * time.Millisecond)
}
