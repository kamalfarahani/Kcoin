package kcoin

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"math"
	"math/big"
)

const targetBits = 24
const hashBits = 256

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1).Lsh(
		big.NewInt(1), uint(hashBits-targetBits))

	return &ProofOfWork{
		block:  block,
		target: target,
	}
}

func (proofOfWork *ProofOfWork) prepareData(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			proofOfWork.block.PervBlockHash,
			proofOfWork.block.Data,
			IntToBytes(proofOfWork.block.Timestamp),
			IntToBytes(int64(targetBits)),
			IntToBytes(nonce),
		}, []byte{})

	return data
}

func (proofOfWork *ProofOfWork) Mine() (int64, []byte, error) {
	var nonce int64
	for nonce = 0; nonce < math.MaxInt64; nonce++ {
		data := proofOfWork.prepareData(nonce)
		hash := sha256.Sum256(data)

		if isHashValid(hash[:], proofOfWork.target) {
			return nonce, hash[:], nil
		}
	}

	return nonce, []byte{}, errors.New("Can't Mine")
}

func (proofOfWork *ProofOfWork) Validate() bool {
	data := proofOfWork.prepareData(proofOfWork.block.Nonce)
	hash := sha256.Sum256(data)

	return isHashValid(hash[:], proofOfWork.target)
}

func isHashValid(hash []byte, target *big.Int) bool {
	var hashInt *big.Int
	hashInt.SetBytes(hash)

	if hashInt.Cmp(target) == -1 {
		return true
	}

	return false
}
