package random

import (
	"math"
	"math/rand"
)

func GetNewLength(startLength int64) int64 {
	linearLengthInc := randFloat(-5, 5)
	normalizedLength := sigmoidNormalization(linearLengthInc)
	return startLength + int64(math.Round(normalizedLength)) // todo: should we avoid 0 value?
}

// randFloat in range [min..max)
func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func sigmoidNormalization(input float64) float64 {
	return (3.5 / (0.225 + math.Pow(math.E-0.4, -0.9-input))) - 5
}
