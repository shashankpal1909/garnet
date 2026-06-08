package command

import (
	"garnet/internal/resp"
	"garnet/internal/store"
)

// Del handles the Redis DEL command.
// Removes the specified keys. A key is ignored if it does not exist.
// Returns an integer reply: The number of keys that were removed.
func Del(args []resp.Value) []byte {
	if len(args) < 1 {
		return resp.EncodeError("ERR wrong number of arguments for 'del' command")
	}

	deletedCount := int64(0)

	for _, arg := range args {
		key := string(arg.Data.([]byte))

		// Wait, if it has passively expired but is still in the map, does Delete return true?
		// To match Redis behavior strictly, expired items shouldn't count.
		// store.Get(key) returns nil if it doesn't exist or is expired, and handles the deletion!
		item := store.Get(key)
		if item != nil {
			// Actually delete it now
			store.Delete(key)
			deletedCount++
		}
	}

	return resp.EncodeInteger(deletedCount)
}
