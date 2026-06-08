package command

import (
	"bytes"
	"testing"

	"garnet/internal/resp"
	"garnet/internal/store"
)

func TestGetCommand(t *testing.T) {
	// Setup test data
	store.Put("existing_key", store.NewItem("existing_value", -1))
	store.Put("expired_key", store.NewItem("expired_value", 1)) // Expires almost immediately

	t.Run("GET existing key", func(t *testing.T) {
		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("existing_key")},
		}

		result := Get(args)
		expected := resp.EncodeBulkString("existing_value")
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("GET non-existent key", func(t *testing.T) {
		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("unknown_key")},
		}

		result := Get(args)
		if !bytes.Equal(result, resp.RESP_NIL) {
			t.Errorf("Expected nil response, got %q", result)
		}
	})

	t.Run("GET wrong number of args", func(t *testing.T) {
		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("key1")},
			{Type: resp.BulkString, Data: []byte("key2")},
		}

		result := Get(args)
		expected := resp.EncodeError("ERR wrong number of arguments for 'get' command")
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected error response, got %q", result)
		}
	})
}
