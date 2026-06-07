package command_test

import (
	"bytes"
	"garnet/internal/command"
	"garnet/internal/resp"
	"testing"
)

func TestPing(t *testing.T) {
	tests := []struct {
		name string
		args []resp.Value
		want []byte
	}{
		{
			name: "No Arguments",
			args: []resp.Value{},
			want: []byte("+PONG\r\n"),
		},
		{
			name: "One Argument",
			args: []resp.Value{
				{Type: resp.BulkString, Data: []byte("hello")},
			},
			want: []byte("$5\r\nhello\r\n"),
		},
		{
			name: "Multiple Arguments (fallback)",
			args: []resp.Value{
				{Type: resp.BulkString, Data: []byte("hello")},
				{Type: resp.BulkString, Data: []byte("world")},
			},
			want: []byte("+PONG\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := command.Ping(tt.args)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("Ping() = %q, want %q", got, tt.want)
			}
		})
	}
}
