package event

import (
	"context"
	"time"

	def "github.com/VadimGossip/drs_storage_tester/internal/service"
)

var _ def.EventService = (*service)(nil)

type genEventsDetails struct {
	packSize     int
	lastPackSize int
	packMark     int
	pps          int
	total        int
	tickerStep   time.Duration
}

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) buildGenTaskDetails(rps, pps int) *genEventsDetails {
	t := genEventsDetails{
		packSize:   rps / pps,
		packMark:   1,
		pps:        1,
		tickerStep: 1000 * time.Millisecond,
	}
	t.lastPackSize = rps - t.packSize
	if t.packSize > 0 {
		t.tickerStep = time.Duration(1000/pps) * time.Millisecond
		t.pps = pps
	}

	return &t
}

func (s *service) RunEventGeneration(ctx context.Context, total, rps, pps int) chan int {
	events := make(chan int)
	t := s.buildGenTaskDetails(rps, pps)
	t.total = total
	go func() {
		defer close(events)
		ticker := time.NewTicker(t.tickerStep)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				n := t.packSize
				if t.packMark == pps {
					n = t.lastPackSize
				}
				events <- n
				t.total -= n
				if t.total == 0 {
					return
				}
			}
		}
	}()
	return events
}
