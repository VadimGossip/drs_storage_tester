package imitator

import (
	"context"
	"github.com/VadimGossip/tj-drs-storage/internal/rate"
	"github.com/sirupsen/logrus"
	"time"
)

type Service interface {
	RunTests(ctx context.Context) error
}

type service struct {
	rate rate.Service
}

var _ Service = (*service)(nil)

func NewService(rate rate.Service) *service {
	return &service{rate: rate}
}

func (s *service) RunTests(ctx context.Context) error {
	ts := time.Now()

	var gwgrId int64 = 4728
	var aNumber uint64 = 525594178906
	var bNumber uint64 = 524423388739
	_, _, err := s.rate.FindRate(ctx, gwgrId, ts.Unix(), 0, aNumber, bNumber)
	if err != nil {
		return err
	}

	logrus.Infof("Test param %+v, duration %s", Task{
		RequestsPerSec: 1,
		PackPerSec:     1,
		Total:          1,
	}, time.Since(ts))

	return nil
}
