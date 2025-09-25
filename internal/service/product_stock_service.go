package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
	"time"
)

type ProductStockService struct {
	stockRepo    *repository.ProductStockRepository
	trackService *ProductStockTrackService
}

func NewProductStockService() *ProductStockService {
	return &ProductStockService{
		stockRepo:    repository.NewProductStockRepository(),
		trackService: NewProductStockTrackService(),
	}
}

// Business logic methods
func (s *ProductStockService) GetAllProductStocks() (interface{}, error) {
	return s.stockRepo.GetAllProductStocks()
}

func (s *ProductStockService) GetProductStocksByProduct(productID uint) (interface{}, error) {
	if productID == 0 {
		return nil, errors.New("invalid product ID")
	}

	// Check if product exists
	productExists, err := s.stockRepo.CheckProductExists(productID)
	if err != nil {
		return nil, err
	}
	if !productExists {
		return nil, errors.New("product not found")
	}

	return s.stockRepo.GetProductStocksByProduct(productID)
}

func (s *ProductStockService) GetProductStockByID(id uint) (interface{}, error) {
	return s.stockRepo.GetProductStockByID(id)
}

func (s *ProductStockService) CreateProductStock(productBatchID, productID, locationID uint, quantity *float64, userID uint) (interface{}, error) {
	// Validate required fields
	if productBatchID == 0 || productID == 0 {
		return nil, errors.New("product batch ID and product ID are required")
	}

	// Check if product exists
	productExists, err := s.stockRepo.CheckProductExists(productID)
	if err != nil {
		return nil, err
	}
	if !productExists {
		return nil, errors.New("product not found")
	}

	// Check if product batch exists
	batchExists, err := s.stockRepo.CheckProductBatchExists(productBatchID)
	if err != nil {
		return nil, err
	}
	if !batchExists {
		return nil, errors.New("product batch not found")
	}

	// Check if location exists (if provided)
	if locationID > 0 {
		// Location check would need to be implemented in repository
		// For now, just skip this validation
	}

	// Prepare quantity (default to 0 if not provided)
	defaultQuantity := float64(0)
	if quantity != nil {
		defaultQuantity = *quantity
	}

	// Create new product stock
	stock := &model.ProductStock{
		ProductBatchID: productBatchID,
		ProductID:      productID,
		LocationID:     locationID,
		Quantity:       &defaultQuantity,
		UserIns:        &userID,
		UserUpdt:       &userID,
	}

	err = s.stockRepo.CreateProductStock(stock)
	if err != nil {
		return nil, err
	}

	// Create tracking record if trackService is available
	if s.trackService != nil && stock.ID > 0 {
		// For now, skip tracking since we need to fix the track service structure first
		// trackReq := CreateProductStockTrackRequest{
		// 	ProductStockID: stock.ID,
		// 	Quantity:       &defaultQuantity,
		// 	Action:         "CREATE",
		// }
		// _, _ = s.trackService.CreateProductStockTrack(trackReq, userID)
	}

	// Return created stock
	return s.stockRepo.GetProductStockByID(stock.ID)
}

func (s *ProductStockService) UpdateProductStock(id, productBatchID, productID, locationID uint, quantity *float64, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid product stock ID")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if stock exists
	_, err := s.stockRepo.GetProductStockModelByID(id)
	if err != nil {
		return nil, err
	}

	// Validate product ID if being updated
	if productID != 0 {
		productExists, err := s.stockRepo.CheckProductExists(productID)
		if err != nil {
			return nil, err
		}
		if !productExists {
			return nil, errors.New("product not found")
		}
	}

	// Validate location ID if being updated
	if locationID > 0 {
		// Location check would need to be implemented in repository
		// For now, just skip this validation
	}

	// Prepare update data
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"updated_at": time.Now(),
	}

	// Update provided fields
	if productBatchID != 0 {
		updateData["product_batch_id"] = productBatchID
	}
	if productID != 0 {
		updateData["product_id"] = productID
	}
	if locationID > 0 {
		updateData["location_id"] = locationID
	}
	if quantity != nil {
		if *quantity < 0 {
			return nil, errors.New("quantity cannot be negative")
		}
		updateData["quantity"] = *quantity
	}

	err = s.stockRepo.UpdateProductStock(id, updateData)
	if err != nil {
		return nil, err
	}

	// Create tracking record if quantity changed (commented for now)
	// if s.trackService != nil && quantity != nil && newQuantity != oldQuantity {
	//   // Add tracking logic here when track service is fixed
	// }

	// Return updated stock
	return s.stockRepo.GetProductStockByID(id)
}

func (s *ProductStockService) DeleteProductStock(id uint, userID uint) error {
	// Check if stock exists
	_, err := s.stockRepo.GetProductStockModelByID(id)
	if err != nil {
		return err
	}

	// Create tracking record before deletion (commented for now)
	// if s.trackService != nil {
	//   // Add tracking logic here when track service is fixed
	// }

	return s.stockRepo.DeleteProductStockWithAudit(id, userID)
}

// Additional business logic methods
func (s *ProductStockService) GetProductStocksByLocation(locationID uint) (interface{}, error) {
	if locationID == 0 {
		return nil, errors.New("invalid location ID")
	}

	// Location check would need to be implemented in repository
	// For now, just skip this validation and return empty or implement basic query
	return []interface{}{}, nil // Placeholder return
}
