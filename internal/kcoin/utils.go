package kcoin

import (
	"bytes"
	"encoding/binary"
	"log"
)

// IntToBytes converts an int64 to a byte array
func IntToBytes(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func containsUint(slice []uint, elemet uint) bool {
	for _, a := range slice {
		if a == elemet {
			return true
		}
	}
	return false
}

// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func panicIfErrNotNil(err error) {
	if err != nil {
		panic(err)
	}
}
