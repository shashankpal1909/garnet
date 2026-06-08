package store

import (
	"testing"
	"time"
)

func TestStore_PutAndGet(t *testing.T) {
	// Clean store just in case
	store = *New()

	item := NewItem("test_value", -1)
	Put("test_key", item)

	retrieved := Get("test_key")
	if retrieved == nil {
		t.Fatalf("Expected to retrieve item, got nil")
	}

	if retrieved.Value.(string) != "test_value" {
		t.Errorf("Expected 'test_value', got '%v'", retrieved.Value)
	}
}

func TestStore_GetNonExistent(t *testing.T) {
	store = *New()
	retrieved := Get("non_existent")
	if retrieved != nil {
		t.Fatalf("Expected nil for non-existent key, got %v", retrieved)
	}
}

func TestStore_Expiration(t *testing.T) {
	store = *New()

	// 50ms TTL
	item := NewItem("expire_me", 50)
	Put("expire_key", item)

	// Should be available immediately
	retrieved := Get("expire_key")
	if retrieved == nil {
		t.Fatalf("Expected to retrieve item immediately, got nil")
	}

	// Wait for TTL to pass
	time.Sleep(60 * time.Millisecond)

	// Should be expired now
	expired := Get("expire_key")
	if expired != nil {
		t.Fatalf("Expected nil for expired key, got %v", expired)
	}

	// Make sure it was actually deleted from the map
	if _, exists := store.data["expire_key"]; exists {
		t.Errorf("Expected key to be deleted from internal map after expiration")
	}
}
