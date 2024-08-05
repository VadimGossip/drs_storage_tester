package data

import (
	"context"
	"sync"

	"github.com/VadimGossip/tj-drs-storage/internal/domain"
)

type Source interface {
	GetTaskRequests(ctx context.Context, limit int64) ([]domain.TaskRequest, error)
}

type Service interface {
	GetTaskRequest() domain.TaskRequest
}

type service struct {
	dbSource Source
	requests []domain.TaskRequest
	mu       sync.Mutex
	mark     int
}

var _ Service = (*service)(nil)

func NewService(dbSource Source) *service {
	s := &service{dbSource: dbSource}
	return s
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
