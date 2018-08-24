package kcoin

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"

	"golang.org/x/crypto/ripemd160"
)

func (wallet *Wallet) GetAddress() []byte {
	pubKeyHash := hashPubKey(wallet.PrivateKey.PublicKey)
	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := calculateChecksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)

	encoded := Base58Encode(fullPayload)

	return encoded
}

func getPubKeyHashFromAddress(address []byte) ([]byte, error) {
	decodedAddress := Base58Decode(address)
	//fix it later
	if len(address) < addressChecksumLen+2 {
		return nil, errors.New("Invalid address")
	}

	return decodedAddress[1 : len(decodedAddress)-addressChecksumLen], nil
}

func hashPubKey(pubKey ecdsa.PublicKey) []byte {
	xyAppendedPubKey := getXYAppendedPubKey(pubKey)
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
