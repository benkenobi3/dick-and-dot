package random

import (
	"crypto/rand"
	"math/big"
	mrand "math/rand"
)

const (
	BlessingSize = 50
)

func GetNewLength(startLength int64, canBeBlessed bool) (newLength int64, wasBlessed bool) {
	blessingRand, _ := rand.Int(rand.Reader, big.NewInt(30))
	if canBeBlessed && blessingRand.Int64() == 0 { // once per month on average
		return startLength + BlessingSize, true
	}

	defaultRand, _ := rand.Int(rand.Reader, big.NewInt(16))
	lengthToAdd := defaultRand.Int64() - 5
	if lengthToAdd == 0 {
		if mrand.Int()%2 == 0 {
			lengthToAdd++
		} else {
			lengthToAdd--
		}
	}
	return startLength + lengthToAdd, false
}
