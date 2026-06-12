package kcoin

import (
	"bytes"
	"encoding/gob"
)

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
