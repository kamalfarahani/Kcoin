package kcoin

import (
	"github.com/boltdb/bolt"
)

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (blockchain *Blockchain) AddBlock(transactions []*Transaction) {
	lastHash := getLastBlockHash(blockchain.db)
	newBlock := NewBlock(transactions, lastHash)
	addBlockToDB(blockchain.db, newBlock)
	blockchain.tip = newBlock.Hash
}

func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		currentHash: blockchain.tip,
		db:          blockchain.db,
	}
}

//Fix this later
func NewBlockchain(address string) *Blockchain {
	db, tip := createDBIfNotExist(address)

	return &Blockchain{
		tip: tip,
		db:  db,
	}
}
