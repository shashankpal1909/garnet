package command_test

import (
	"bytes"
	"garnet/internal/command"
	"garnet/internal/resp"
	"testing"
)

func TestEcho(t *testing.T) {
	tests := []struct {
		name string
		args []resp.Value
		want []byte
	}{
		{
			name: "Valid Echo",
			args: []resp.Value{
				{Type: resp.BulkString, Data: []byte("hello world")},
			},
			want: []byte("$11\r\nhello world\r\n"),
		},
		{
			name: "No Arguments",
			args: []resp.Value{},
			want: []byte("-ERR wrong number of arguments for 'echo' command\r\n"),
		},
		{
			name: "Too Many Arguments",
			args: []resp.Value{
				{Type: resp.BulkString, Data: []byte("hello")},
				{Type: resp.BulkString, Data: []byte("world")},
			},
			want: []byte("-ERR wrong number of arguments for 'echo' command\r\n"),
		},
		{
			name: "Wrong Argument Type",
			args: []resp.Value{
				{Type: resp.Integer, Data: int64(100)},
			},
			want: []byte("-ERR wrong type of argument for 'echo' command\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := command.Echo(tt.args)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("Echo() = %q, want %q", got, tt.want)
			}
		})
	}
}
