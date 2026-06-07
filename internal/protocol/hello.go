package protocol

import "strings"

// Handle routes incoming commands
// to their corresponding handlers.
func Handle(input string) string {
	switch strings.ToUpper(input) {

	case "HELLO":
		return "Hello from Garnet"

	default:
		return "ERR unknown command"
	}
}
