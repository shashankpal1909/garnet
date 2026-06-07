package command

import (
	"strings"

	"garnet/internal/resp"
)

type CommandFunc func(args []resp.Value) []byte

var registry = map[string]CommandFunc{
	"PING":  Ping,
	"HELLO": Hello,
}

// Dispatch routes incoming commands to their corresponding handlers.
func Dispatch(value resp.Value) []byte {
	if value.Type != resp.Array {
		return resp.EncodeError("ERR expected array")
	}

	array, ok := value.Data.([]resp.Value)
	if !ok || len(array) == 0 {
		return resp.EncodeError("ERR empty command")
	}

	cmdValue := array[0]
	if cmdValue.Type != resp.BulkString {
		return resp.EncodeError("ERR expected bulk string for command")
	}

	cmdBytes, ok := cmdValue.Data.([]byte)
	if !ok {
		return resp.EncodeError("ERR invalid command format")
	}

	cmdName := strings.ToUpper(string(cmdBytes))

	handler, exists := registry[cmdName]
	if !exists {
		return resp.EncodeError("ERR unknown command")
	}

	// Pass arguments to the handler (everything after the command name)
	args := array[1:]
	return handler(args)
}
