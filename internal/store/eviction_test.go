package store

import (
	"testing"

	"garnet/internal/config"
)

func TestEviction_NoEviction(t *testing.T) {
	Init(&config.Config{
		MaxKeys:        2,
		EvictionPolicy: config.NoEviction,
	})

	err := Put("k1", NewItem("v1", -1))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = Put("k2", NewItem("v2", -1))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Should fail because limit is 2 and policy is noeviction
	err = Put("k3", NewItem("v3", -1))
	if err != ErrMaxKeysExceeded {
		t.Errorf("expected ErrMaxKeysExceeded, got %v", err)
	}
}

func TestEviction_AllKeysRandom(t *testing.T) {
	Init(&config.Config{
		MaxKeys:        2,
		EvictionPolicy: config.AllKeysRandom,
	})

	Put("k1", NewItem("v1", -1))
	Put("k2", NewItem("v2", -1))

	// Should succeed by evicting a random key
	err := Put("k3", NewItem("v3", -1))
	if err != nil {
		t.Errorf("expected successful eviction and insert, got %v", err)
	}

	if len(store.data) != 2 {
		t.Errorf("expected exactly 2 keys, got %d", len(store.data))
	}
}

func TestEviction_VolatileRandom(t *testing.T) {
	Init(&config.Config{
		MaxKeys:        2,
		EvictionPolicy: config.VolatileRandom,
	})

	// Put 2 persistent keys
	Put("p1", NewItem("v1", -1))
	Put("p2", NewItem("v2", -1))

	// Trying to put a third should fail because volatile-random can only evict keys with a TTL.
	// Since p1 and p2 have no TTL, they cannot be evicted.
	err := Put("k3", NewItem("v3", -1))
	if err != ErrMaxKeysExceeded {
		t.Errorf("expected ErrMaxKeysExceeded because no keys have TTL, got %v", err)
	}

	// Clear the store and try again with TTL keys
	store.data = make(map[string]*Item)
	store.expires = make(map[string]struct{})

	// Add one volatile and one persistent key
	Put("v1", NewItem("v1", 1000))
	Put("p1", NewItem("p1", -1))

	// Putting a third key should now successfully evict the volatile key "v1"
	err = Put("p2", NewItem("p2", -1))
	if err != nil {
		t.Errorf("expected success, got %v", err)
	}

	if Get("v1") != nil {
		t.Errorf("expected v1 to be evicted, but it was found")
	}
}

func TestEviction_UnknownPolicy(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic on unknown policy")
		}
	}()
	NewEvictor("invalid-policy")
}

func TestEviction_EmptyStore(t *testing.T) {
	// Test that evictors return 0 when store is empty
	s := New(0, &NoOpEvictor{})

	r := &AllKeysRandomEvictor{}
	if deleted := r.Evict(s); deleted != 0 {
		t.Errorf("expected 0, got %d", deleted)
	}

	v := &VolatileRandomEvictor{}
	if deleted := v.Evict(s); deleted != 0 {
		t.Errorf("expected 0, got %d", deleted)
	}
}
