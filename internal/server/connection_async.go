//go:build linux

package server

import (
	"syscall"

	"garnet/internal/command"
	"garnet/internal/logger"
	"garnet/internal/resp"
)

type AsyncConnection struct {
	fd     int
	buffer []byte
}

func NewAsyncConnection(fd int) *AsyncConnection {
	return &AsyncConnection{
		fd:     fd,
		buffer: make([]byte, 0, 4096),
	}
}

func (c *AsyncConnection) Read() error {
	buf := make([]byte, 4096)
	
	// Read as much data as is available from the non-blocking socket
	n, err := syscall.Read(c.fd, buf)
	if err != nil {
		if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
			return nil
		}
		return err
	}

	if n == 0 {
		return syscall.ECONNRESET // Client closed the connection
	}

	c.buffer = append(c.buffer, buf[:n]...)

	// Process all complete RESP commands currently in the buffer
	for {
		if len(c.buffer) == 0 {
			break
		}

		val, consumed, err := resp.DecodeFromBytes(c.buffer)
		if err != nil {
			if err == resp.ErrIncomplete {
				break // Wait for more data in the next EPOLLIN event
			}
			return err // Protocol error, disconnect client
		}

		// Remove the parsed command from the buffer
		c.buffer = c.buffer[consumed:]

		// Execute the command if it's an Array
		if val.Type == resp.Array {
			logger.Logger.Printf(
				"received command from fd %d: %+v",
				c.fd,
				val,
			)

			response := command.Dispatch(val)
			c.Write(response)
		}
	}

	return nil
}

func (c *AsyncConnection) Write(data []byte) {
	var total int
	for total < len(data) {
		n, err := syscall.Write(c.fd, data[total:])
		if err != nil {
			if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
				continue
			}
			return
		}
		total += n
	}
}
