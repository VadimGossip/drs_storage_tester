package tarantool

import (
	"context"
	"fmt"
	"github.com/VadimGossip/drs_storage_tester/internal/model"
	"time"

	db "github.com/VadimGossip/platform_common/pkg/db/tarantool"
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
func (r *repository) FindRate(_ context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber string) (model.RateBase, time.Duration, error) {
	ts := time.Now()
	resp, err := r.db.DB().Do(tarantool.NewCallRequest(findRateFunc).Args([]interface{}{gwgrId, dateAt, dir, aNumber, bNumber})).Get()
	if err != nil {
		return model.RateBase{}, time.Since(ts), err
	}
	if len(resp) != 0 {
		return model.RateBase{RmsrId: 1, PriceBase: 1.11}, time.Since(ts), nil
	}

	return model.RateBase{}, time.Since(ts), fmt.Errorf("unexpected response length %d", len(resp))
}

func (r *repository) FindSupRates(_ context.Context, gwgrIds []int64, dateAt int64, aNumber, bNumber string) (map[int64]model.RateBase, time.Duration, error) {
	ts := time.Now()
	resp, err := r.db.DB().Do(tarantool.NewCallRequest(findSupRatesFunc).Args([]interface{}{gwgrIds, dateAt, aNumber, bNumber})).Get()
	if err != nil {
		return nil, time.Since(ts), err
	}
	if len(resp) == 0 {
		return nil, time.Since(ts), fmt.Errorf("unexpected response length %d", len(resp))
	}
	result := map[int64]model.RateBase{
		1: {RmsrId: 1, PriceBase: 1.11},
	}

	return result, time.Since(ts), nil
}
