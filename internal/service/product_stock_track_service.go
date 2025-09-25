package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
	"time"
)

type ProductStockTrackService struct {
	trackRepo *repository.ProductStockTrackRepository
	stockRepo *repository.ProductStockRepository
}

type CreateProductStockTrackRequest struct {
	ProductStockID uint      `json:"product_stock_id" validate:"required"`
	ProductID      *uint     `json:"product_id,omitempty"`
	ProductBatchID *uint     `json:"product_batch_id,omitempty"`
	Quantity       *float64  `json:"quantity" validate:"omitempty,gt=0"`
	Operation      *string   `json:"operation" validate:"omitempty,oneof=Plus Minus"`
	Stock          *float64  `json:"stock" validate:"omitempty,gte=0"`
	Action         string    `json:"action,omitempty"` // CREATE, UPDATE, DELETE
}

type UpdateProductStockTrackRequest struct {
	Quantity  *float64 `json:"quantity,omitempty" validate:"omitempty,gt=0"`
	Operation *string  `json:"operation,omitempty" validate:"omitempty,oneof=Plus Minus"`
	Stock     *float64 `json:"stock,omitempty" validate:"omitempty,gte=0"`
}

func NewProductStockTrackService() *ProductStockTrackService {
	return &ProductStockTrackService{
		trackRepo: repository.NewProductStockTrackRepository(),
		stockRepo: repository.NewProductStockRepository(),
	}
}

// Business logic methods
func (s *ProductStockTrackService) GetAllProductStockTracks() (interface{}, error) {
	return s.trackRepo.GetAllProductStockTracks()
}

func (s *ProductStockTrackService) GetProductStockTracksByStock(stockID uint) (interface{}, error) {
	if stockID == 0 {
		return nil, errors.New("invalid stock ID")
	}

	// Check if stock exists
	_, err := s.stockRepo.GetProductStockModelByID(stockID)
	if err != nil {
		return nil, errors.New("product stock not found")
	}

	return s.trackRepo.GetProductStockTracksByStock(stockID)
}

func (s *ProductStockTrackService) GetProductStockTracksByProduct(productID uint) (interface{}, error) {
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

	return s.trackRepo.GetProductStockTracksByProduct(productID)
}

func (s *ProductStockTrackService) GetProductStockTrackByID(id uint) (interface{}, error) {
	return s.trackRepo.GetProductStockTrackByID(id)
}

func (s *ProductStockTrackService) CreateProductStockTrack(req CreateProductStockTrackRequest, userID uint) (interface{}, error) {
	// Validate required fields
	if req.ProductStockID == 0 {
		return nil, errors.New("product stock ID is required")
	}

	// Check if product stock exists and get details
	stock, err := s.stockRepo.GetProductStockModelByID(req.ProductStockID)
	if err != nil {
		return nil, err
	}

	// Set default values from stock if not provided
	productID := stock.ProductID
	if req.ProductID != nil {
		productID = *req.ProductID
	}

	productBatchID := stock.ProductBatchID
	if req.ProductBatchID != nil {
		productBatchID = *req.ProductBatchID
	}

	// Set default operation if not provided
	operation := "Plus"
	if req.Operation != nil {
		operation = *req.Operation
	}

	// Calculate current stock if not provided
	currentStock := float64(0)
	if stock.Quantity != nil {
		currentStock = *stock.Quantity
	}
	if req.Stock != nil {
		currentStock = *req.Stock
	}

	// Calculate quantity value (default to 0 if not provided)
	quantity := float64(0)
	if req.Quantity != nil {
		quantity = *req.Quantity
	}

	// Create new product stock track
	track := &model.ProductStockTrack{
		ProductStockID: req.ProductStockID,
		ProductID:      productID,
		ProductBatchID: productBatchID,
		Date:           time.Now(),
		Quantity:       quantity,
		Operation:      operation,
		Stock:          currentStock,
		UserIns:        &userID,
		UserUpdt:       &userID,
	}

	err = s.trackRepo.CreateProductStockTrack(track)
	if err != nil {
		return nil, err
	}

	// Return created track
	return s.trackRepo.GetProductStockTrackByID(track.ID)
}

func (s *ProductStockTrackService) UpdateProductStockTrack(id uint, req UpdateProductStockTrackRequest, userID uint) (interface{}, error) {
	// Check if track exists
	_, err := s.trackRepo.GetProductStockTrackModelByID(id)
	if err != nil {
		return nil, err
	}

	// Prepare update data
	updateData := map[string]interface{}{
		"user_updt":  &userID,
		"updated_at": time.Now(),
	}

	// Update provided fields
	if req.Quantity != nil {
		if *req.Quantity <= 0 {
			return nil, errors.New("quantity must be greater than 0")
		}
		updateData["quantity"] = *req.Quantity
	}
	if req.Operation != nil {
		if *req.Operation != "Plus" && *req.Operation != "Minus" {
			return nil, errors.New("operation must be 'Plus' or 'Minus'")
		}
		updateData["operation"] = *req.Operation
	}
	if req.Stock != nil {
		if *req.Stock < 0 {
			return nil, errors.New("stock cannot be negative")
		}
		updateData["stock"] = *req.Stock
	}

	err = s.trackRepo.UpdateProductStockTrack(id, updateData)
	if err != nil {
		return nil, err
	}

	// Return updated track
	return s.trackRepo.GetProductStockTrackByID(id)
}

func (s *ProductStockTrackService) DeleteProductStockTrack(id uint, userID uint) error {
	// Check if track exists
	_, err := s.trackRepo.GetProductStockTrackModelByID(id)
	if err != nil {
		return err
	}

	return s.trackRepo.DeleteProductStockTrackWithAudit(id, userID)
}