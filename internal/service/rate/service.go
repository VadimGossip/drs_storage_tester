package rate

import (
	"github.com/VadimGossip/drs_storage_tester/internal/repository"
	def "github.com/VadimGossip/drs_storage_tester/internal/service"
)

var _ def.RateService = (*service)(nil)

type service struct {
	rateRepository repository.RateRepository
}

func NewService(rateRepository repository.RateRepository) *service {
	return &service{
		rateRepository: rateRepository,
	}
}
