package rate

import (
	"context"
	"strconv"
	"time"

	"github.com/VadimGossip/drs_storage_tester/internal/model"
)

func (s *service) FindRate(ctx context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber uint64) (model.RateBase, time.Duration, error) {
	return s.rateRepository.FindRate(ctx, gwgrId, dateAt, dir, strconv.Itoa(int(aNumber)), strconv.Itoa(int(bNumber)))
}
