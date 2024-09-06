package rate

import (
	"github.com/VadimGossip/tj-drs-storage/internal/repository"
	def "github.com/VadimGossip/tj-drs-storage/internal/service"
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
