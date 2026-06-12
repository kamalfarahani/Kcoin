package kcoin

import (
	"log"

	"github.com/boltdb/bolt"
)

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

func (blockchain *Blockchain) AddBlock(transactions []*Transaction) {
	for _, tx := range transactions {
		txManager := NewTransactionManager(tx)
		if !txManager.Verify(blockchain) {
			log.Println("ERROR: Invalid transaction")
			return
		}
	}

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
func NewBlockchain(address []byte) *Blockchain {
	db, tip := createDBIfNotExist(address)

	return &Blockchain{
		tip: tip,
		db:  db,
	}
}
