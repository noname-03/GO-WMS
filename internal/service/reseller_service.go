package service

import (
	"myapp/internal/repository"
)

type ResellerService struct {
	resellerRepo *repository.ResellerRepository
}

func NewResellerService() *ResellerService {
	return &ResellerService{
		resellerRepo: repository.NewResellerRepository(),
	}
}

// Business logic methods
func (s *ResellerService) GetAllResellers() (interface{}, error) {
	return s.resellerRepo.GetAllResellers()
}
