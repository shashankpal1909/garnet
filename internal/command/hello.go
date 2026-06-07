package command

import (
	"garnet/internal/resp"
)

// Hello handles the HELLO command
func Hello(args []resp.Value) []byte {
	// The HELLO command can take a protocol version, but for now we just return a simple greeting
	return resp.EncodeSimpleString("Hello from Garnet")
}
