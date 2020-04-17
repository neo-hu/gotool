package limiter

import (
	"context"
	"golang.org/x/time/rate"
	"time"
)

type Rate interface {
	WaitN(ctx context.Context, n int) error
	Burst() int
}

func NewRate(bytesPerSec, burstLimit int) Rate {
	limiter := rate.NewLimiter(rate.Limit(bytesPerSec), burstLimit)
	limiter.AllowN(time.Now(), burstLimit) // spend initial burst
	return limiter
}