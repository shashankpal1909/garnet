package store

import (
	"testing"
	"time"
)

func TestStore_PutAndGet(t *testing.T) {
	// Clean store just in case
	store = *New(0, &NoOpEvictor{})

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
	store = *New(0, &NoOpEvictor{})
	retrieved := Get("non_existent")
	if retrieved != nil {
		t.Fatalf("Expected nil for non-existent key, got %v", retrieved)
	}
}

func TestStore_Expiration(t *testing.T) {
	store = *New(0, &NoOpEvictor{})

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

func TestStore_ActiveExpire(t *testing.T) {
	store = *New(0, &NoOpEvictor{})

	// Put 10 persistent keys
	for i := 0; i < 10; i++ {
		Put("persistent_"+string(rune(i)), NewItem("val", -1))
	}

	// Put 30 expiring keys with 10ms TTL
	for i := 0; i < 30; i++ {
		Put("expire_"+string(rune(i)), NewItem("val", 10))
	}

	// Before sleep, ActiveExpire should delete 0
	deleted := ActiveExpire()
	if deleted != 0 {
		t.Errorf("Expected 0 deleted, got %d", deleted)
	}

	// Wait for TTLs to pass
	time.Sleep(20 * time.Millisecond)

	// ActiveExpire should delete all 30 since the sampling algorithm loops
	// when > 25% of sampled keys are expired
	deleted = ActiveExpire()
	if deleted != 30 {
		t.Errorf("Expected 30 deleted, got %d", deleted)
	}

	// Verify maps are clean
	if len(store.expires) != 0 {
		t.Errorf("Expected expires map to be empty, got %d", len(store.expires))
	}
	if len(store.data) != 10 {
		t.Errorf("Expected exactly 10 persistent keys remaining, got %d", len(store.data))
	}
}

func TestStore_Delete(t *testing.T) {
	store = *New(0, &NoOpEvictor{})
	Put("to_delete", NewItem("val", -1))

	if !Delete("to_delete") {
		t.Errorf("Expected true when deleting existing key")
	}
	if Delete("non_existent") {
		t.Errorf("Expected false when deleting non-existent key")
	}
}

func TestStore_UpdateExpire(t *testing.T) {
	store = *New(0, &NoOpEvictor{})
	Put("update_me", NewItem("val", -1))

	// Update to expire
	UpdateExpire("update_me", time.Now().UnixMilli()+10)
	if _, exists := store.expires["update_me"]; !exists {
		t.Errorf("Expected key to be added to expires map")
	}

	// Update to persistent
	UpdateExpire("update_me", -1)
	if _, exists := store.expires["update_me"]; exists {
		t.Errorf("Expected key to be removed from expires map")
	}

	// Update non-existent
	UpdateExpire("non_existent", -1)
}
