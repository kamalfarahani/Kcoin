package kcoin

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
)

type TransactionInput struct {
	TxID        []byte
	OutputIndex int
	PubKey      *ecdsa.PublicKey
	Signature   []byte
}

func (txInput *TransactionInput) GetSignature() (*big.Int, *big.Int) {
	r := big.Int{}
	s := big.Int{}
	sigLen := len(txInput.Signature)
	r.SetBytes(txInput.Signature[:(sigLen / 2)])
	s.SetBytes(txInput.Signature[(sigLen / 2):])

	return &r, &s
}

func (txInput *TransactionInput) UsesKey(pubKeyHash []byte) bool {
	if txInput.PubKey == nil {
		return false
	}

	lockingHash := hashPubKey(*txInput.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
