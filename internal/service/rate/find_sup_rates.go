package rate

import (
	"context"
	"github.com/VadimGossip/drs_storage_tester/internal/model"
	"strconv"
	"time"
)

func (s *service) FindSupRates(ctx context.Context, gwgrIds []int64, dateAt int64, aNumber, bNumber uint64) (map[int64]model.RateBase, time.Duration, error) {
	return s.rateRepository.FindSupRates(ctx, gwgrIds, dateAt, strconv.Itoa(int(aNumber)), strconv.Itoa(int(bNumber)))
}
