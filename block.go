package kcoin

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

const genesisData = "Genesis Block"

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int64
}

func (block *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	headers := bytes.Join(
		[][]byte{
			block.PrevBlockHash,
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
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}

	proofOfWork := NewProofOfWork(block)
	nonce, hash, err := proofOfWork.Mine()

	if err != nil {
		panic(err)
	}

	block.Hash = hash
	block.Nonce = nonce

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock(genesisData, []byte{})
}
