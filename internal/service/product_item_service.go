package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
	"time"
)

type ProductItemService struct {
	itemRepo     *repository.ProductItemRepository
	stockRepo    *repository.ProductStockRepository
	trackService *ProductItemTrackService
}

func NewProductItemService() *ProductItemService {
	return &ProductItemService{
		itemRepo:     repository.NewProductItemRepository(),
		stockRepo:    repository.NewProductStockRepository(),
		trackService: NewProductItemTrackService(),
	}
}

// Business logic methods
func (s *ProductItemService) GetAllProductItems() (interface{}, error) {
	return s.itemRepo.GetAllProductItems()
}

func (s *ProductItemService) GetProductItemsByStock(stockID uint) (interface{}, error) {
	if stockID == 0 {
		return nil, errors.New("invalid stock ID")
	}

	// Check if stock exists
	_, err := s.stockRepo.GetProductStockModelByID(stockID)
	if err != nil {
		return nil, errors.New("product stock not found")
	}

	return s.itemRepo.GetProductItemsByStock(stockID)
}

func (s *ProductItemService) GetProductItemsByProduct(productID uint) (interface{}, error) {
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

	return s.itemRepo.GetProductItemsByProduct(productID)
}

func (s *ProductItemService) GetProductItemsByLocation(locationID uint) (interface{}, error) {
	if locationID == 0 {
		return nil, errors.New("invalid location ID")
	}

	// Location functionality would need specific repository method
	// For now, return placeholder
	return []interface{}{}, nil
}

func (s *ProductItemService) GetProductItemByID(id uint) (interface{}, error) {
	return s.itemRepo.GetProductItemByID(id)
}

func (s *ProductItemService) CreateProductItem(productStockID, productID, productBatchID uint, locationID *uint, stockIn, stockOut, quantity *float64, userID uint) (interface{}, error) {
	// Validate required fields
	if productStockID == 0 || productID == 0 || productBatchID == 0 {
		return nil, errors.New("product stock ID, product ID, and product batch ID are required")
	}

	// Check if product stock exists
	_, err := s.stockRepo.GetProductStockModelByID(productStockID)
	if err != nil {
		return nil, errors.New("product stock not found")
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
	if locationID != nil && *locationID > 0 {
		// Location check would be implemented when needed
	}

	// Validate stock in/out values
	if stockIn != nil && stockOut != nil {
		if *stockIn > 0 && *stockOut > 0 {
			return nil, errors.New("cannot have both stock in and stock out values")
		}
	}

	// Calculate quantity if not provided
	defaultQuantity := float64(0)
	if quantity != nil {
		defaultQuantity = *quantity
	} else if stockIn != nil {
		defaultQuantity = *stockIn
	} else if stockOut != nil {
		defaultQuantity = -*stockOut // Negative for stock out
	}

	// Create new product item
	item := &model.ProductItem{
		ProductStockID: productStockID,
		ProductID:      productID,
		ProductBatchID: productBatchID,
		StockIn:        stockIn,
		StockOut:       stockOut,
		Quantity:       &defaultQuantity,
		UserIns:        &userID,
		UserUpdt:       &userID,
	}

	err = s.itemRepo.CreateProductItem(item)
	if err != nil {
		return nil, err
	}

	// Create tracking record if trackService is available (commented for now)
	// if s.trackService != nil && item.ID > 0 {
	//   // Add tracking logic here when track service is fixed
	// }

	// Return created item
	return s.itemRepo.GetProductItemByID(item.ID)
}

func (s *ProductItemService) UpdateProductItem(id uint, productStockID, productID *uint, locationID *uint, stockIn, stockOut, quantity *float64, userID uint) (interface{}, error) {
	// Check if item exists
	_, err := s.itemRepo.GetProductItemModelByID(id)
	if err != nil {
		return nil, err
	}

	// Validate product stock ID if being updated
	if productStockID != nil && *productStockID > 0 {
		_, err := s.stockRepo.GetProductStockModelByID(*productStockID)
		if err != nil {
			return nil, errors.New("product stock not found")
		}
	}

	// Validate product ID if being updated
	if productID != nil && *productID > 0 {
		productExists, err := s.stockRepo.CheckProductExists(*productID)
		if err != nil {
			return nil, err
		}
		if !productExists {
			return nil, errors.New("product not found")
		}
	}

	// Validate location ID if being updated
	if locationID != nil && *locationID > 0 {
		// Location validation would be implemented when needed
	}

	// Validate stock in/out values
	if stockIn != nil && stockOut != nil {
		if *stockIn > 0 && *stockOut > 0 {
			return nil, errors.New("cannot have both stock in and stock out values")
		}
	}

	// Prepare update data
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"updated_at": time.Now(),
	}

	// Update provided fields
	if productStockID != nil {
		updateData["product_stock_id"] = *productStockID
	}
	if productID != nil {
		updateData["product_id"] = *productID
	}
	if stockIn != nil {
		updateData["stock_in"] = *stockIn
	}
	if stockOut != nil {
		updateData["stock_out"] = *stockOut
	}
	if quantity != nil {
		updateData["quantity"] = *quantity
	}

	err = s.itemRepo.UpdateProductItem(id, updateData)
	if err != nil {
		return nil, err
	}

	// Create tracking record if needed (commented for now)
	// if s.trackService != nil {
	//   // Add tracking logic here when track service is fixed
	// }

	// Return updated item
	return s.itemRepo.GetProductItemByID(id)
}

func (s *ProductItemService) DeleteProductItem(id uint, userID uint) error {
	// Check if item exists
	// existingItem, err := s.itemRepo.GetProductItemModelByID(id)
	// if err != nil {
	// 	return err
	// }

	// Create tracking record before deletion
	// if s.trackService != nil {
	// 	trackReq := CreateProductItemTrackRequest{
	// 		ProductItemID: id,
	// 		StockIn:       existingItem.StockIn,
	// 		StockOut:      existingItem.StockOut,
	// 		Quantity:      existingItem.Quantity,
	// 		Action:        "DELETE",
	// 	}
	// 	_, _ = s.trackService.CreateProductItemTrack(trackReq, userID)
	// }

	return s.itemRepo.DeleteProductItemWithAudit(id, userID)
}

// Additional business logic methods
func (s *ProductItemService) GetProductItemsByBatch(batchID uint) (interface{}, error) {
	if batchID == 0 {
		return nil, errors.New("invalid batch ID")
	}

	// Check if batch exists
	batchExists, err := s.stockRepo.CheckProductBatchExists(batchID)
	if err != nil {
		return nil, err
	}
	if !batchExists {
		return nil, errors.New("product batch not found")
	}

	return []interface{}{}, nil // Placeholder - method not implemented in repository
}

func (s *ProductItemService) ProcessStockMovement(stockID uint, stockIn, stockOut *float64, locationID *uint, userID uint) (interface{}, error) {
	if stockID == 0 {
		return nil, errors.New("invalid stock ID")
	}

	if (stockIn == nil || *stockIn <= 0) && (stockOut == nil || *stockOut <= 0) {
		return nil, errors.New("either stock in or stock out must be provided")
	}

	if stockIn != nil && stockOut != nil && *stockIn > 0 && *stockOut > 0 {
		return nil, errors.New("cannot process both stock in and stock out in the same transaction")
	}

	// Get stock details
	stock, err := s.stockRepo.GetProductStockModelByID(stockID)
	if err != nil {
		return nil, err
	}

	// Create product item using individual parameters
	return s.CreateProductItem(
		stockID,
		stock.ProductID,
		stock.ProductBatchID,
		locationID,
		stockIn,
		stockOut,
		nil, // quantity will be calculated from stockIn/stockOut
		userID,
	)
}
