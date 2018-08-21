package kcoin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
}

func NewWallet() *Wallet {
	return &Wallet{
		PrivateKey: newKey(),
	}
}

func newKey() ecdsa.PrivateKey {
	curve := elliptic.P256()
	prKey, err := ecdsa.GenerateKey(curve, rand.Reader)

	panicIfErrNotNil(err)

	return *prKey
}
