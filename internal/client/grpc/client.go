package grpc

import (
	"context"
	"time"

	"github.com/VadimGossip/drs_storage_tester/internal/model"
)

type RateClient interface {
	FindRate(ctx context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber uint64) (model.RateBase, time.Duration, error)
	FindSupRates(ctx context.Context, gwgrIds []int64, dateAt int64, aNumber, bNumber uint64) (map[int64]model.RateBase, time.Duration, error)
}
