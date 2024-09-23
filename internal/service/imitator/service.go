package imitator

import (
	"context"
	"fmt"
	"github.com/VadimGossip/drs_storage_tester/internal/config"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/VadimGossip/drs_storage_tester/internal/model"
	def "github.com/VadimGossip/drs_storage_tester/internal/service"
	"github.com/VadimGossip/drs_storage_tester/internal/service/event"
	"github.com/VadimGossip/drs_storage_tester/pkg/util"
)

var _ def.ImitatorService = (*service)(nil)

type service struct {
	cfg  config.ImitatorConfig
	rate def.RateService
	data def.DataService
}

func NewService(cfg config.ImitatorConfig,
	rate def.RateService,
	data def.DataService) *service {
	return &service{cfg: cfg,
		rate: rate,
		data: data}
}

func (s *service) addDurationToSummary(durSummary *model.DurationSummary, dur time.Duration, mu *sync.Mutex) {
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

func (s *service) sendFindRateRequest(ctx context.Context, req model.TaskRequest, summary *model.TaskSummary, mu *sync.Mutex) error {
	ts := time.Now()
	_, _, dbDur, err := s.rate.FindRate(ctx, req.GwgrId, ts.Unix(), 0, req.Anumber, req.Bnumber)
	if err != nil {
		return err
	}

	s.addDurationToSummary(summary.TotalDuration, time.Since(ts), mu)
	s.addDurationToSummary(summary.DbDuration, dbDur, mu)

	return nil
}
func (s *service) sendFindSupRatesRequest(ctx context.Context, supGwgrIds []int64, req model.TaskRequest, summary *model.TaskSummary, mu *sync.Mutex) error {
	ts := time.Now()
	_, dbDur, err := s.rate.FindSupRates(ctx, supGwgrIds, ts.Unix(), req.Anumber, req.Bnumber)
	if err != nil {
		return err
	}

	s.addDurationToSummary(summary.TotalDuration, time.Since(ts), mu)
	s.addDurationToSummary(summary.DbDuration, dbDur, mu)

	return nil
}

func (s *service) RunTests(ctx context.Context) error {
	task := model.Task{
		RequestsPerSec: s.cfg.RequestPerSecond(),
		PackPerSec:     s.cfg.PackPerSecond(),
		Summary: &model.TaskSummary{
			Total: s.cfg.TotalRequests(),
			DbDuration: &model.DurationSummary{
				EMA: util.NewEMA(0.01),
			},
			TotalDuration: &model.DurationSummary{
				EMA: util.NewEMA(0.01),
			},
		},
	}

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
				if s.cfg.RequestType() == s.cfg.AllSupplierRequestType() {
					if err := s.sendFindSupRatesRequest(ctx, s.data.GetSupGwgrIds(), s.data.GetTaskRequest(), task.Summary, mu); err != nil {
						logrus.Errorf("sendFindSupRatesRequest err %s", err)
					}
					return
				}
				if err := s.sendFindRateRequest(ctx, s.data.GetTaskRequest(), task.Summary, mu); err != nil {
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
