package kcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"log"
)

const subsidy = 10

type Transaction struct {
	ID      []byte
	Inputs  []TransactionInput
	Outputs []TransactionOutput
}

func NewTransaction(inputs []TransactionInput, outputs []TransactionOutput) *Transaction {
	transaction := &Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}

	transaction.SetID()

	return transaction
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

func SendCoin(fromAddress, toAddress string, amount int) {
	blockchain := NewBlockchain(fromAddress)
	defer blockchain.db.Close()

	tx, err := NewUTXOTransaction(*blockchain, fromAddress, toAddress, amount)
	if err != nil {
		log.Println(err)
		return
	}

	blockchain.AddBlock([]*Transaction{tx})
}

func NewUTXOTransaction(
	blockchain Blockchain,
	fromAddress, toAddress string, amount int) (*Transaction, error) {

	if GetBalance(blockchain, fromAddress) < amount {
		return nil, errors.New("ERROR: Not enough funds")
	}

	var inputs []TransactionInput
	var outputs []TransactionOutput
	accumulated := 0
	txToUnspentOutputIndexs := FindUnspentTransactions(blockchain, fromAddress)

Work:
	for tx, unspentIndexes := range txToUnspentOutputIndexs {
		for outputIndex, output := range tx.Outputs {
			if containsUint(unspentIndexes, uint(outputIndex)) {
				newInput := TransactionInput{tx.ID, outputIndex, fromAddress}
				inputs = append(inputs, newInput)
				accumulated += output.Value
			}

			if accumulated >= amount {
				break Work
			}
		}
	}

	outputs = append(outputs, TransactionOutput{amount, toAddress})
	if accumulated > amount {
		outputs = append(outputs, TransactionOutput{accumulated - amount, fromAddress})
	}

	return NewTransaction(inputs, outputs), nil
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
			unspentTxOutputs := findTransactionUnspentOutputIndexes(
				address, *tx, txIDToSpentOutputIndexes)

			txToUnspentOutputIndexs[tx] = append(
				txToUnspentOutputIndexs[tx],
				unspentTxOutputs...)

			for _, input := range tx.Inputs {
				if isInputSpent(address, input) {
					inputTxID := hex.EncodeToString(input.TxID)
					txIDToSpentOutputIndexes[inputTxID] = append(
						txIDToSpentOutputIndexes[inputTxID],
						input.OutputIndex)
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

func isInputSpent(address string, input TransactionInput) bool {
	if input.CanUnlockOutputWith(address) {
		return true
	}

	return false
}
