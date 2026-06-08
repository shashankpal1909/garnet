package store

import "time"

// Item represents a stored value and its expiration metadata.
type Item struct {
	Value interface{}
	// ExpiresAt holds the Unix epoch timestamp in milliseconds when the item should expire.
	// A value of -1 indicates that the item does not expire.
	ExpiresAt int64
}

// Store holds the global in-memory dataset.
// Note: Because the Garnet async event loop is single-threaded,
// we do not need a sync.Mutex or sync.RWMutex here!
type Store struct {
	data    map[string]*Item
	expires map[string]struct{}
}

// store is the global singleton instance.
var store Store

func init() {
	store = *New()
}

// New creates and initializes a new Store.
func New() *Store {
	return &Store{
		data:    make(map[string]*Item),
		expires: make(map[string]struct{}),
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
	if item.ExpiresAt != -1 {
		store.expires[k] = struct{}{}
	} else {
		delete(store.expires, k)
	}
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
		delete(store.expires, k)
		return nil
	}

	return item
}

// Delete removes a key from the store. Returns true if the key existed.
func Delete(k string) bool {
	if _, exists := store.data[k]; exists {
		delete(store.data, k)
		delete(store.expires, k)
		return true
	}
	return false
}

// UpdateExpire modifies the expiration of an existing key.
func UpdateExpire(k string, expiresAt int64) {
	if item, exists := store.data[k]; exists {
		item.ExpiresAt = expiresAt
		if expiresAt != -1 {
			store.expires[k] = struct{}{}
		} else {
			delete(store.expires, k)
		}
	}
}

// ActiveExpire samples keys with TTLs and deletes expired ones.
// Returns the number of keys deleted.
func ActiveExpire() int {
	totalDeleted := 0

	for {
		sampled := 0
		expiredCount := 0

		// Iterate over expires map. Go's map iteration is pseudo-random.
		for k := range store.expires {
			item, exists := store.data[k]
			if !exists {
				delete(store.expires, k)
				continue
			}

			if time.Now().UnixMilli() >= item.ExpiresAt {
				delete(store.data, k)
				delete(store.expires, k)
				expiredCount++
				totalDeleted++
			}

			sampled++
			if sampled >= 20 {
				break
			}
		}

		// Stop if less than 25% of sampled keys were expired, or if none were sampled.
		if sampled == 0 || float64(expiredCount)/float64(sampled) < 0.25 {
			break
		}
	}

	return totalDeleted
}
