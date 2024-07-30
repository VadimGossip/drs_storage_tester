package rate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type Repository interface {
	GetBRmsg(ctx context.Context, key bRmsgKey, dateAt int64) (int64, error)
	GetARmsg(ctx context.Context, key aRmsgKey, dateAt int64) (int64, error)
	GetRateRmsvRmsr(ctx context.Context, key rateKey, dateAt int64) (int64, int64, error)
	GetRateValue(ctx context.Context, rmsvId int64) (Rate, error)
	GetCurrencyRate(ctx context.Context, currencyId int64, dateAt int64) (float64, error)
}

type repository struct {
	kdb *redis.Client
}

var _ Repository = (*repository)(nil)

func NewRepository(kdb *redis.Client) *repository {
	return &repository{kdb: kdb}
}

func (r *repository) MSet(ctx context.Context, pairs []interface{}) error {
	var maxSize = 1000000
	for len(pairs) > 0 {
		if maxSize > len(pairs) {
			maxSize = len(pairs)
		}
		if err := r.kdb.MSet(ctx, pairs[:maxSize]...).Err(); err != nil {
			return err
		}
		pairs = pairs[maxSize:]
	}
	return nil
}

func (r *repository) buildRAKeyStr(key aRmsgKey) string {
	return RAObjectKey + ":" + strconv.FormatInt(key.gwgrId, 10) + ":" + strconv.Itoa(int(key.direction)) + ":" + strconv.FormatInt(key.bRmsgId, 10) + ":" + strconv.FormatUint(key.code, 10)
}

func (r *repository) buildRBKeyStr(key bRmsgKey) string {
	return RBObjectKey + ":" + strconv.FormatInt(key.gwgrId, 10) + ":" + strconv.Itoa(int(key.direction)) + ":" + strconv.FormatUint(key.code, 10)
}

func (r *repository) LoadRateAGroups(ctx context.Context, data map[aRmsgKey][]IdHistItem) (time.Duration, error) {
	ts := time.Now()
	resultSlice := make([]interface{}, 0, len(data)*2)
	for key, value := range data {
		v, err := json.Marshal(value)
		if err != nil {
			return time.Since(ts), err
		}

		resultSlice = append(resultSlice, r.buildRAKeyStr(key), v)
	}
	if err := r.MSet(ctx, resultSlice); err != nil {
		return time.Since(ts), err
	}

	return time.Since(ts), nil
}

func (r *repository) LoadRateBGroups(ctx context.Context, data map[bRmsgKey][]IdHistItem) (time.Duration, error) {
	ts := time.Now()
	resultSlice := make([]interface{}, 0, len(data)*2)
	for key, value := range data {
		v, err := json.Marshal(value)
		if err != nil {
			return time.Since(ts), err
		}
		resultSlice = append(resultSlice, r.buildRBKeyStr(key), v)
	}

	if err := r.MSet(ctx, resultSlice); err != nil {
		return time.Since(ts), err
	}

	return time.Since(ts), nil
}

func (r *repository) LoadRates(ctx context.Context, data map[rateKey][]RmsRateHistItem) (time.Duration, error) {
	ts := time.Now()
	resultSlice := make([]interface{}, 0, len(data)*2)
	for key, value := range data {
		fullKey := RTSObjectKey + ":" + strconv.FormatInt(key.gwgrId, 10) + ":" + strconv.Itoa(int(key.direction)) + ":" + strconv.FormatInt(key.aRmsgId, 10) + ":" + strconv.FormatInt(key.bRmsgId, 10)
		v, err := json.Marshal(value)
		if err != nil {
			return time.Since(ts), err
		}

		resultSlice = append(resultSlice, fullKey, v)
	}
	if err := r.MSet(ctx, resultSlice); err != nil {
		return time.Since(ts), err
	}

	return time.Since(ts), nil
}

func (r *repository) LoadRateValues(ctx context.Context, data map[int64]Rate) (time.Duration, error) {
	ts := time.Now()
	resultSlice := make([]interface{}, 0, len(data)*2)
	for key, value := range data {
		fullKey := RVObjectKey + ":" + strconv.FormatInt(key, 10)
		v, err := json.Marshal(value)
		if err != nil {
			return time.Since(ts), err
		}

		resultSlice = append(resultSlice, fullKey, v)
	}
	if err := r.MSet(ctx, resultSlice); err != nil {
		return time.Since(ts), err
	}

	return time.Since(ts), nil
}

func (r *repository) LoadCurrencyRates(ctx context.Context, data map[int64][]CurrencyRateHist) (time.Duration, error) {
	ts := time.Now()
	resultSlice := make([]interface{}, 0, len(data)*2)
	for key, value := range data {
		fullKey := CURRTSObjectKey + ":" + strconv.FormatInt(key, 10)
		v, err := json.Marshal(value)
		if err != nil {
			return time.Since(ts), err
		}

		resultSlice = append(resultSlice, fullKey, v)
	}
	if err := r.MSet(ctx, resultSlice); err != nil {
		return time.Since(ts), err
	}

	return time.Since(ts), nil
}

func (r *repository) GetBRmsg(ctx context.Context, key bRmsgKey, dateAt int64) (int64, error) {
	keys := make([]string, 0)
	for key.code > 0 {
		keys = append(keys, r.buildRBKeyStr(key))
		key.code /= 10
	}

	values, err := r.kdb.MGet(ctx, keys...).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, fmt.Errorf("can't find B-code rate group")
		}
	}

	for i := range values {
		if values[i] != nil {
			result := make([]IdHistItem, 0)
			err = json.Unmarshal([]byte(values[i].(string)), &result)
			if err != nil {
				return 0, err
			}
			for _, value := range result {
				if dateAt >= value.DBegin && dateAt < value.DEnd {
					return value.Id, nil
				}
			}
		}
	}
	return 0, fmt.Errorf("cannot find B-code rate group")
}

func (r *repository) GetARmsg(ctx context.Context, key aRmsgKey, dateAt int64) (int64, error) {
	keys := make([]string, 0)
	for key.code > 0 {
		keys = append(keys, r.buildRAKeyStr(key))
		key.code /= 10
	}

	values, err := r.kdb.MGet(ctx, keys...).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return -2, nil
		}
		return 0, err
	}

	for i := range values {
		if values[i] != nil {
			result := make([]IdHistItem, 0)
			err = json.Unmarshal([]byte(values[i].(string)), &result)
			if err != nil {
				return 0, err
			}
			for _, value := range result {
				if dateAt >= value.DBegin && dateAt < value.DEnd {
					return value.Id, nil
				}
			}
		}
	}
	return -2, nil
}

func (r *repository) GetRateRmsvRmsr(ctx context.Context, key rateKey, dateAt int64) (int64, int64, error) {
	result := make([]RmsRateHistItem, 0)
	keyStr := RTSObjectKey + ":" + strconv.Itoa(int(key.gwgrId)) + ":" + strconv.Itoa(int(key.direction)) + ":" + strconv.Itoa(int(key.aRmsgId)) + ":" + strconv.Itoa(int(key.bRmsgId))

	val, err := r.kdb.Get(ctx, keyStr).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, 0, fmt.Errorf("can't find rate")
		}
	}

	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		return 0, 0, err
	}
	for _, value := range result {
		if dateAt >= value.DBegin && dateAt < value.DEnd {
			return value.RmsrId, value.RmsvId, nil
		}
	}

	return 0, 0, fmt.Errorf("cannot find rate")
}

func (r *repository) GetRateValue(ctx context.Context, rmsvId int64) (Rate, error) {
	var rateVal Rate
	keyStr := RVObjectKey + ":" + strconv.Itoa(int(rmsvId))

	val, err := r.kdb.Get(ctx, keyStr).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return Rate{}, fmt.Errorf("can't find rate value")
		}
	}
	err = json.Unmarshal([]byte(val), &rateVal)
	if err != nil {
		return Rate{}, err
	}

	return rateVal, nil
}

func (r *repository) GetCurrencyRate(ctx context.Context, currencyId int64, dateAt int64) (float64, error) {
	result := make([]CurrencyRateHist, 0)
	keyStr := CURRTSObjectKey + ":" + strconv.Itoa(int(currencyId))

	val, err := r.kdb.Get(ctx, keyStr).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, fmt.Errorf("can't find currency rate")
		}
	}

	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		return 0, err
	}

	for _, value := range result {
		if dateAt >= value.DBegin && dateAt < value.DEnd {
			return value.CurrencyRate, nil
		}
	}

	return 0, fmt.Errorf("can't find currency rate")
}
