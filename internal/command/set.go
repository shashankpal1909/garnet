package command

import (
	"strconv"
	"strings"

	"garnet/internal/resp"
	"garnet/internal/store"
)

// Set handles the Redis SET command.
// It supports basic SET key value, as well as the EX (expiration in seconds) option.
func Set(args []resp.Value) []byte {
	// A SET command must have at least a key and a value
	if len(args) <= 1 {
		return resp.EncodeError("ERR wrong number of arguments for 'set' command")
	}

	key := string(args[0].Data.([]byte))
	val := string(args[1].Data.([]byte))

	// By default, keys do not expire (-1 indicates no expiration)
	exDurationMs := int64(-1)

	// Parse optional arguments (EX for now)
	for i := 2; i < len(args); i++ {
		opt := strings.ToUpper(string(args[i].Data.([]byte)))
		switch opt {
		case "EX":
			// EX specifies an expiration in seconds. We need one more argument for the value.
			if i+1 >= len(args) {
				return resp.EncodeError("ERR syntax error")
			}

			exDurationSec, err := strconv.ParseInt(string(args[i+1].Data.([]byte)), 10, 64)
			if err != nil {
				return resp.EncodeError("ERR invalid TTL value")
			}

			// Convert seconds to milliseconds for the store
			exDurationMs = exDurationSec * 1000

			// Skip the value argument since we just processed it
			i++

		default:
			return resp.EncodeError("ERR syntax error")
		}
	}

	// Store the item with the computed TTL
	store.Put(key, store.NewItem(val, exDurationMs))
	return resp.RESP_OK
}
