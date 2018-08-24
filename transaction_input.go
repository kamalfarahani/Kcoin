package kcoin

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
)

type TransactionInput struct {
	TxID                []byte
	OutputIndex         int
	XYAppendedPubKey    []byte
	XYAppendedSignature []byte
}

func (txInput *TransactionInput) GetSignature() (*big.Int, *big.Int) {
	r := big.Int{}
	s := big.Int{}
	sigLen := len(txInput.XYAppendedSignature)
	r.SetBytes(txInput.XYAppendedSignature[:(sigLen / 2)])
	s.SetBytes(txInput.XYAppendedSignature[(sigLen / 2):])

	return &r, &s
}

func (txInput *TransactionInput) GetPublickKey() *ecdsa.PublicKey {
	curve := elliptic.P256()
	x := big.Int{}
	y := big.Int{}
	keyLen := len(txInput.XYAppendedPubKey)
	x.SetBytes(txInput.XYAppendedPubKey[:(keyLen / 2)])
	y.SetBytes(txInput.XYAppendedPubKey[(keyLen / 2):])

	return &ecdsa.PublicKey{
		Curve: curve,
		X:     &x,
		Y:     &y,
	}
}

func (txInput *TransactionInput) UsesKey(pubKeyHash []byte) bool {
	if len(txInput.XYAppendedPubKey) == 0 {
		return false
	}

	lockingHash := hashPubKey(*txInput.GetPublickKey())

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func getXYAppendedPubKey(pubKey ecdsa.PublicKey) []byte {
	return append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)
}
