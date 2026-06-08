package server

import (
	"io"
	"net"

	"garnet/internal/command"
	"garnet/internal/logger"
	"garnet/internal/resp"
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

	decoder := resp.NewDecoder(conn)

	for {
		value, err := decoder.Decode()
		if err != nil {
			if err != io.EOF {
				logger.Logger.Printf(
					"connection read error: %v",
					err,
				)
			}
			return
		}

		logger.Logger.Printf(
			"received command from %s: %+v",
			conn.RemoteAddr(),
			value,
		)

		response := command.Dispatch(value)

		if _, err := conn.Write(response); err != nil {

			logger.Logger.Printf(
				"write failed: %v",
				err,
			)

			return
		}
	}
}
