package random

import (
	"github.com/benkenobi3/dick-and-dot/internal/database/repository"
	"math/rand"
	"time"
)

const (
	dickTimeoutNs = time.Hour * 8
)

func GetNewLength(startLength int64) int64 {
	return startLength + rand.Int63n(16) - 5 // todo: should we avoid 0 value?
}

func TimeBeforeReadyToGrow(dick repository.Dick) *time.Duration {
	now := time.Now().UTC()
	ableToGrowAgainAt := dick.UpdatedAt.Add(dickTimeoutNs)
	var timeLeft *time.Duration
	if now.Before(ableToGrowAgainAt) {
		t := ableToGrowAgainAt.Sub(now)
		timeLeft = &t
	}

	return timeLeft
}
