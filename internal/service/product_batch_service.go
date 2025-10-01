package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
	"time"
)

type ProductBatchService struct {
	batchRepo    *repository.ProductBatchRepository
	trackService *ProductBatchTrackService
}

func NewProductBatchService() *ProductBatchService {
	return &ProductBatchService{
		batchRepo:    repository.NewProductBatchRepository(),
		trackService: NewProductBatchTrackService(),
	}
}

// Business logic methods
func (s *ProductBatchService) GetAllProductBatches() (interface{}, error) {
	return s.batchRepo.GetAllProductBatches()
}

func (s *ProductBatchService) GetProductBatchesByProduct(productID uint) (interface{}, error) {
	if productID == 0 {
		return nil, errors.New("invalid product ID")
	}

	// Check if product exists
	productExists, err := s.batchRepo.CheckProductExists(productID)
	if err != nil {
		return nil, err
	}
	if !productExists {
		return nil, errors.New("product not found")
	}

	return s.batchRepo.GetProductBatchesByProduct(productID)
}

func (s *ProductBatchService) GetProductBatchByID(id uint) (interface{}, error) {
	batch, err := s.batchRepo.GetProductBatchByID(id)
	if err != nil {
		return nil, err
	}
	return batch, nil
}

func (s *ProductBatchService) CreateProductBatch(productID uint, codeBatch *string, unitPrice *float64, expDate time.Time, description *string, userID uint) (interface{}, error) {
	if productID == 0 {
		return nil, errors.New("product ID is required")
	}

	if expDate.IsZero() {
		return nil, errors.New("expiry date is required")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if product exists
	productExists, err := s.batchRepo.CheckProductExists(productID)
	if err != nil {
		return nil, err
	}
	if !productExists {
		return nil, errors.New("product not found")
	}

	// Note: No duplicate check for batches since they can have different exp_date, code_batch, etc.

	batch := &model.ProductBatch{
		ProductID:   productID,
		CodeBatch:   codeBatch,
		UnitPrice:   unitPrice,
		ExpDate:     expDate,
		Description: description,
		UserIns:     &userID, // Set pointer to userID
	}

	err = s.batchRepo.CreateProductBatch(batch)
	if err != nil {
		return nil, err
	}

	// Fetch the created batch with relationships
	createdBatch, err := s.batchRepo.GetProductBatchByID(batch.ID)
	if err != nil {
		return nil, err
	}

	// Create tracking record for creation
	batchModel, err := s.batchRepo.GetProductBatchModelByID(batch.ID)
	if err == nil {
		err = s.trackService.TrackCreate(batchModel, userID)
		if err != nil {
			// Log error but don't fail the creation
			// Consider using a logger here
			// log.Printf("Failed to create tracking record: %v", err)
		}
	}

	return createdBatch, nil
}

func (s *ProductBatchService) UpdateProductBatch(id uint, productID uint, codeBatch *string, unitPrice *float64, expDate time.Time, description *string, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid product batch ID")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if batch exists using model for business logic
	oldBatch, err := s.batchRepo.GetProductBatchModelByID(id)
	if err != nil {
		return nil, errors.New("product batch not found")
	}

	// If product ID is being changed, check if new product exists
	if productID != 0 && productID != oldBatch.ProductID {
		productExists, err := s.batchRepo.CheckProductExists(productID)
		if err != nil {
			return nil, err
		}
		if !productExists {
			return nil, errors.New("product not found")
		}
	}

	// Use existing product if not provided
	if productID == 0 {
		productID = oldBatch.ProductID
	}

	// Prepare update data with audit trail
	updateData := make(map[string]interface{})
	if productID != oldBatch.ProductID {
		updateData["product_id"] = productID
	}
	if codeBatch != nil {
		updateData["code_batch"] = codeBatch
	}
	if unitPrice != nil {
		updateData["unit_price"] = unitPrice
	}
	if !expDate.IsZero() {
		updateData["exp_date"] = expDate
	}
	if description != nil {
		updateData["description"] = description
	}
	// Always set the user who updated
	updateData["user_updt"] = userID

	err = s.batchRepo.UpdateProductBatch(id, updateData)
	if err != nil {
		return nil, err
	}

	// Create tracking record for update using actual changes
	err = s.trackService.TrackUpdateFromChanges(updateData, oldBatch, userID)
	if err != nil {
		// Log error but don't fail the update
		// Consider using a logger here
		// log.Printf("Failed to create tracking record: %v", err)
	}

	updatedBatch, err := s.batchRepo.GetProductBatchByID(id)
	if err != nil {
		return nil, err
	}

	return updatedBatch, nil
}

func (s *ProductBatchService) DeleteProductBatch(id uint, userID uint) error {
	if id == 0 {
		return errors.New("invalid product batch ID")
	}

	if userID == 0 {
		return errors.New("user ID is required for audit trail")
	}

	// Check if batch exists using model for tracking
	batchToDelete, err := s.batchRepo.GetProductBatchModelByID(id)
	if err != nil {
		return errors.New("product batch not found")
	}

	// Create tracking record for deletion (before actual deletion)
	err = s.trackService.TrackDelete(batchToDelete, userID)
	if err != nil {
		// Log error but don't fail the deletion
		// Consider using a logger here
		// log.Printf("Failed to create tracking record: %v", err)
	}

	return s.batchRepo.DeleteProductBatchWithAudit(id, userID)
}

// GetDeletedProductBatches returns all soft deleted product batches
func (s *ProductBatchService) GetDeletedProductBatches() (interface{}, error) {
	return s.batchRepo.GetDeletedProductBatches()
}

// RestoreProductBatch restores a soft deleted product batch
func (s *ProductBatchService) RestoreProductBatch(id uint, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid product batch ID")
	}
	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	err := s.batchRepo.RestoreProductBatch(id, userID)
	if err != nil {
		return nil, err
	}

	restoredBatch, err := s.batchRepo.GetProductBatchByID(id)
	if err != nil {
		return nil, err
	}
	return restoredBatch, nil
}
