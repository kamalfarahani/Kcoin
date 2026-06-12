package kcoin

import (
	"bytes"
	"encoding/gob"
)

func SerializeBlock(block *Block) []byte {
	var resultWriter bytes.Buffer
	encoder := gob.NewEncoder(&resultWriter)

	err := encoder.Encode(block)
	panicIfErrNotNil(err)

	return resultWriter.Bytes()
}

func DeserializeBlock(encodedData []byte) *Block {
	var resultBlock Block

	decoder := gob.NewDecoder(bytes.NewReader(encodedData))
	err := decoder.Decode(&resultBlock)
	panicIfErrNotNil(err)

	return &resultBlock
}
