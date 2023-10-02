package random

import "math/rand"

func GetNewLength(startLength int64) int64 {
	return startLength + rand.Int63n(16) - 5 // todo: should we avoid 0 value?
}
