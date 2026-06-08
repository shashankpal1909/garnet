package command

import (
	"garnet/internal/resp"
)

// Echo handles the ECHO command.
// It expects exactly one argument and returns that argument as a bulk string.
func Echo(args []resp.Value) []byte {
	// Ensure exactly one argument is provided
	if len(args) != 1 {
		return resp.EncodeError("ERR wrong number of arguments for 'echo' command")
	}

	// Extract the argument data as a byte slice
	value, ok := args[0].Data.([]byte)
	if !ok {
		return resp.EncodeError("ERR wrong type of argument for 'echo' command")
	}

	// Return the argument formatted as a RESP bulk string
	return resp.EncodeBulkString(string(value))
}
