package kcoin

type Blockchain struct {
	blocks []*Block
}

func (blockchain *Blockchain) Add(data string) {
	prevBlock := blockchain.blocks[len(blockchain.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	blockchain.blocks = append(blockchain.blocks, newBlock)
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		blocks: []*Block{NewGenesisBlock()},
	}
}
