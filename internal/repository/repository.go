package repository

import (
	"time"
)

type RateRepository interface {
	FindRate(gwgrId, dateAt int64, dir uint8, aNumber, bNumber string) (int64, float64, time.Duration, error)
	FindSupRates(gwgrIds []int64, dateAt int64, aNumber, bNumber string) (int64, time.Duration, error)
}
