package rate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	db "github.com/VadimGossip/drs_storage_tester/internal/client/db/keydb"
	"github.com/VadimGossip/drs_storage_tester/internal/model"
	def "github.com/VadimGossip/drs_storage_tester/internal/repository"
	"github.com/VadimGossip/drs_storage_tester/pkg/util"
)

var _ def.RateRepository = (*repository)(nil)

type Repository interface {
	FindRate(ctx context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber string) (int64, float64, time.Duration, error)
	FindSupRates(ctx context.Context, gwgrIds []int64, dateAt int64, aNumber, bNumber string) (int64, time.Duration, error)
}

type repository struct {
	db db.Client
}

var _ Repository = (*repository)(nil)

func NewRepository(db db.Client) *repository {
	return &repository{db: db}
}

func (r *repository) buildRAKeyStr(key model.ARmsgKey) string {
	return model.RAObjectKey + ":" + strconv.FormatInt(key.GwgrId, 10) + ":" + strconv.Itoa(int(key.Direction)) + ":" + strconv.FormatInt(key.BRmsgId, 10) + ":" + key.Code
}

func (r *repository) buildRBKeyStr(key model.BRmsgKey) string {
	return model.RBObjectKey + ":" + strconv.FormatInt(key.GwgrId, 10) + ":" + strconv.Itoa(int(key.Direction)) + ":" + key.Code
}

func (r *repository) getBRmsg(ctx context.Context, key model.BRmsgKey, dateAt int64) (int64, time.Duration, error) {
	keys := make([]string, 0)
	for i := len(key.Code); i > 0; i-- {
		keys = append(keys, r.buildRBKeyStr(key))
		key.Code = key.Code[:i-1]
	}

	values, dur, err := r.db.DB().MGetWithDur(ctx, keys...)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, dur, fmt.Errorf("can't find B-code rate group")
		}
	}

	for i := range values {
		if values[i] != nil {
			result := make([]model.IdHistItem, 0)
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

func (r *repository) getARmsg(ctx context.Context, key model.ARmsgKey, dateAt int64) (int64, time.Duration, error) {
	keys := make([]string, 0)
	for i := len(key.Code); i > 0; i-- {
		keys = append(keys, r.buildRAKeyStr(key))
		key.Code = key.Code[:i-1]
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
			result := make([]model.IdHistItem, 0)
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

func (r *repository) getRateRmsvRmsr(ctx context.Context, key model.RateKey, dateAt int64) (int64, int64, time.Duration, error) {
	result := make([]model.RmsRateHistItem, 0)
	keyStr := model.RTSObjectKey + ":" + strconv.Itoa(int(key.GwgrId)) + ":" + strconv.Itoa(int(key.Direction)) + ":" + strconv.Itoa(int(key.ARmsgId)) + ":" + strconv.Itoa(int(key.BRmsgId))

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

func (r *repository) getRateValue(ctx context.Context, rmsvId int64) (model.Rate, time.Duration, error) {
	var rateVal model.Rate
	keyStr := model.RVObjectKey + ":" + strconv.Itoa(int(rmsvId))

	val, dur, err := r.db.DB().GetWithDur(ctx, keyStr)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return model.Rate{}, dur, fmt.Errorf("can't find rate value")
		}
		return model.Rate{}, dur, err
	}
	err = json.Unmarshal([]byte(val), &rateVal)
	if err != nil {
		return model.Rate{}, dur, err
	}

	return rateVal, dur, nil
}

func (r *repository) getCurrencyRate(ctx context.Context, currencyId int64, dateAt int64) (float64, time.Duration, error) {
	result := make([]model.CurrencyRateHist, 0)
	keyStr := model.CURRTSObjectKey + ":" + strconv.Itoa(int(currencyId))

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

func (r *repository) FindRate(ctx context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber string) (int64, float64, time.Duration, error) {
	var totalDur time.Duration
	bRmsgId, dur, err := r.getBRmsg(ctx, model.BRmsgKey{
		GwgrId:    gwgrId,
		Direction: dir,
		Code:      bNumber,
	}, dateAt)
	if err != nil {
		return 0, 0, dur, err
	}

	aRmsgId, dur, err := r.getARmsg(ctx, model.ARmsgKey{
		GwgrId:    gwgrId,
		Direction: dir,
		BRmsgId:   bRmsgId,
		Code:      aNumber,
	}, dateAt)
	totalDur += dur

	rmsrId, rmsvId, dur, err := r.getRateRmsvRmsr(ctx, model.RateKey{
		GwgrId:    gwgrId,
		Direction: dir,
		ARmsgId:   aRmsgId,
		BRmsgId:   bRmsgId,
	}, dateAt)
	totalDur += dur
	if err != nil {
		return 0, 0, totalDur, err
	}

	rv, dur, err := r.getRateValue(ctx, rmsvId)
	totalDur += dur
	if err != nil {
		return 0, 0, totalDur, err
	}

	currencyRate, dur, err := r.getCurrencyRate(ctx, rv.CurrencyId, dateAt)
	totalDur += dur
	if err != nil {
		return 0, 0, totalDur, err
	}

	return rmsrId, util.RoundFloat(rv.Price*currencyRate, 7), totalDur, nil
}

func (r *repository) FindSupRates(ctx context.Context, gwgrIds []int64, dateAt int64, aNumber, bNumber string) (int64, time.Duration, error) {
	return 0, 0, fmt.Errorf("unimplemented")
}
