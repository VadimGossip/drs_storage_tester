package rate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"

	db "github.com/VadimGossip/tj-drs-storage/internal/client/db/keydb"
)

type Repository interface {
	GetBRmsg(ctx context.Context, key bRmsgKey, dateAt int64) (int64, time.Duration, error)
	GetARmsg(ctx context.Context, key aRmsgKey, dateAt int64) (int64, time.Duration, error)
	GetRateRmsvRmsr(ctx context.Context, key rateKey, dateAt int64) (int64, int64, time.Duration, error)
	GetRateValue(ctx context.Context, rmsvId int64) (Rate, time.Duration, error)
	GetCurrencyRate(ctx context.Context, currencyId int64, dateAt int64) (float64, time.Duration, error)
}

type repository struct {
	db db.Client
}

var _ Repository = (*repository)(nil)

func NewRepository(db db.Client) *repository {
	return &repository{db: db}
}

func (r *repository) buildRAKeyStr(key aRmsgKey) string {
	return RAObjectKey + ":" + strconv.FormatInt(key.gwgrId, 10) + ":" + strconv.Itoa(int(key.direction)) + ":" + strconv.FormatInt(key.bRmsgId, 10) + ":" + strconv.FormatUint(key.code, 10)
}

func (r *repository) buildRBKeyStr(key bRmsgKey) string {
	return RBObjectKey + ":" + strconv.FormatInt(key.gwgrId, 10) + ":" + strconv.Itoa(int(key.direction)) + ":" + strconv.FormatUint(key.code, 10)
}

func (r *repository) GetBRmsg(ctx context.Context, key bRmsgKey, dateAt int64) (int64, time.Duration, error) {
	keys := make([]string, 0)
	for key.code > 0 {
		keys = append(keys, r.buildRBKeyStr(key))
		key.code /= 10
	}
	values, dur, err := r.db.DB().MGetWithDur(ctx, keys...)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, dur, fmt.Errorf("can't find B-code rate group")
		}
	}

	for i := range values {
		if values[i] != nil {
			result := make([]IdHistItem, 0)
			err = json.Unmarshal([]byte(values[i].(string)), &result)
			if err != nil {
				return 0, dur, err
			}
			for _, value := range result {
				if dateAt >= value.DBegin && dateAt < value.DEnd {
					return value.Id, dur, nil
				}
			}
		}
	}
	return 0, dur, fmt.Errorf("cannot find B-code rate group")
}

func (r *repository) GetARmsg(ctx context.Context, key aRmsgKey, dateAt int64) (int64, time.Duration, error) {
	keys := make([]string, 0)
	for key.code > 0 {
		keys = append(keys, r.buildRAKeyStr(key))
		key.code /= 10
	}

	values, dur, err := r.db.DB().MGetWithDur(ctx, keys...)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return -2, dur, nil
		}
		return 0, dur, err
	}

	for i := range values {
		if values[i] != nil {
			result := make([]IdHistItem, 0)
			err = json.Unmarshal([]byte(values[i].(string)), &result)
			if err != nil {
				return 0, dur, err
			}
			for _, value := range result {
				if dateAt >= value.DBegin && dateAt < value.DEnd {
					return value.Id, dur, nil
				}
			}
		}
	}
	return -2, dur, nil
}

func (r *repository) GetRateRmsvRmsr(ctx context.Context, key rateKey, dateAt int64) (int64, int64, time.Duration, error) {
	result := make([]RmsRateHistItem, 0)
	keyStr := RTSObjectKey + ":" + strconv.Itoa(int(key.gwgrId)) + ":" + strconv.Itoa(int(key.direction)) + ":" + strconv.Itoa(int(key.aRmsgId)) + ":" + strconv.Itoa(int(key.bRmsgId))

	val, dur, err := r.db.DB().GetWithDur(ctx, keyStr)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, 0, dur, fmt.Errorf("can't find rate")
		}
		return 0, 0, dur, err
	}

	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		return 0, 0, dur, err
	}
	for _, value := range result {
		if dateAt >= value.DBegin && dateAt < value.DEnd {
			return value.RmsrId, value.RmsvId, dur, nil
		}
	}

	return 0, 0, dur, fmt.Errorf("cannot find rate")
}

func (r *repository) GetRateValue(ctx context.Context, rmsvId int64) (Rate, time.Duration, error) {
	var rateVal Rate
	keyStr := RVObjectKey + ":" + strconv.Itoa(int(rmsvId))

	val, dur, err := r.db.DB().GetWithDur(ctx, keyStr)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return Rate{}, dur, fmt.Errorf("can't find rate value")
		}
		return Rate{}, dur, err
	}
	err = json.Unmarshal([]byte(val), &rateVal)
	if err != nil {
		return Rate{}, dur, err
	}

	return rateVal, dur, nil
}

func (r *repository) GetCurrencyRate(ctx context.Context, currencyId int64, dateAt int64) (float64, time.Duration, error) {
	result := make([]CurrencyRateHist, 0)
	keyStr := CURRTSObjectKey + ":" + strconv.Itoa(int(currencyId))

	val, dur, err := r.db.DB().GetWithDur(ctx, keyStr)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, dur, fmt.Errorf("can't find currency rate")
		}
		return 0, dur, err
	}

	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		return 0, dur, err
	}

	for _, value := range result {
		if dateAt >= value.DBegin && dateAt < value.DEnd {
			return value.CurrencyRate, dur, nil
		}
	}

	return 0, dur, fmt.Errorf("can't find currency rate")
}
