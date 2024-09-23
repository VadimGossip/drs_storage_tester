package data

import (
	"context"
	"sync"

	"github.com/VadimGossip/drs_storage_tester/internal/model"
	def "github.com/VadimGossip/drs_storage_tester/internal/service"
)

var _ def.DataService = (*service)(nil)

type service struct {
	requestService def.RequestService
	supGwgrIds     []int64
	requests       []model.TaskRequest
	mu             sync.Mutex
	mark           int
}

func NewService(requestService def.RequestService) *service {
	s := &service{requestService: requestService}
	return s
}

func (s *service) Refresh(ctx context.Context, limit int64) error {
	var err error
	s.requests, err = s.requestService.GetTaskRequests(ctx, limit)
	if err != nil {
		return err
	}

	s.supGwgrIds, err = s.requestService.GetSupGwgrIds(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetTaskRequest() model.TaskRequest {
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
