package main

import (
	"fmt"

	"github.com/kamalfarahani/Kcoin/internal/kcoin"
)

func main() {
	wallet := kcoin.NewWallet()
	address := wallet.GetAddress()
	fmt.Printf("New wallet address: %s\n", address)

	blockchain := kcoin.NewBlockchain(address)
	manager := kcoin.NewBlockchainManager(blockchain)

	balance := manager.GetBalance(address)
	fmt.Printf("Balance of %s: %d\n", address, balance)
}
