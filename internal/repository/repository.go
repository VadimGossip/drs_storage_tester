package repository

import (
	"context"
	"github.com/VadimGossip/drs_storage_tester/internal/model"
	"time"
)

type RateRepository interface {
	FindRate(ctx context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber string) (model.RateBase, time.Duration, error)
	FindSupRates(ctx context.Context, gwgrIds []int64, dateAt int64, aNumber, bNumber string) (map[int64]model.RateBase, time.Duration, error)
}

type RequestRepository interface {
	GetTaskRequests(ctx context.Context, limit int64) ([]model.TaskRequest, error)
	GetSupGwgrIds(ctx context.Context) ([]int64, error)
}
