package rate

import (
	"context"
	"github.com/VadimGossip/tj-drs-storage/pkg/util"
)

type Service interface {
	FindRate(ctx context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber uint64) (int64, float64, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) FindRate(ctx context.Context, gwgrId, dateAt int64, dir uint8, aNumber, bNumber uint64) (int64, float64, error) {
	bRmsgId, err := s.repo.GetBRmsg(ctx, bRmsgKey{
		gwgrId:    gwgrId,
		direction: dir,
		code:      bNumber,
	}, dateAt)
	if err != nil {
		return 0, 0, err
	}
	aRmsgId, err := s.repo.GetARmsg(ctx, aRmsgKey{
		gwgrId:    gwgrId,
		direction: dir,
		bRmsgId:   bRmsgId,
		code:      aNumber,
	}, dateAt)

	rmsrId, rmsvId, err := s.repo.GetRateRmsvRmsr(ctx, rateKey{
		gwgrId:    gwgrId,
		direction: dir,
		aRmsgId:   aRmsgId,
		bRmsgId:   bRmsgId,
	}, dateAt)
	if err != nil {
		return 0, 0, err
	}

	rv, err := s.repo.GetRateValue(ctx, rmsvId)
	if err != nil {
		return 0, 0, err
	}

	currencyRate, err := s.repo.GetCurrencyRate(ctx, rv.CurrencyId, dateAt)
	if err != nil {
		return 0, 0, err
	}

	return rmsrId, util.RoundFloat(rv.Price*currencyRate, 7), nil
}
