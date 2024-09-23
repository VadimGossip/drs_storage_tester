package rate

import (
	"context"
	"strconv"
	"time"
)

func (s *service) FindRate(ctx context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber uint64) (int64, float64, time.Duration, error) {
	return s.rateRepository.FindRate(ctx, gwgrId, dateAt, dir, strconv.Itoa(int(aNumber)), strconv.Itoa(int(bNumber)))
}

func (s *service) FindSupRates(ctx context.Context, gwgrIds []int64, dateAt int64, aNumber, bNumber uint64) (int64, time.Duration, error) {
	return s.rateRepository.FindSupRates(ctx, gwgrIds, dateAt, strconv.Itoa(int(aNumber)), strconv.Itoa(int(bNumber)))
}
