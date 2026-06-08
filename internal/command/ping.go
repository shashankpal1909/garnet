package command

import (
	"garnet/internal/resp"
)

// Ping handles the PING command.
// It returns "PONG" if no arguments are provided,
// or echoes the argument as a bulk string if exactly one argument is given.
func Ping(args []resp.Value) []byte {
	// Return PONG if no arguments are provided
	if len(args) == 0 {
		return resp.EncodeSimpleString("PONG")
	}

	// Return the argument as a bulk string if exactly one is provided
	if len(args) == 1 {
		if arg, ok := args[0].Data.([]byte); ok {
			return resp.EncodeBulkString(string(arg))
		}
	}

	// Return error for wrong number of arguments or invalid types
	return resp.EncodeError("ERR wrong number of arguments for 'ping' command")
}
