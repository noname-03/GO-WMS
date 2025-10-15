package service

import (
	"myapp/internal/repository"
)

type InventoryService struct {
	inventoryRepo *repository.InventoryRepository
}

func NewInventoryService() *InventoryService {
	return &InventoryService{
		inventoryRepo: repository.NewInventoryRepository(),
	}
}

// GetInventoryStock returns inventory stock information with optional filters
func (s *InventoryService) GetInventoryStock(brandID, categoryID, productID, productBatchID, locationID *uint, barcode *string) (interface{}, error) {
	return s.inventoryRepo.GetInventoryStock(brandID, categoryID, productID, productBatchID, locationID, barcode)
}
