package kcoin

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
)

func init() {
	gob.Register(elliptic.P256())
}

func SerializeTransaction(transaction Transaction) []byte {
	var encodeWriter bytes.Buffer
	encoder := gob.NewEncoder(&encodeWriter)
	err := encoder.Encode(transaction)
	panicIfErrNotNil(err)

	return encodeWriter.Bytes()
}

func DeserializeTransaction(encodedData []byte) Transaction {
	var resultBlock Transaction

	decoder := gob.NewDecoder(bytes.NewReader(encodedData))
	err := decoder.Decode(&resultBlock)
	panicIfErrNotNil(err)

	return resultBlock
}
