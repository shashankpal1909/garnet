package command_test

import (
	"bytes"
	"garnet/internal/command"
	"garnet/internal/resp"
	"testing"
)

func TestHello(t *testing.T) {
	got := command.Hello([]resp.Value{})
	want := []byte("+Hello from Garnet\r\n")

	if !bytes.Equal(got, want) {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
