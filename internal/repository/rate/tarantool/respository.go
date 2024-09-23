package tarantool

import (
	"context"
	"fmt"
	"time"

	db "github.com/VadimGossip/drs_storage_tester/internal/client/db/tarantool"
	"github.com/tarantool/go-tarantool/v2"

	def "github.com/VadimGossip/drs_storage_tester/internal/repository"
)

var _ def.RateRepository = (*repository)(nil)

const (
	findRateFunc     string = "rates.find_rate"
	findSupRatesFunc string = "rates.find_sup_rates"
)

type repository struct {
	db db.Client
}

func NewRepository(db db.Client) *repository {
	return &repository{
		db: db,
	}
}
func (r *repository) FindRate(_ context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber string) (int64, float64, time.Duration, error) {
	ts := time.Now()
	resp, err := r.db.DB().Do(tarantool.NewCallRequest(findRateFunc).Args([]interface{}{gwgrId, dateAt, dir, aNumber, bNumber})).Get()
	if err != nil {
		return 0, 0, time.Since(ts), err
	}
	if len(resp) == 2 {
		return 1, 1, time.Since(ts), nil
	}

	return 0, 0, time.Since(ts), fmt.Errorf("unexpected response length %d", len(resp))
}

func (r *repository) FindSupRates(_ context.Context, gwgrIds []int64, dateAt int64, aNumber, bNumber string) (int64, time.Duration, error) {
	ts := time.Now()
	resp, err := r.db.DB().Do(tarantool.NewCallRequest(findSupRatesFunc).Args([]interface{}{gwgrIds, dateAt, aNumber, bNumber})).Get()
	if err != nil {
		return 0, time.Since(ts), err
	}
	if len(resp) == 0 {
		return 0, time.Since(ts), fmt.Errorf("unexpected response length %d", len(resp))
	}

	return int64(len(resp)), time.Since(ts), nil
}
