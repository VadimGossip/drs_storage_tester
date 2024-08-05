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
	return &service{rate: rate}
}

func (s *service) addDurationToSummary(summary *domain.TaskSummary, dur time.Duration, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	if summary.Duration.Min > dur {
		summary.Duration.Min = dur
	}
	if summary.Duration.Max < dur {
		summary.Duration.Max = dur
	}
	summary.Duration.EMA.Add(float64(dur.Nanoseconds()))
	summary.Duration.Histogram[(util.RoundFloat(float64(dur.Milliseconds()/100), 0)*100)+100]++
	return
}

func (s *service) sendDbRequest(ctx context.Context, summary *domain.TaskSummary, mu *sync.Mutex) error {
	req := s.data.GetTaskRequest()
	ts := time.Now()

	_, _, err := s.rate.FindRate(ctx, req.GwgrId, ts.Unix(), 0, req.Anumber, req.Bnumber)
	if err != nil {
		return err
	}

	s.addDurationToSummary(summary, time.Since(ts), mu)
	return nil
}

func (s *service) RunTests(ctx context.Context, task *domain.Task) error {
	rateCtrl := event.NewService()
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for e := range rateCtrl.RunEventGeneration(ctx, task.Summary.Total, task.RequestsPerSec, task.PackPerSec) {
		for i := 0; i < e; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				if err := s.sendDbRequest(ctx, task.Summary, mu); err != nil {
					logrus.Errorf("sendDbRequest err %s", err)
				}
			}(wg)
		}
	}
	wg.Wait()
	fmt.Printf("Histogram %+v\n", task.Summary.Duration.Histogram)
	fmt.Printf("Request EMA Answer Duration %+v\n", time.Duration(task.Summary.Duration.EMA.Value()))
	fmt.Printf("Request Min Answer Duration %+v\n", task.Summary.Duration.Min)
	fmt.Printf("Request Max Answer Duration %+v\n", task.Summary.Duration.Max)
	return nil
}
