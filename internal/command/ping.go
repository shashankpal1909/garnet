package command

import (
	"garnet/internal/resp"
)

// Ping handles the PING command
func Ping(args []resp.Value) []byte {
	if len(args) == 0 {
		return resp.EncodeSimpleString("PONG")
	} else if len(args) == 1 {
		if arg, ok := args[0].Data.([]byte); ok {
			return resp.EncodeBulkString(string(arg))
		}
	}
	// Fallback if args format is unexpected
	return resp.EncodeSimpleString("PONG")
}
