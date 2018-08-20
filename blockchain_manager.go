package kcoin

import (
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

func SendCoin(fromAddress, toAddress string, amount int) {
	blockchain := NewBlockchain(fromAddress)
	bManager := NewBlockchainManager(blockchain)

	tx, err := bManager.NewUTXOTransaction(fromAddress, toAddress, amount)
	if err != nil {
		log.Println(err)
		return
	}

	blockchain.AddBlock([]*Transaction{tx})
}

func (bManager *BlockchainManager) NewUTXOTransaction(
	fromAddress, toAddress string, amount int) (*Transaction, error) {
	if bManager.GetBalance(fromAddress) < amount {
		return nil, errors.New("ERROR: Not enough funds")
	}

	var inputs []TransactionInput
	accumulated := 0
	txToUnspentOutputIndexs := bManager.FindUnspentTransactions(fromAddress)

	for tx, unspentIndexes := range txToUnspentOutputIndexs {
		newInputs, newAccumulated := makeInputsFromUnspentTxOutputs(
			fromAddress, amount, *tx, unspentIndexes, accumulated)

		inputs = append(inputs, newInputs...)
		accumulated = newAccumulated
		if accumulated >= amount {
			break
		}
	}

	outputs := makeOutputsForUTXOTransaction(
		fromAddress, toAddress, amount, accumulated-amount)

	return NewTransaction(inputs, outputs), nil
}

func (bManager *BlockchainManager) GetBalance(address string) int {
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

func (bManager *BlockchainManager) FindUnspentTransactions(address string) map[*Transaction][]uint {
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
	address string,
	transaction Transaction,
	txIDToSpentOutputIndexes map[string][]uint) []uint {

	var unspentOutputIndexes []uint
	txID := hex.EncodeToString(transaction.ID)

	for outputIndex, output := range transaction.Outputs {
		if !containsUint(txIDToSpentOutputIndexes[txID], uint(outputIndex)) {
			if output.CanBeUnlockedWith(address) {
				unspentOutputIndexes = append(unspentOutputIndexes, uint(outputIndex))
			}
		}
	}

	return unspentOutputIndexes
}

func isInputSpent(address string, input TransactionInput) bool {
	if input.CanUnlockOutputWith(address) {
		return true
	}

	return false
}

func makeInputsFromUnspentTxOutputs(
	address string,
	amount int,
	tx Transaction,
	unspentIndexes []uint,
	accumulated int) ([]TransactionInput, int) {

	var resultInputs []TransactionInput
	resultAccumulated := accumulated

	for outputIndex, output := range tx.Outputs {
		if containsUint(unspentIndexes, uint(outputIndex)) {
			newInput := TransactionInput{tx.ID, outputIndex, address}
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
	fromAddress, toAddress string,
	amount, exchange int) []TransactionOutput {

	resultOutputs := []TransactionOutput{
		TransactionOutput{
			Value:        amount,
			ScriptPubKey: toAddress,
		},
	}

	if exchange != 0 {
		resultOutputs = append(
			resultOutputs, TransactionOutput{
				Value:        exchange,
				ScriptPubKey: fromAddress,
			})
	}

	return resultOutputs
}
