package random

import (
	"math/rand/v2"
)

const (
	BlessingSize  = 50
	BlessingRange = 30
)

func GetLengthToAdd(canBeBlessed bool) (lengthToAdd int64, wasBlessed bool) {
	if canBeBlessed && rand.Int64N(BlessingRange) == 0 { // once per month on average
		return BlessingSize, true
	}

	lengthToAdd = rand.Int64N(16) - 5 // [-5,10]
	for lengthToAdd == 0 {
		lengthToAdd = rand.Int64N(16) - 5
	}
	return
}
