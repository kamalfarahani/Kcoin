package kcoin

import "github.com/boltdb/bolt"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (blockchain *Blockchain) AddBlock(data string) {
	lastHash := getLastBlockHash(blockchain.db)
	newBlock := NewBlock(data, lastHash)
	addBlockToDB(blockchain.db, newBlock)
	blockchain.tip = newBlock.Hash
}

func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		currentHash: blockchain.tip,
		db:          blockchain.db,
	}
}

func NewBlockchain() *Blockchain {
	db, tip := createDBIfNotExist()

	return &Blockchain{
		tip: tip,
		db:  db,
	}
}
