package kcoin

import (
	"bytes"
	"crypto/ecdsa"
)

type TransactionInput struct {
	TxID        []byte
	OutputIndex int
	PubKey      *ecdsa.PublicKey
	Signature   []byte
}

func (txInput *TransactionInput) UsesKey(pubKeyHash []byte) bool {
	if txInput.PubKey == nil {
		return false
	}

	lockingHash := hashPubKey(*txInput.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
