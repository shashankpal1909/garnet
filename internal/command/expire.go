package command

import (
	"strconv"
	"time"

	"garnet/internal/resp"
	"garnet/internal/store"
)

// Expire handles the Redis EXPIRE command.
// Sets a timeout on key in seconds.
// Returns an integer reply: 1 if the timeout was set, 0 if the key does not exist.
func Expire(args []resp.Value) []byte {
	if len(args) != 2 {
		return resp.EncodeError("ERR wrong number of arguments for 'expire' command")
	}

	key := string(args[0].Data.([]byte))
	secondsStr := string(args[1].Data.([]byte))

	seconds, err := strconv.ParseInt(secondsStr, 10, 64)
	if err != nil {
		return resp.EncodeError("ERR value is not an integer or out of range")
	}

	item := store.Get(key)
	if item == nil {
		// Key does not exist
		return resp.EncodeInteger(0)
	}

	// Update expiration
	// If seconds is <= 0, Redis deletes the key immediately.
	if seconds <= 0 {
		store.Delete(key)
	} else {
		store.UpdateExpire(key, time.Now().UnixMilli()+(seconds*1000))
	}

	return resp.EncodeInteger(1)
}
