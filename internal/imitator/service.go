package imitator

import (
	"context"
	"fmt"
	"github.com/VadimGossip/tj-drs-storage/internal/data"
	"github.com/VadimGossip/tj-drs-storage/internal/domain"
	"github.com/VadimGossip/tj-drs-storage/internal/event"
	"github.com/VadimGossip/tj-drs-storage/internal/rate"
	"github.com/VadimGossip/tj-drs-storage/pkg/util"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Service interface {
	RunTests(ctx context.Context, task *domain.Task) error
}

type service struct {
	rate rate.Service
	data data.Service
}

var _ Service = (*service)(nil)

func NewService(rate rate.Service, data data.Service) *service {
	return &service{rate: rate, data: data}
}

func (s *service) addDurationToSummary(durSummary *domain.DurationSummary, dur time.Duration, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	if durSummary.Min > dur {
		durSummary.Min = dur
	}
	if durSummary.Max < dur {
		durSummary.Max = dur
	}
	durSummary.EMA.Add(float64(dur.Nanoseconds()))
	durSummary.Histogram[(util.RoundFloat(float64(dur.Milliseconds()/100), 0)*100)+100]++
	return
}

func (s *service) sendFindRateRequest(ctx context.Context, req domain.TaskRequest, summary *domain.TaskSummary, mu *sync.Mutex) error {
	ts := time.Now()
	_, _, dbDur, err := s.rate.FindRate(ctx, req.GwgrId, ts.Unix(), 0, req.Anumber, req.Bnumber)
	if err != nil {
		return err
	}

	s.addDurationToSummary(summary.TotalDuration, time.Since(ts), mu)
	s.addDurationToSummary(summary.DbDuration, dbDur, mu)

	return nil
}
func (s *service) sendFindSupRatesRequest(ctx context.Context, supGwgrIds []int64, req domain.TaskRequest, summary *domain.TaskSummary, mu *sync.Mutex) error {
	ts := time.Now()
	_, dbDur, err := s.rate.FindSupRates(ctx, supGwgrIds, ts.Unix(), req.Anumber, req.Bnumber)
	if err != nil {
		return err
	}

	s.addDurationToSummary(summary.TotalDuration, time.Since(ts), mu)
	s.addDurationToSummary(summary.DbDuration, dbDur, mu)

	return nil
}

func (s *service) RunTests(ctx context.Context, task *domain.Task) error {
	if err := s.data.Refresh(ctx, int64(task.Summary.Total)); err != nil {
		return err
	}
	rateCtrl := event.NewService()
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for e := range rateCtrl.RunEventGeneration(ctx, task.Summary.Total, task.RequestsPerSec, task.PackPerSec) {
		for i := 0; i < e; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				if err := s.sendFindSupRatesRequest(ctx, s.data.GetSupGwgrIds(), s.data.GetTaskRequest(), task.Summary, mu); err != nil {
					logrus.Errorf("sendFindSupRatesRequest err %s", err)
				}
			}(wg)
		}
	}
	wg.Wait()
	fmt.Printf("Total Histogram %+v\n", task.Summary.TotalDuration.Histogram)
	fmt.Printf("Db Histogram %+v\n", task.Summary.TotalDuration.Histogram)
	fmt.Printf("Request EMA Answer TotalDuration %+v(%+v)\n", time.Duration(task.Summary.TotalDuration.EMA.Value()), time.Duration(task.Summary.DbDuration.EMA.Value()))
	fmt.Printf("Request Min Answer Duration %+v(%+v)\n", task.Summary.TotalDuration.Min, task.Summary.DbDuration.Min)
	fmt.Printf("Request Max Answer Duration %+v(%+v)\n", task.Summary.TotalDuration.Max, task.Summary.DbDuration.Max)
	return nil
}
