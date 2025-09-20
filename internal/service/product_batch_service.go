package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
	"time"
)

type ProductBatchService struct {
	batchRepo *repository.ProductBatchRepository
}

func NewProductBatchService() *ProductBatchService {
	return &ProductBatchService{
		batchRepo: repository.NewProductBatchRepository(),
	}
}

// Business logic methods
func (s *ProductBatchService) GetAllProductBatches() ([]model.ProductBatch, error) {
	return s.batchRepo.GetAllProductBatches()
}

func (s *ProductBatchService) GetProductBatchesByProduct(productID uint) ([]model.ProductBatch, error) {
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

func (s *ProductBatchService) GetProductBatchByID(id uint) (*model.ProductBatch, error) {
	batch, err := s.batchRepo.GetProductBatchByID(id)
	if err != nil {
		return nil, err
	}
	return &batch, nil
}

func (s *ProductBatchService) CreateProductBatch(productID uint, codeBatch *string, unitPrice *float64, expDate time.Time, description *string, userID uint) (*model.ProductBatch, error) {
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

	return &createdBatch, nil
}

func (s *ProductBatchService) UpdateProductBatch(id uint, productID uint, codeBatch *string, unitPrice *float64, expDate time.Time, description *string, userID uint) (*model.ProductBatch, error) {
	if id == 0 {
		return nil, errors.New("invalid product batch ID")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if batch exists
	batch, err := s.batchRepo.GetProductBatchByID(id)
	if err != nil {
		return nil, errors.New("product batch not found")
	}

	// If product ID is being changed, check if new product exists
	if productID != 0 && productID != batch.ProductID {
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
		productID = batch.ProductID
	}

	// Prepare update data with audit trail
	updateData := make(map[string]interface{})
	if productID != batch.ProductID {
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

	updatedBatch, err := s.batchRepo.GetProductBatchByID(id)
	if err != nil {
		return nil, err
	}
	return &updatedBatch, nil
}

func (s *ProductBatchService) DeleteProductBatch(id uint, userID uint) error {
	if id == 0 {
		return errors.New("invalid product batch ID")
	}

	if userID == 0 {
		return errors.New("user ID is required for audit trail")
	}

	// Check if batch exists
	_, err := s.batchRepo.GetProductBatchByID(id)
	if err != nil {
		return errors.New("product batch not found")
	}

	return s.batchRepo.DeleteProductBatchWithAudit(id, userID)
}
