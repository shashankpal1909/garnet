package command

import (
	"bytes"
	"testing"

	"garnet/internal/resp"
	"garnet/internal/store"
)

func TestSetCommand(t *testing.T) {
	// store.Delete("test_key") // Just to make sure it's clean, but let's actually just run commands

	t.Run("Basic SET", func(t *testing.T) {
		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("mykey")},
			{Type: resp.BulkString, Data: []byte("myvalue")},
		}

		result := Set(args)
		if !bytes.Equal(result, resp.RESP_OK) {
			t.Errorf("Expected +OK\\r\\n, got %q", result)
		}

		item := store.Get("mykey")
		if item == nil || item.Value.(string) != "myvalue" {
			t.Errorf("Item was not stored correctly")
		}
		if item.ExpiresAt != -1 {
			t.Errorf("Expected no expiration, got %v", item.ExpiresAt)
		}
	})

	t.Run("SET with EX", func(t *testing.T) {
		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("ttlkey")},
			{Type: resp.BulkString, Data: []byte("ttlvalue")},
			{Type: resp.BulkString, Data: []byte("EX")},
			{Type: resp.BulkString, Data: []byte("10")},
		}

		result := Set(args)
		if !bytes.Equal(result, resp.RESP_OK) {
			t.Errorf("Expected +OK\\r\\n, got %q", result)
		}

		item := store.Get("ttlkey")
		if item == nil || item.Value.(string) != "ttlvalue" {
			t.Errorf("Item was not stored correctly")
		}
		if item.ExpiresAt == -1 {
			t.Errorf("Expected expiration to be set, got -1")
		}
	})

	t.Run("SET wrong number of args", func(t *testing.T) {
		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("onlykey")},
		}

		result := Set(args)
		expected := resp.EncodeError("ERR wrong number of arguments for 'set' command")
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected error response, got %q", result)
		}
	})

	t.Run("SET EX syntax error missing value", func(t *testing.T) {
		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("k")},
			{Type: resp.BulkString, Data: []byte("v")},
			{Type: resp.BulkString, Data: []byte("EX")},
		}

		result := Set(args)
		expected := resp.EncodeError("ERR syntax error")
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected syntax error, got %q", result)
		}
	})

	t.Run("SET EX invalid TTL", func(t *testing.T) {
		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("k")},
			{Type: resp.BulkString, Data: []byte("v")},
			{Type: resp.BulkString, Data: []byte("EX")},
			{Type: resp.BulkString, Data: []byte("notanumber")},
		}

		result := Set(args)
		expected := resp.EncodeError("ERR invalid TTL value")
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected invalid TTL error, got %q", result)
		}
	})
}
