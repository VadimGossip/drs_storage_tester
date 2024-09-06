package repository

import (
	"time"
)

type RateRepository interface {
	FindRate(gwgrId, dateAt int64, dir uint8, aNumber, bNumber string) (int64, float64, time.Duration, error)
}
