package kcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
)

const subsidy = 10

type Transaction struct {
	ID      []byte
	Inputs  []TransactionInput
	Outputs []TransactionOutput
}

func (transaction *Transaction) SetID() {
	var encodeWriter bytes.Buffer
	encoder := gob.NewEncoder(&encodeWriter)
	err := encoder.Encode(transaction)
	panicIfErrNotNil(err)

	hash := sha256.Sum256(encodeWriter.Bytes())
	transaction.ID = hash[:]
}

func (transaction *Transaction) IsCoinbase() bool {
	return len(transaction.Inputs) == 1 &&
		len(transaction.Inputs[0].TxID) == 0 &&
		transaction.Inputs[0].OutputIndex == -1
}

func NewCoinbaseTransaction(to, data string) *Transaction {
	txIn := TransactionInput{
		TxID:        []byte{},
		OutputIndex: -1,
		ScriptSig:   data,
	}

	txOut := TransactionOutput{
		Value:        subsidy,
		ScriptPubKey: to,
	}

	transaction := Transaction{
		ID:      nil,
		Inputs:  []TransactionInput{txIn},
		Outputs: []TransactionOutput{txOut},
	}

	transaction.SetID()

	return &transaction
}

func GetBalance(blockchain Blockchain, address string) int {
	balance := 0
	txToUnspentOutputIndexs := FindUnspentTransactions(blockchain, address)
	for tx, unspentIndexs := range txToUnspentOutputIndexs {
		for outpuIndex, output := range tx.Outputs {
			if containsUint(unspentIndexs, uint(outpuIndex)) {
				balance += output.Value
			}
		}
	}

	return balance
}

func FindUnspentTransactions(blockchain Blockchain, address string) map[*Transaction][]uint {
	txToUnspentOutputIndexs := make(map[*Transaction][]uint)
	txIDToSpentOutputIndexes := make(map[string][]int)
	blockchainIterator := blockchain.Iterator()

	for {
		block := blockchainIterator.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			unspentTxOutputs := findTransactionUnspentOutputIndexes(
				address, *tx, txIDToSpentOutputIndexes)

			txToUnspentOutputIndexs[tx] = append(
				txToUnspentOutputIndexs[tx],
				unspentTxOutputs...)

			txIDToSpentOutputIndexes[txID] = append(
				txIDToSpentOutputIndexes[txID],
				findTransactionSpentOutputIndexes(address, *tx)...)
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
	txIDToSpentOutputIndexes map[string][]int) []uint {

	var unspentOutputIndexes []uint
	txID := hex.EncodeToString(transaction.ID)

	for outputIndex, output := range transaction.Outputs {
		if !isOutputSpent(outputIndex, txIDToSpentOutputIndexes[txID]) {
			if output.CanBeUnlockedWith(address) {
				unspentOutputIndexes = append(unspentOutputIndexes, uint(outputIndex))
			}
		}
	}

	return unspentOutputIndexes
}

func isOutputSpent(outputIndex int, spentOutputIndexes []int) bool {
	if spentOutputIndexes == nil {
		return false
	}

	for _, spentOutIndex := range spentOutputIndexes {
		if spentOutIndex == outputIndex {
			return true
		}
	}

	return false
}

func findTransactionSpentOutputIndexes(address string, transaction Transaction) []int {
	if transaction.IsCoinbase() {
		return []int{}
	}

	var spentOutputIndexes []int

	for _, inputTx := range transaction.Inputs {
		if inputTx.CanUnlockOutputWith(address) {
			spentOutputIndexes = append(spentOutputIndexes, inputTx.OutputIndex)
		}
	}

	return spentOutputIndexes
}
