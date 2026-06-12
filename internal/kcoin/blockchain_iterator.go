package kcoin

import "github.com/boltdb/bolt"

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (iterator *BlockchainIterator) Next() *Block {
	block := getBlock(iterator.db, iterator.currentHash)
	iterator.currentHash = block.PrevBlockHash

	return block
}
