package kcoin

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
)

type TransactionManager struct {
	tx *Transaction
}

func NewTransactionManager(transaction *Transaction) *TransactionManager {
	return &TransactionManager{
		tx: transaction,
	}
}

func (tManager *TransactionManager) Sign(prKey ecdsa.PrivateKey) {
	if tManager.tx.IsCoinbase() {
		return
	}

	txCopy := tManager.trimmedCopy()
	for i := 0; i < len(txCopy.Inputs); i++ {
		r, s, err := ecdsa.Sign(rand.Reader, &prKey, txCopy.ID)
		panicIfErrNotNil(err)

		tManager.tx.Inputs[i].Signature = append(r.Bytes(), s.Bytes()...)
	}
}

func (tManager *TransactionManager) Verify(blockchain *Blockchain) bool {
	for _, input := range tManager.tx.Inputs {
		if tManager.tx.IsCoinbase() || !isInputValid(input, blockchain) {
			return false
		}

		txCopy := tManager.trimmedCopy()
		if bytes.Compare(tManager.tx.ID, txCopy.ID) != 0 {
			return false
		}

		r, s := input.GetSignature()
		if !ecdsa.Verify(input.PubKey, tManager.tx.ID, r, s) {
			return false
		}
	}

	return true
}

func (tManager *TransactionManager) trimmedCopy() Transaction {
	var inputs []TransactionInput
	var outputs []TransactionOutput

	for _, input := range tManager.tx.Inputs {
		inputs = append(inputs, TransactionInput{
			TxID:        input.TxID,
			OutputIndex: input.OutputIndex,
			PubKey:      input.PubKey,
			Signature:   nil,
		})
	}

	for _, output := range tManager.tx.Outputs {
		outputs = append(outputs, TransactionOutput{
			Value:      output.Value,
			PubKeyHash: output.PubKeyHash,
		})
	}

	txCopy := Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}

	txCopy.SetID()

	return txCopy
}

func isInputValid(input TransactionInput, blockchain *Blockchain) bool {
	prevTx, err := findTransactionByID(input.TxID, blockchain)
	if err != nil {
		return false
	}

	outputPubKey := prevTx.Outputs[input.OutputIndex].PubKeyHash
	inputPubKeyHash := hashPubKey(*input.PubKey)
	if bytes.Compare(outputPubKey, inputPubKeyHash) != 0 {
		return false
	}

	if isOutputSpent(input.TxID, input.OutputIndex, blockchain) {
		return false
	}

	return true
}

func findTransactionByID(txID []byte, blockchain *Blockchain) (Transaction, error) {
	blockchainIterator := blockchain.Iterator()
	for {
		block := blockchainIterator.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(txID, tx.ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Can't find transaction")
}

func isOutputSpent(txID []byte, outputIndex int, blockchain *Blockchain) bool {
	blockchainIterator := blockchain.Iterator()
	for {
		block := blockchainIterator.Next()

		for _, tx := range block.Transactions {
			for _, input := range tx.Inputs {
				if bytes.Compare(input.TxID, txID) == 0 &&
					input.OutputIndex == outputIndex {
					return true
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return false
}
