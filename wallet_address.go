package kcoin

import (
	"crypto/ecdsa"
	"crypto/sha256"

	"github.com/itchyny/base58-go"
	"golang.org/x/crypto/ripemd160"
)

func (wallet *Wallet) GetAddress() []byte {
	pubKeyHash := hashPubKey(wallet.PrivateKey.PublicKey)
	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := calculateChecksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)

	encoder := base58.BitcoinEncoding
	encoded, err := encoder.Encode(fullPayload)
	panicIfErrNotNil(err)

	return encoded
}

func hashPubKey(pubKey ecdsa.PublicKey) []byte {
	xyAppendedPubKey := append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)
	pubKeySHA256 := sha256.Sum256(xyAppendedPubKey)

	ripemdHasher := ripemd160.New()
	_, err := ripemdHasher.Write(pubKeySHA256[:])
	panicIfErrNotNil(err)
	pubKeyRIPEMD160 := ripemdHasher.Sum(nil)

	return pubKeyRIPEMD160
}

func calculateChecksum(payload []byte) []byte {
	hashFunc := sha256.Sum256
	firstHash := hashFunc(payload)
	secondHash := hashFunc(firstHash[:])

	return secondHash[:addressChecksumLen]
}
