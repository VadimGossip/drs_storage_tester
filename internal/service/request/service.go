package request

import (
	"context"
	"github.com/VadimGossip/drs_storage_tester/internal/repository"

	"github.com/VadimGossip/drs_storage_tester/internal/model"
	def "github.com/VadimGossip/drs_storage_tester/internal/service"
)

var _ def.RequestService = (*service)(nil)

type service struct {
	requestRepo repository.RequestRepository
}

func NewService(requestRepo repository.RequestRepository) *service {
	return &service{requestRepo: requestRepo}
}

func (s *service) GetTaskRequests(ctx context.Context, limit int64) ([]model.TaskRequest, error) {
	return s.requestRepo.GetTaskRequests(ctx, limit)
}

func (s *service) GetSupGwgrIds(ctx context.Context) ([]int64, error) {
	return s.requestRepo.GetSupGwgrIds(ctx)
}
