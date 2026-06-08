package resp_test

import (
	"bytes"
	"garnet/internal/resp"
	"testing"
)

func TestEncodeSimpleString(t *testing.T) {
	got := resp.EncodeSimpleString("PONG")
	want := []byte("+PONG\r\n")

	if !bytes.Equal(got, want) {
		t.Errorf("EncodeSimpleString() = %q, want %q", got, want)
	}
}

func TestEncodeError(t *testing.T) {
	got := resp.EncodeError("ERR unknown command")
	want := []byte("-ERR unknown command\r\n")

	if !bytes.Equal(got, want) {
		t.Errorf("EncodeError() = %q, want %q", got, want)
	}
}

func TestEncodeBulkString(t *testing.T) {
	got := resp.EncodeBulkString("PONG")
	want := []byte("$4\r\nPONG\r\n")

	if !bytes.Equal(got, want) {
		t.Errorf("EncodeBulkString() = %q, want %q", got, want)
	}
}

func TestEncodeInteger(t *testing.T) {
	got := resp.EncodeInteger(42)
	want := []byte(":42\r\n")

	if !bytes.Equal(got, want) {
		t.Errorf("EncodeInteger() = %q, want %q", got, want)
	}
}
