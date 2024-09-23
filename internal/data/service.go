package data

import (
	"context"
	"sync"

	"github.com/VadimGossip/drs_storage_tester/internal/domain"
)

type Source interface {
	GetTaskRequests(ctx context.Context, limit int64) ([]domain.TaskRequest, error)
	GetSupGwgrIds(ctx context.Context) ([]int64, error)
}

type Service interface {
	GetTaskRequest() domain.TaskRequest
	GetSupGwgrIds() []int64
	Refresh(ctx context.Context, limit int64) error
}

type service struct {
	dbSource   Source
	supGwgrIds []int64
	requests   []domain.TaskRequest
	mu         sync.Mutex
	mark       int
}

var _ Service = (*service)(nil)

func NewService(dbSource Source) *service {
	s := &service{dbSource: dbSource}
	return s
}

func (s *service) Refresh(ctx context.Context, limit int64) error {
	var err error
	s.requests, err = s.dbSource.GetTaskRequests(ctx, limit)
	if err != nil {
		return err
	}

	s.supGwgrIds, err = s.dbSource.GetSupGwgrIds(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetTaskRequest() domain.TaskRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.mark+1 < len(s.requests) {
		s.mark++
	} else {
		s.mark = 0
	}
	return s.requests[s.mark]
}

func (s *service) GetSupGwgrIds() []int64 {
	return s.supGwgrIds
}
