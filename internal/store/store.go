package store

import "time"

// Item represents a stored value and its expiration metadata.
type Item struct {
	Value     interface{}
	// ExpiresAt holds the Unix epoch timestamp in milliseconds when the item should expire.
	// A value of -1 indicates that the item does not expire.
	ExpiresAt int64
}

// Store holds the global in-memory dataset.
// Note: Because the Garnet async event loop is single-threaded,
// we do not need a sync.Mutex or sync.RWMutex here!
type Store struct {
	data map[string]*Item
}

// store is the global singleton instance.
var store Store

func init() {
	store = *New()
}

// New creates and initializes a new Store.
func New() *Store {
	return &Store{
		data: make(map[string]*Item),
	}
}

// NewItem creates a new Item. If durationMs is > 0, the expiration time
// is computed dynamically relative to the current time.
func NewItem(value interface{}, durationMs int64) *Item {
	expiresAt := int64(-1)
	if durationMs > 0 {
		expiresAt = time.Now().UnixMilli() + durationMs
	}
	return &Item{
		Value:     value,
		ExpiresAt: expiresAt,
	}
}

// Put inserts or updates an item in the store.
func Put(k string, item *Item) {
	store.data[k] = item
}

// Get retrieves an item by key. 
// If the item has expired, it removes the item and returns nil.
func Get(k string) *Item {
	item, exists := store.data[k]
	if !exists {
		return nil
	}

	// Check if the item has expired
	if item.ExpiresAt != -1 && time.Now().UnixMilli() >= item.ExpiresAt {
		delete(store.data, k)
		return nil
	}

	return item
}
