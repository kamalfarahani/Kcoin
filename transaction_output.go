package kcoin

import (
	"bytes"
)

type TransactionOutput struct {
	Value      int
	PubKeyHash []byte
}

func (txOutput *TransactionOutput) LockWithAddress(address []byte) error {
	pubKeyHash, err := getPubKeyHashFromAddress(address)
	if err != nil {
		return err
	}

	txOutput.PubKeyHash = pubKeyHash
	return nil
}

func (txOutput *TransactionOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(txOutput.PubKeyHash, pubKeyHash) == 0
}
