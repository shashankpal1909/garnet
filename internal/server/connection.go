package server

import (
	"bufio"
	"net"
	"strings"

	"garnet/internal/logger"
	"garnet/internal/protocol"
)

// HandleConnection manages the lifecycle
// of a single client connection.
func HandleConnection(conn net.Conn) {
	defer conn.Close()

	logger.Logger.Printf(
		"client connected: %s",
		conn.RemoteAddr(),
	)

	defer logger.Logger.Printf(
		"client disconnected: %s",
		conn.RemoteAddr(),
	)

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		input := strings.TrimSpace(
			scanner.Text(),
		)

		logger.Logger.Printf(
			"received command from %s: %s",
			conn.RemoteAddr(),
			input,
		)

		response := protocol.Handle(input)

		if _, err := conn.Write(
			[]byte(response + "\n"),
		); err != nil {

			logger.Logger.Printf(
				"write failed: %v",
				err,
			)

			return
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Logger.Printf(
			"connection read error: %v",
			err,
		)
	}
}
