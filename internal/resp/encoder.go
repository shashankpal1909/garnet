package resp

import "fmt"

func EncodeSimpleString(value string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", value))
}

func EncodeError(value string) []byte {
	return []byte(fmt.Sprintf("-%s\r\n", value))
}

func EncodeBulkString(value string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value))
}
