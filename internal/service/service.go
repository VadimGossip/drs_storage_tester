package service

import (
	"context"
	"time"
)

type RateService interface {
	FindRate(ctx context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber uint64) (int64, float64, time.Duration, error)
	FindSupRates(ctx context.Context, gwgrIds []int64, dateAt int64, aNumber, bNumber uint64) (int64, time.Duration, error)
}
