package resp

import "fmt"

var RESP_NIL = []byte("$-1\r\n")
var RESP_OK = []byte("+OK\r\n")

func EncodeSimpleString(value string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", value))
}

func EncodeError(value string) []byte {
	return []byte(fmt.Sprintf("-%s\r\n", value))
}

func EncodeBulkString(value string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value))
}

func EncodeInteger(value int64) []byte {
	return []byte(fmt.Sprintf(":%d\r\n", value))
}
