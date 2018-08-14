package kcoin

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PervBlockHash []byte
	Hash          []byte
}

func (block *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	headers := bytes.Join(
		[][]byte{
			block.PervBlockHash,
			block.Data,
			timestamp,
		}, []byte{})

	hash := sha256.Sum256(headers)
	block.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PervBlockHash: prevBlockHash,
		Hash:          []byte{},
	}
	block.SetHash()

	return block
}
