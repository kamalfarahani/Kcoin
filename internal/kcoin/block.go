package kcoin

import (
	"bytes"
	"crypto/sha256"
	"time"
)

const genesisCoinbaseData = "Kamal Genesis Block"

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int64
}

func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte

	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	resultHash := sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return resultHash[:]
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
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

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}
