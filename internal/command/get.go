package command

import (
	"garnet/internal/resp"
	"garnet/internal/store"
)

// Get handles the Redis GET command.
// It retrieves the value of a key from the global store.
// If the key does not exist or has expired, it returns RESP_NIL.
func Get(args []resp.Value) []byte {
	// A GET command requires exactly one argument: the key.
	if len(args) != 1 {
		return resp.EncodeError("ERR wrong number of arguments for 'get' command")
	}

	key := string(args[0].Data.([]byte))
	
	// Retrieve the item from the store.
	// Note: store.Get() internally checks for expiration and will return nil if expired.
	item := store.Get(key)

	// If the item doesn't exist or expired, return a RESP Null Bulk String ($-1\r\n)
	if item == nil {
		return resp.RESP_NIL
	}

	// Return the value as a RESP Bulk String
	return resp.EncodeBulkString(item.Value.(string))
}
