package service

import (
	"context"
	"github.com/VadimGossip/drs_storage_tester/internal/model"
	"time"
)

type RateService interface {
	FindRate(ctx context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber uint64) (model.RateBase, time.Duration, error)
	FindSupRates(ctx context.Context, gwgrIds []int64, dateAt int64, aNumber, bNumber uint64) (map[int64]model.RateBase, time.Duration, error)
}

type RequestService interface {
	GetTaskRequests(ctx context.Context, limit int64) ([]model.TaskRequest, error)
	GetSupGwgrIds(ctx context.Context) ([]int64, error)
}

type EventService interface {
	RunEventGeneration(ctx context.Context, total, rps, pps int) chan int
}

type DataService interface {
	GetTaskRequest() model.TaskRequest
	GetSupGwgrIds() []int64
	Refresh(ctx context.Context, limit int64) error
}

type ImitatorService interface {
	RunTests(ctx context.Context) error
}
