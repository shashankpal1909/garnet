package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Decoder struct {
	reader *bufio.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{reader: bufio.NewReader(r)}
}

func (d *Decoder) Decode() (Value, error) {
	b, err := d.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch Type(b) {
	case Array:
		return d.decodeArray()
	case BulkString:
		return d.decodeBulkString()
	case SimpleString:
		return d.decodeSimpleString()
	case Error:
		return d.decodeError()
	case Integer:
		return d.decodeInteger()
	default:
		return Value{}, fmt.Errorf("unknown RESP type: %q", b)
	}
}

func (d *Decoder) readLine() (line []byte, err error) {
	for {
		b, err := d.reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == '\r' {
			_, err = d.reader.ReadByte() // consume \n
			if err != nil {
				return nil, err
			}
			break
		}
		line = append(line, b)
	}
	return line, nil
}

func (d *Decoder) decodeArray() (Value, error) {
	line, err := d.readLine()
	if err != nil {
		return Value{}, err
	}
	count, err := strconv.Atoi(string(line))
	if err != nil {
		return Value{}, err
	}

	if count == -1 {
		return Value{Type: Array, Data: nil}, nil
	}

	array := make([]Value, count)
	for i := 0; i < count; i++ {
		val, err := d.Decode()
		if err != nil {
			return Value{}, err
		}
		array[i] = val
	}

	return Value{Type: Array, Data: array}, nil
}

func (d *Decoder) decodeBulkString() (Value, error) {
	line, err := d.readLine()
	if err != nil {
		return Value{}, err
	}
	length, err := strconv.Atoi(string(line))
	if err != nil {
		return Value{}, err
	}

	if length == -1 {
		return Value{Type: BulkString, Data: nil}, nil
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(d.reader, buf)
	if err != nil {
		return Value{}, err
	}

	// consume \r\n
	_, err = d.readLine()
	if err != nil {
		return Value{}, err
	}

	return Value{Type: BulkString, Data: buf}, nil
}

func (d *Decoder) decodeSimpleString() (Value, error) {
	line, err := d.readLine()
	if err != nil {
		return Value{}, err
	}
	return Value{Type: SimpleString, Data: string(line)}, nil
}

func (d *Decoder) decodeError() (Value, error) {
	line, err := d.readLine()
	if err != nil {
		return Value{}, err
	}
	return Value{Type: Error, Data: string(line)}, nil
}

func (d *Decoder) decodeInteger() (Value, error) {
	line, err := d.readLine()
	if err != nil {
		return Value{}, err
	}
	num, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return Value{}, err
	}
	return Value{Type: Integer, Data: num}, nil
}
