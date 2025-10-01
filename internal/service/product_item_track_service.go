package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
	"time"
)

type ProductItemTrackService struct {
	trackRepo *repository.ProductItemTrackRepository
	itemRepo  *repository.ProductItemRepository
	stockRepo *repository.ProductStockRepository
}

func NewProductItemTrackService() *ProductItemTrackService {
	return &ProductItemTrackService{
		trackRepo: repository.NewProductItemTrackRepository(),
		itemRepo:  repository.NewProductItemRepository(),
		stockRepo: repository.NewProductStockRepository(),
	}
}

// Business logic methods
func (s *ProductItemTrackService) GetAllProductItemTracks() (interface{}, error) {
	return s.trackRepo.GetAllProductItemTracks()
}

func (s *ProductItemTrackService) GetProductItemTracksByItem(itemID uint) (interface{}, error) {
	if itemID == 0 {
		return nil, errors.New("invalid item ID")
	}

	// Check if item exists
	_, err := s.itemRepo.GetProductItemModelByID(itemID)
	if err != nil {
		return nil, errors.New("product item not found")
	}

	return s.trackRepo.GetProductItemTracksByItem(itemID)
}

func (s *ProductItemTrackService) GetProductItemTracksByStock(stockID uint) (interface{}, error) {
	if stockID == 0 {
		return nil, errors.New("invalid stock ID")
	}

	// Check if stock exists
	_, err := s.stockRepo.GetProductStockModelByID(stockID)
	if err != nil {
		return nil, errors.New("product stock not found")
	}

	return s.trackRepo.GetProductItemTracksByStock(stockID)
}

func (s *ProductItemTrackService) GetProductItemTracksByProduct(productID uint) (interface{}, error) {
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

	return s.trackRepo.GetProductItemTracksByProduct(productID)
}

func (s *ProductItemTrackService) GetProductItemTracksByDateRange(startDate, endDate time.Time) (interface{}, error) {
	if startDate.IsZero() || endDate.IsZero() {
		return nil, errors.New("start date and end date are required")
	}

	if startDate.After(endDate) {
		return nil, errors.New("start date cannot be after end date")
	}

	return s.trackRepo.GetProductItemTracksByDateRange(startDate, endDate)
}

func (s *ProductItemTrackService) GetProductItemTrackByID(id uint) (interface{}, error) {
	return s.trackRepo.GetProductItemTrackByID(id)
}

func (s *ProductItemTrackService) CreateProductItemTrack(productItemID uint, productStockID, productID, productBatchID *uint, date time.Time, unitPrice *string, stockIn, stockOut, quantity *float64, operation *string, stock *float64, description *string, action string, userID uint) (interface{}, error) {
	// Validate required fields
	if productItemID == 0 {
		return nil, errors.New("product item ID is required")
	}

	if date.IsZero() {
		return nil, errors.New("date is required")
	}

	// Check if product item exists and get details
	item, err := s.itemRepo.GetProductItemModelByID(productItemID)
	if err != nil {
		return nil, err
	}

	// Set default values from item if not provided
	defaultProductStockID := item.ProductStockID
	if productStockID != nil {
		defaultProductStockID = *productStockID
	}

	defaultProductID := item.ProductID
	if productID != nil {
		defaultProductID = *productID
	}

	defaultProductBatchID := item.ProductBatchID
	if productBatchID != nil {
		defaultProductBatchID = *productBatchID
	}

	// Determine operation based on stock in/out or provided operation
	defaultOperation := "In"
	if operation != nil {
		defaultOperation = *operation
	} else if stockOut != nil && *stockOut > 0 {
		defaultOperation = "Out"
	} else if stockIn != nil && *stockIn > 0 {
		defaultOperation = "In"
	}

	// Calculate current stock if not provided
	currentStock := float64(0)
	if item.Quantity != nil {
		currentStock = *item.Quantity
	}
	if stock != nil {
		currentStock = *stock
	}

	// Calculate quantity value (default to 0 if not provided)
	defaultQuantity := float64(0)
	if quantity != nil {
		defaultQuantity = *quantity
	}

	// Create new product item track
	track := &model.ProductItemTrack{
		ProductStockID: defaultProductStockID,
		ProductBatchID: defaultProductBatchID,
		ProductID:      defaultProductID,
		Date:           date,
		Quantity:       defaultQuantity,
		Operation:      defaultOperation,
		Stock:          currentStock,
		UnitPrice:      unitPrice,
		Description:    description,
		UserIns:        &userID,
		UserUpdt:       &userID,
	}

	err = s.trackRepo.CreateProductItemTrack(track)
	if err != nil {
		return nil, err
	}

	// Return created track
	return s.trackRepo.GetProductItemTrackByID(track.ID)
}

func (s *ProductItemTrackService) UpdateProductItemTrack(id uint, date *time.Time, unitPrice *string, quantity *float64, operation *string, stock *float64, description *string, userID uint) (interface{}, error) {
	// Check if track exists
	_, err := s.trackRepo.GetProductItemTrackModelByID(id)
	if err != nil {
		return nil, err
	}

	// Prepare update data
	updateData := map[string]interface{}{
		"user_updt":  &userID,
		"updated_at": time.Now(),
	}

	// Update provided fields
	if date != nil && !date.IsZero() {
		updateData["date"] = *date
	}
	if unitPrice != nil {
		updateData["unit_price"] = *unitPrice
	}
	if quantity != nil {
		if *quantity <= 0 {
			return nil, errors.New("quantity must be greater than 0")
		}
		updateData["quantity"] = *quantity
	}
	if operation != nil {
		if *operation != "In" && *operation != "Out" && *operation != "Plus" && *operation != "Minus" {
			return nil, errors.New("operation must be 'In', 'Out', 'Plus', or 'Minus'")
		}
		updateData["operation"] = *operation
	}
	if stock != nil {
		if *stock < 0 {
			return nil, errors.New("stock cannot be negative")
		}
		updateData["stock"] = *stock
	}
	if description != nil {
		updateData["description"] = *description
	}

	err = s.trackRepo.UpdateProductItemTrack(id, updateData)
	if err != nil {
		return nil, err
	}

	// Return updated track
	return s.trackRepo.GetProductItemTrackByID(id)
}

func (s *ProductItemTrackService) DeleteProductItemTrack(id uint, userID uint) error {
	// Check if track exists
	_, err := s.trackRepo.GetProductItemTrackModelByID(id)
	if err != nil {
		return err
	}

	return s.trackRepo.DeleteProductItemTrackWithAudit(id, userID)
}
