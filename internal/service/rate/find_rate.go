package rate

import (
	"context"
	"strconv"
	"time"
)

func (s *service) FindRate(_ context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber uint64) (int64, float64, time.Duration, error) {
	return s.rateRepository.FindRate(gwgrId, dateAt, dir, strconv.Itoa(int(aNumber)), strconv.Itoa(int(bNumber)))
}

func (s *service) FindSupRates(_ context.Context, gwgrIds []int64, dateAt int64, aNumber, bNumber uint64) (int64, time.Duration, error) {
	return s.rateRepository.FindSupRates(gwgrIds, dateAt, strconv.Itoa(int(aNumber)), strconv.Itoa(int(bNumber)))
}
