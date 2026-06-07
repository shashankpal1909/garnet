package server_test

import (
	"bytes"
	"garnet/internal/server"
	"net"
	"testing"
	"time"
)

func TestHandleConnection(t *testing.T) {
	clientConn, serverConn := net.Pipe()

	// Run HandleConnection in the background
	go server.HandleConnection(serverConn)

	// Send a valid PING command from the client
	_, err := clientConn.Write([]byte("*1\r\n$4\r\nPING\r\n"))
	if err != nil {
		t.Fatalf("Failed to write to client connection: %v", err)
	}

	// Wait for the response
	buf := make([]byte, 1024)
	clientConn.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err := clientConn.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read from client connection: %v", err)
	}

	got := buf[:n]
	want := []byte("+PONG\r\n")

	if !bytes.Equal(got, want) {
		t.Errorf("Expected response %q, got %q", want, got)
	}

	// Close client connection which should terminate HandleConnection
	clientConn.Close()
}
