package kcoin

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log"
)

type BlockchainManager struct {
	blockchain *Blockchain
}

func NewBlockchainManager(blockchain *Blockchain) *BlockchainManager {
	return &BlockchainManager{
		blockchain: blockchain,
	}
}

func (bManager *BlockchainManager) SendCoin(wallet Wallet, toAddress []byte, amount int) {
	tx, err := bManager.NewUTXOTransaction(wallet, toAddress, amount)
	if err != nil {
		log.Println(err)
		return
	}

	bManager.blockchain.AddBlock([]*Transaction{tx})
}

func (bManager *BlockchainManager) NewUTXOTransaction(
	wallet Wallet, toAddress []byte, amount int) (*Transaction, error) {
	if bManager.GetBalance(wallet.GetAddress()) < amount {
		return nil, errors.New("ERROR: Not enough funds")
	}

	var inputs []TransactionInput
	accumulated := 0
	txToUnspentOutputIndexs := bManager.FindUnspentTransactions(wallet.GetAddress())

	for tx, unspentIndexes := range txToUnspentOutputIndexs {
		newInputs, newAccumulated := makeInputsFromUnspentTxOutputs(
			wallet.PrivateKey.PublicKey, amount, *tx, unspentIndexes, accumulated)

		inputs = append(inputs, newInputs...)
		accumulated = newAccumulated
		if accumulated >= amount {
			break
		}
	}

	outputs := makeOutputsForUTXOTransaction(
		wallet.GetAddress(), toAddress, amount, accumulated-amount)

	txManager := NewTransactionManager(NewTransaction(inputs, outputs))
	txManager.Sign(wallet.PrivateKey)
	return txManager.tx, nil
}

func (bManager *BlockchainManager) GetBalance(address []byte) int {
	balance := 0
	txToUnspentOutputIndexs := bManager.FindUnspentTransactions(address)
	for tx, unspentIndexs := range txToUnspentOutputIndexs {
		for outpuIndex, output := range tx.Outputs {
			if containsUint(unspentIndexs, uint(outpuIndex)) {
				balance += output.Value
			}
		}
	}

	return balance
}

func (bManager *BlockchainManager) FindUnspentTransactions(address []byte) map[*Transaction][]uint {
	txToUnspentOutputIndexs := make(map[*Transaction][]uint)
	txIDToSpentOutputIndexes := make(map[string][]uint)
	blockchainIterator := bManager.blockchain.Iterator()

	for {
		block := blockchainIterator.Next()
		for _, tx := range block.Transactions {
			unspentTxOutputs := findTransactionUnspentOutputIndexes(
				address, *tx, txIDToSpentOutputIndexes)

			txToUnspentOutputIndexs[tx] = append(
				txToUnspentOutputIndexs[tx],
				unspentTxOutputs...)

			for _, input := range tx.Inputs {
				if isInputSpent(address, input) && !tx.IsCoinbase() {
					inputTxID := hex.EncodeToString(input.TxID)
					txIDToSpentOutputIndexes[inputTxID] = append(
						txIDToSpentOutputIndexes[inputTxID],
						uint(input.OutputIndex))
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return txToUnspentOutputIndexs
}

func findTransactionUnspentOutputIndexes(
	address []byte,
	transaction Transaction,
	txIDToSpentOutputIndexes map[string][]uint) []uint {

	var unspentOutputIndexes []uint
	txID := hex.EncodeToString(transaction.ID)

	for outputIndex, output := range transaction.Outputs {
		if !containsUint(txIDToSpentOutputIndexes[txID], uint(outputIndex)) {
			pubKeyHash, _ := getPubKeyHashFromAddress(address)
			if output.IsLockedWithKey(pubKeyHash) {
				unspentOutputIndexes = append(unspentOutputIndexes, uint(outputIndex))
			}
		}
	}

	return unspentOutputIndexes
}

func isInputSpent(address []byte, input TransactionInput) bool {
	pubKeyHash, _ := getPubKeyHashFromAddress(address)
	if input.UsesKey(pubKeyHash) {
		return true
	}

	return false
}

func makeInputsFromUnspentTxOutputs(
	pubKey ecdsa.PublicKey,
	amount int,
	tx Transaction,
	unspentIndexes []uint,
	accumulated int) ([]TransactionInput, int) {

	var resultInputs []TransactionInput
	resultAccumulated := accumulated

	for outputIndex, output := range tx.Outputs {
		if containsUint(unspentIndexes, uint(outputIndex)) {
			newInput := TransactionInput{tx.ID, outputIndex, &pubKey, nil}
			resultInputs = append(resultInputs, newInput)
			resultAccumulated += output.Value
		}

		if resultAccumulated >= amount {
			return resultInputs, resultAccumulated
		}
	}

	return resultInputs, resultAccumulated
}

func makeOutputsForUTXOTransaction(
	fromAddress, toAddress []byte,
	amount, exchange int) []TransactionOutput {

	toPubKey, _ := getPubKeyHashFromAddress(toAddress)
	resultOutputs := []TransactionOutput{
		TransactionOutput{
			Value:      amount,
			PubKeyHash: toPubKey,
		},
	}

	fromPubKey, _ := getPubKeyHashFromAddress(fromAddress)
	if exchange != 0 {
		resultOutputs = append(
			resultOutputs, TransactionOutput{
				Value:      exchange,
				PubKeyHash: fromPubKey,
			})
	}

	return resultOutputs
}
