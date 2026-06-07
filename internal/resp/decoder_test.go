package resp_test

import (
	"bytes"
	"garnet/internal/resp"
	"reflect"
	"testing"
)

func TestDecoder(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    resp.Value
		wantErr bool
	}{
		{
			name:  "Simple String",
			input: []byte("+OK\r\n"),
			want:  resp.Value{Type: resp.SimpleString, Data: "OK"},
		},
		{
			name:  "Error",
			input: []byte("-Error message\r\n"),
			want:  resp.Value{Type: resp.Error, Data: "Error message"},
		},
		{
			name:  "Integer",
			input: []byte(":1000\r\n"),
			want:  resp.Value{Type: resp.Integer, Data: int64(1000)},
		},
		{
			name:  "Bulk String",
			input: []byte("$5\r\nhello\r\n"),
			want:  resp.Value{Type: resp.BulkString, Data: []byte("hello")},
		},
		{
			name:  "Null Bulk String",
			input: []byte("$-1\r\n"),
			want:  resp.Value{Type: resp.BulkString, Data: nil},
		},
		{
			name:  "Array",
			input: []byte("*2\r\n$4\r\nECHO\r\n$11\r\nhello world\r\n"),
			want: resp.Value{
				Type: resp.Array,
				Data: []resp.Value{
					{Type: resp.BulkString, Data: []byte("ECHO")},
					{Type: resp.BulkString, Data: []byte("hello world")},
				},
			},
		},
		{
			name:  "Null Array",
			input: []byte("*-1\r\n"),
			want:  resp.Value{Type: resp.Array, Data: nil},
		},
		{
			name:    "Invalid Type",
			input:   []byte("?unknown\r\n"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoder := resp.NewDecoder(bytes.NewReader(tt.input))
			got, err := decoder.Decode()

			if (err != nil) != tt.wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}
