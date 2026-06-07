package command_test

import (
	"bytes"
	"garnet/internal/command"
	"garnet/internal/resp"
	"testing"
)

func TestDispatch(t *testing.T) {
	tests := []struct {
		name  string
		input resp.Value
		want  []byte
	}{
		{
			name:  "Invalid Type",
			input: resp.Value{Type: resp.SimpleString, Data: "PING"},
			want:  []byte("-ERR expected array\r\n"),
		},
		{
			name:  "Empty Array",
			input: resp.Value{Type: resp.Array, Data: []resp.Value{}},
			want:  []byte("-ERR empty command\r\n"),
		},
		{
			name: "Invalid Command Type",
			input: resp.Value{
				Type: resp.Array,
				Data: []resp.Value{
					{Type: resp.Integer, Data: int64(1)},
				},
			},
			want: []byte("-ERR expected bulk string for command\r\n"),
		},
		{
			name: "Unknown Command",
			input: resp.Value{
				Type: resp.Array,
				Data: []resp.Value{
					{Type: resp.BulkString, Data: []byte("FOO")},
				},
			},
			want: []byte("-ERR unknown command\r\n"),
		},
		{
			name: "Valid PING Command",
			input: resp.Value{
				Type: resp.Array,
				Data: []resp.Value{
					{Type: resp.BulkString, Data: []byte("PING")},
				},
			},
			want: []byte("+PONG\r\n"),
		},
		{
			name: "Valid HELLO Command",
			input: resp.Value{
				Type: resp.Array,
				Data: []resp.Value{
					{Type: resp.BulkString, Data: []byte("hello")},
				},
			},
			want: []byte("+Hello from Garnet\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := command.Dispatch(tt.input)
			if !bytes.Equal(got, tt.want) {
				t.Errorf("Dispatch() = %q, want %q", got, tt.want)
			}
		})
	}
}
