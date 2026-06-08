package command

import (
	"time"

	"garnet/internal/resp"
	"garnet/internal/store"
)

// TTL handles the Redis TTL command.
// It returns the remaining time to live of a key that has a timeout in seconds.
func TTL(args []resp.Value) []byte {
	// TTL requires exactly one argument: the key.
	if len(args) != 1 {
		return resp.EncodeError("ERR wrong number of arguments for 'ttl' command")
	}

	key := string(args[0].Data.([]byte))

	// Retrieve the item. store.Get handles checking and removing expired items.
	item := store.Get(key)

	// Return -2 if the key does not exist (or has already expired).
	if item == nil {
		return resp.EncodeInteger(-2)
	}

	// Return -1 if the key exists but has no associated expire.
	if item.ExpiresAt == -1 {
		return resp.EncodeInteger(-1)
	}

	// Calculate remaining time to live in seconds
	remainingMs := item.ExpiresAt - time.Now().UnixMilli()

	// Convert to seconds, rounding down
	ttlSec := remainingMs / 1000

	// If it's technically expired right this millisecond but hasn't been cleaned up,
	// or it's < 1 second, it should still say 0 or expire
	if ttlSec < 0 {
		return resp.EncodeInteger(-2) // Should never really happen because of store.Get()
	}

	return resp.EncodeInteger(ttlSec)
}
