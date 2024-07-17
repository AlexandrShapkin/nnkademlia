package nnkademlia

import (
	"crypto/sha1"
	"math/big"
)

func IdFromString(str string) *big.Int {
	hash := sha1.Sum([]byte(str))

	id := big.NewInt(0)
	id.SetBytes(hash[:])
	return id
}

func DistanceBetween(from *big.Int, to *big.Int) *big.Int {
	return big.NewInt(0).Xor(from, to)
}