package command

import (
	"bytes"
	"testing"

	"garnet/internal/resp"
	"garnet/internal/store"
)

func TestTTLCommand(t *testing.T) {
	t.Run("TTL on non-existent key", func(t *testing.T) {
		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("missing_key")},
		}

		result := TTL(args)
		expected := resp.EncodeInteger(-2)
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("TTL on key without expiration", func(t *testing.T) {
		store.Put("forever_key", store.NewItem("val", -1))

		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("forever_key")},
		}

		result := TTL(args)
		expected := resp.EncodeInteger(-1)
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("TTL on key with expiration", func(t *testing.T) {
		// Set expiration to 50 seconds from now (50000 ms)
		store.Put("expire_key", store.NewItem("val", 50000))

		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("expire_key")},
		}

		result := TTL(args)

		// It could be 49 or 50 depending on timing, but generally 50 if fast enough
		expected50 := resp.EncodeInteger(50)
		expected49 := resp.EncodeInteger(49)

		if !bytes.Equal(result, expected50) && !bytes.Equal(result, expected49) {
			t.Errorf("Expected %q or %q, got %q", expected50, expected49, result)
		}
	})

	t.Run("TTL wrong number of arguments", func(t *testing.T) {
		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("key1")},
			{Type: resp.BulkString, Data: []byte("key2")},
		}

		result := TTL(args)
		expected := resp.EncodeError("ERR wrong number of arguments for 'ttl' command")
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected error response, got %q", result)
		}
	})
}
