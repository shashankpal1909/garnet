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
		{
			name:    "Invalid Integer Format",
			input:   []byte(":not-a-number\r\n"),
			wantErr: true,
		},
		{
			name:    "Invalid Bulk String Length",
			input:   []byte("$not-a-number\r\nhello\r\n"),
			wantErr: true,
		},
		{
			name:    "Invalid Array Length",
			input:   []byte("*not-a-number\r\n"),
			wantErr: true,
		},
		{
			name:    "Incomplete Data",
			input:   []byte("*1\r\n$4\r\nPING"),
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

func TestDecodeFromBytes(t *testing.T) {
	input := []byte("*2\r\n$4\r\nECHO\r\n$11\r\nhello world\r\n")
	val, consumed, err := resp.DecodeFromBytes(input)
	if err != nil {
		t.Fatalf("DecodeFromBytes error: %v", err)
	}
	if consumed != len(input) {
		t.Errorf("Expected consumed %d, got %d", len(input), consumed)
	}
	if val.Type != resp.Array || len(val.Data.([]resp.Value)) != 2 {
		t.Errorf("DecodeFromBytes returned incorrect value")
	}

	// Test incomplete
	_, _, err = resp.DecodeFromBytes([]byte("*1\r\n$4\r\nPI"))
	if err != resp.ErrIncomplete {
		t.Errorf("Expected ErrIncomplete, got %v", err)
	}

	// Test multiple commands
	multiInput := []byte("+OK\r\n+PONG\r\n")
	val1, c1, err := resp.DecodeFromBytes(multiInput)
	if err != nil || val1.Data != "OK" {
		t.Errorf("DecodeFromBytes failed on first command")
	}
	val2, _, err := resp.DecodeFromBytes(multiInput[c1:])
	if err != nil || val2.Data != "PONG" {
		t.Errorf("DecodeFromBytes failed on second command")
	}
}
