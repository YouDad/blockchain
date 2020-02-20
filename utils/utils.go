package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"reflect"
)

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func intToByte(num int) []byte {
	var bytes [4]byte

	bytes[0] = byte(num >> 0)
	if num < 1<<8 {
		return bytes[:1]
	}
	bytes[1] = byte(num >> 8)
	if num < 1<<16 {
		return bytes[:2]
	}
	bytes[2] = byte(num >> 16)
	if num < 1<<24 {
		return bytes[:3]
	}
	bytes[3] = byte(num >> 24)
	return bytes[:4]
}

func InterfaceIsNil(i interface{}) bool {
	defer func() {
		recover()
	}()
	vi := reflect.ValueOf(i)
	return vi.IsNil()
}
