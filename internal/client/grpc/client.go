package grpc

import (
	"context"
	"time"
)

type RateClient interface {
	FindRate(ctx context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber string) (int64, float64, time.Duration, error)
	FindSupRates(ctx context.Context, gwgrIds []int64, dateAt int64, aNumber, bNumber string) (int64, time.Duration, error)
}
