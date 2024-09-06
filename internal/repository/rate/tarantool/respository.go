package tarantool

import (
	"fmt"
	"time"

	db "github.com/VadimGossip/tj-drs-storage/internal/client/db/tarantool"

	def "github.com/VadimGossip/tj-drs-storage/internal/repository"

	"github.com/tarantool/go-tarantool/v2"
)

var _ def.RateRepository = (*repository)(nil)

const (
	findRateFunc string = "rates.find_rate"
)

type repository struct {
	db db.Client
}

func NewRepository(db db.Client) *repository {
	return &repository{
		db: db,
	}
}
func (r *repository) FindRate(gwgrId, dateAt int64, dir uint8, aNumber, bNumber string) (int64, float64, time.Duration, error) {
	ts := time.Now()
	resp, err := r.db.DB().Do(tarantool.NewCallRequest("rates.find_rate").Args([]interface{}{gwgrId, dateAt, dir, aNumber, bNumber})).Get()
	if err != nil {
		return 0, 0, time.Since(ts), err
	}
	if len(resp) == 2 {
		return 1, 1, time.Since(ts), nil
	}

	return 0, 0, time.Since(ts), fmt.Errorf("unexpected response length %d", len(resp))
}
