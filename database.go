package kcoin

import "github.com/boltdb/bolt"

const dbFilePath = "./assets/blockchain.db"
const dbFileMode = 0600
const blocksBucketName = "blocks"
const lastHashKey = "l"

func createDBIfNotExist(address string) (*bolt.DB, []byte) {
	var tip []byte
	db, outErr := bolt.Open(dbFilePath, dbFileMode, nil)

	outErr = db.Update(func(tx *bolt.Tx) error {
		blocksBucket := tx.Bucket([]byte(blocksBucketName))
		var inErr error

		if blocksBucket == nil {
			coinbaseTx := NewCoinbaseTransaction(address, genesisCoinbaseData)
			genesis := NewGenesisBlock(coinbaseTx)
			blocksBucket, inErr = tx.CreateBucket([]byte(blocksBucketName))
			inErr = blocksBucket.Put(genesis.Hash, SerializeBlock(genesis))
			inErr = blocksBucket.Put([]byte(lastHashKey), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = blocksBucket.Get([]byte(lastHashKey))
		}

		return inErr
	})

	panicIfErrNotNil(outErr)

	return db, tip
}

func getBlock(db *bolt.DB, hash []byte) *Block {
	var block *Block

	err := db.View(func(tx *bolt.Tx) error {
		blocksBucket := tx.Bucket([]byte(blocksBucketName))
		encodedBlock := blocksBucket.Get(hash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	panicIfErrNotNil(err)

	return block
}

func getLastBlockHash(db *bolt.DB) []byte {
	var lastHash []byte
	err := db.View(func(tx *bolt.Tx) error {
		blocksBucket := tx.Bucket([]byte(blocksBucketName))
		lastHash = blocksBucket.Get([]byte(lastHashKey))

		return nil
	})

	panicIfErrNotNil(err)

	return lastHash
}

func addBlockToDB(db *bolt.DB, newBlock *Block) {
	outErr := db.Update(func(tx *bolt.Tx) error {
		blocksBucket := tx.Bucket([]byte(blocksBucketName))
		inErr := blocksBucket.Put(newBlock.Hash, SerializeBlock(newBlock))
		inErr = blocksBucket.Put([]byte(lastHashKey), newBlock.Hash)

		return inErr
	})

	panicIfErrNotNil(outErr)
}
