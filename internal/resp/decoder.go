package resp

import (
	"bytes"
	"errors"
	"io"
	"strconv"
)

// ErrIncomplete is returned when the byte slice doesn't contain a full RESP payload.
// The caller should buffer more bytes and try again.
var ErrIncomplete = errors.New("incomplete resp payload")

// DecodeFromBytes attempts to parse exactly one complete RESP value from the provided byte slice.
// It returns the parsed Value, the number of bytes consumed, and any error.
func DecodeFromBytes(data []byte) (Value, int, error) {
	if len(data) == 0 {
		return Value{}, 0, ErrIncomplete
	}

	b := data[0]
	switch Type(b) {
	case Array:
		return decodeArray(data)
	case BulkString:
		return decodeBulkString(data)
	case SimpleString:
		return decodeSimpleString(data)
	case Error:
		return decodeError(data)
	case Integer:
		return decodeInteger(data)
	default:
		return Value{}, 0, errors.New("unknown RESP type")
	}
}

func readLine(data []byte, offset int) (line []byte, newOffset int, err error) {
	idx := bytes.IndexByte(data[offset:], '\r')
	if idx == -1 {
		return nil, 0, ErrIncomplete
	}
	// Check if \n exists after \r
	if offset+idx+1 >= len(data) {
		return nil, 0, ErrIncomplete
	}
	if data[offset+idx+1] != '\n' {
		return nil, 0, errors.New("invalid line ending")
	}

	line = data[offset : offset+idx]
	return line, offset + idx + 2, nil
}

func decodeArray(data []byte) (Value, int, error) {
	line, offset, err := readLine(data, 1)
	if err != nil {
		return Value{}, 0, err
	}

	count, err := strconv.Atoi(string(line))
	if err != nil {
		return Value{}, 0, err
	}

	if count == -1 {
		return Value{Type: Array, Data: nil}, offset, nil
	}

	array := make([]Value, count)
	for i := 0; i < count; i++ {
		val, consumed, err := DecodeFromBytes(data[offset:])
		if err != nil {
			return Value{}, 0, err
		}
		array[i] = val
		offset += consumed
	}

	return Value{Type: Array, Data: array}, offset, nil
}

func decodeBulkString(data []byte) (Value, int, error) {
	line, offset, err := readLine(data, 1)
	if err != nil {
		return Value{}, 0, err
	}

	length, err := strconv.Atoi(string(line))
	if err != nil {
		return Value{}, 0, err
	}

	if length == -1 {
		return Value{Type: BulkString, Data: nil}, offset, nil
	}

	if offset+length+2 > len(data) {
		return Value{}, 0, ErrIncomplete
	}

	buf := make([]byte, length)
	copy(buf, data[offset:offset+length])
	offset += length

	if data[offset] != '\r' || data[offset+1] != '\n' {
		return Value{}, 0, errors.New("invalid bulk string ending")
	}
	offset += 2

	return Value{Type: BulkString, Data: buf}, offset, nil
}

func decodeSimpleString(data []byte) (Value, int, error) {
	line, offset, err := readLine(data, 1)
	if err != nil {
		return Value{}, 0, err
	}
	return Value{Type: SimpleString, Data: string(line)}, offset, nil
}

func decodeError(data []byte) (Value, int, error) {
	line, offset, err := readLine(data, 1)
	if err != nil {
		return Value{}, 0, err
	}
	return Value{Type: Error, Data: string(line)}, offset, nil
}

func decodeInteger(data []byte) (Value, int, error) {
	line, offset, err := readLine(data, 1)
	if err != nil {
		return Value{}, 0, err
	}
	num, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return Value{}, 0, err
	}
	return Value{Type: Integer, Data: num}, offset, nil
}

// Decoder provides a blocking interface for reading RESP values from an io.Reader.
// It wraps the non-blocking DecodeFromBytes logic.
type Decoder struct {
	reader io.Reader
	buffer []byte
}

// NewDecoder creates a new decoder wrapping the provided io.Reader.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		reader: r,
		buffer: make([]byte, 0, 4096),
	}
}

// Decode blocks until a full RESP value can be parsed from the reader.
func (d *Decoder) Decode() (Value, error) {
	buf := make([]byte, 4096)

	for {
		// Attempt to parse from the current buffer
		if len(d.buffer) > 0 {
			val, consumed, err := DecodeFromBytes(d.buffer)
			if err == nil {
				// Successfully parsed a value!
				// Remove the consumed bytes from our buffer
				d.buffer = d.buffer[consumed:]
				return val, nil
			}
			if err != ErrIncomplete {
				// Fatal protocol error
				return Value{}, err
			}
			// If err == ErrIncomplete, we just need more data, so we continue to read.
		}

		// Read more data from the connection (this blocks)
		n, err := d.reader.Read(buf)
		if err != nil {
			return Value{}, err // Returns io.EOF on disconnect
		}

		// Append the new data to our buffer and loop to try parsing again
		d.buffer = append(d.buffer, buf[:n]...)
	}
}
