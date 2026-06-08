package command

import (
	"bytes"
	"testing"

	"garnet/internal/resp"
	"garnet/internal/store"
)

func TestDelCommand(t *testing.T) {
	t.Run("DEL existing keys", func(t *testing.T) {
		store.Put("key1", store.NewItem("val1", -1))
		store.Put("key2", store.NewItem("val2", -1))

		args := []resp.Value{
			{Type: resp.BulkString, Data: []byte("key1")},
			{Type: resp.BulkString, Data: []byte("key2")},
			{Type: resp.BulkString, Data: []byte("missing_key")},
		}

		result := Del(args)
		expected := resp.EncodeInteger(2)
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected %q, got %q", expected, result)
		}

		if store.Get("key1") != nil || store.Get("key2") != nil {
			t.Errorf("Keys were not actually deleted")
		}
	})

	t.Run("DEL missing arguments", func(t *testing.T) {
		args := []resp.Value{}
		result := Del(args)
		expected := resp.EncodeError("ERR wrong number of arguments for 'del' command")
		if !bytes.Equal(result, expected) {
			t.Errorf("Expected error, got %q", result)
		}
	})
}
