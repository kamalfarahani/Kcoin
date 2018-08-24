package kcoin

import (
	"crypto/sha256"
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
	hash := sha256.Sum256(SerializeTransaction(*transaction))
	transaction.ID = hash[:]
}

func (transaction *Transaction) IsCoinbase() bool {
	return len(transaction.Inputs) == 1 &&
		len(transaction.Inputs[0].TxID) == 0 &&
		transaction.Inputs[0].OutputIndex == -1
}

func NewCoinbaseTransaction(toAddress []byte, data string) *Transaction {
	txIn := TransactionInput{
		TxID:                []byte{},
		OutputIndex:         -1,
		XYAppendedPubKey:    []byte(data),
		XYAppendedSignature: nil,
	}

	//Handle err later
	pubKeyHash, _ := getPubKeyHashFromAddress(toAddress)
	txOut := TransactionOutput{
		Value:      subsidy,
		PubKeyHash: pubKeyHash,
	}

	return NewTransaction(
		[]TransactionInput{txIn},
		[]TransactionOutput{txOut})
}
