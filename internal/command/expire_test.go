package command

import (
	"bytes"
	"testing"
	"time"

	"garnet/internal/resp"
	"garnet/internal/store"
)

func TestExpireCommand(t *testing.T) {
	t.Run("EXPIRE existing key", func(t *testing.T) {
		store.Put("expire_test", store.NewItem("val", -1))

		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("expire_test")},
			{Type: resp.BulkString, Data: []byte("10")},
		}

		result := Expire(args)
		expected := resp.EncodeInteger(1)
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected %q, got %q", expected, result)
		}

		item := store.Get("expire_test")
		if item == nil {
			t.Fatalf("Item was accidentally deleted")
		}
		if item.ExpiresAt == -1 {
			t.Errorf("Expiration was not updated")
		}

		remainingSec := (item.ExpiresAt - time.Now().UnixMilli()) / 1000
		if remainingSec < 9 || remainingSec > 10 {
			t.Errorf("Expected ~10s expiration, got %ds", remainingSec)
		}
	})

	t.Run("EXPIRE non-existent key", func(t *testing.T) {
		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("missing_key")},
			{Type: resp.BulkString, Data: []byte("10")},
		}

		result := Expire(args)
		expected := resp.EncodeInteger(0)
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected %q, got %q", expected, result)
		}
	})

	t.Run("EXPIRE with 0 seconds", func(t *testing.T) {
		store.Put("delete_me", store.NewItem("val", -1))

		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("delete_me")},
			{Type: resp.BulkString, Data: []byte("0")},
		}

		result := Expire(args)
		expected := resp.EncodeInteger(1)
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected %q, got %q", expected, result)
		}

		if store.Get("delete_me") != nil {
			t.Errorf("Key should have been deleted immediately")
		}
	})

	t.Run("EXPIRE invalid number", func(t *testing.T) {
		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("key")},
			{Type: resp.BulkString, Data: []byte("not_a_number")},
		}

		result := Expire(args)
		expected := resp.EncodeError("ERR value is not an integer or out of range")
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected error, got %q", result)
		}
	})
}
