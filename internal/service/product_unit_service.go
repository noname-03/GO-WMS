package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
)

type ProductUnitService struct {
	productUnitRepo  *repository.ProductUnitRepository
	trackUnitService *ProductUnitTrackService
}

func NewProductUnitService() *ProductUnitService {
	return &ProductUnitService{
		productUnitRepo:  repository.NewProductUnitRepository(),
		trackUnitService: NewProductUnitTrackService(),
	}
}

// Business logic methods
func (s *ProductUnitService) GetAllProductUnits() (interface{}, error) {
	return s.productUnitRepo.GetAllProductUnits()
}

func (s *ProductUnitService) GetProductUnitsByProduct(productID uint) (interface{}, error) {
	if productID == 0 {
		return nil, errors.New("invalid product ID")
	}

	// Check if product exists
	productExists, err := s.productUnitRepo.CheckProductExists(productID)
	if err != nil {
		return nil, err
	}
	if !productExists {
		return nil, errors.New("product not found")
	}

	return s.productUnitRepo.GetProductUnitsByProduct(productID)
}

func (s *ProductUnitService) GetProductUnitByID(id uint) (interface{}, error) {
	unit, err := s.productUnitRepo.GetProductUnitByID(id)
	if err != nil {
		return nil, err
	}
	return unit, nil
}

func (s *ProductUnitService) CreateProductUnit(productID uint, name *string, quantity *float64, unitPrice *float64, barcode *string, description *string, userID uint) (interface{}, error) {
	if productID == 0 {
		return nil, errors.New("product ID is required")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if product exists
	productExists, err := s.productUnitRepo.CheckProductExists(productID)
	if err != nil {
		return nil, err
	}
	if !productExists {
		return nil, errors.New("product not found")
	}

	// Check if product unit with same name exists for this product (if name is provided)
	if name != nil && *name != "" {
		exists, err := s.productUnitRepo.CheckProductUnitExists(productID, *name, 0)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("product unit with this name already exists for this product")
		}
	}

	if barcode != nil && *barcode != "" {
		exists, err := s.productUnitRepo.CheckBarcodeProductUnitExists(productID, *barcode, 0)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("product unit with this barcode already exists for this product")
		}
	}

	// Basic validation
	if quantity != nil && *quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}

	if unitPrice != nil && *unitPrice < 0 {
		return nil, errors.New("UnitPrice cannot be negative")
	}

	productUnit := &model.ProductUnit{
		ProductID:   productID,
		Name:        name,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		Barcode:     barcode,
		Description: description,
		UserIns:     &userID,
	}

	err = s.productUnitRepo.CreateProductUnit(productUnit)
	if err != nil {
		return nil, err
	}

	// Fetch the created product unit with product details
	createdProductUnit, err := s.productUnitRepo.GetProductUnitByIDModel(productUnit.ID)
	if err == nil {
		err = s.trackUnitService.TrackCreate(createdProductUnit, userID)
		if err != nil {
			return nil, err
		}
	}
	return createdProductUnit, nil
}

func (s *ProductUnitService) UpdateProductUnit(id uint, productID uint, name *string, quantity *float64, UnitPrice *float64, barcode *string, description *string, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid product unit ID")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if product unit exists
	oldBatch, err := s.productUnitRepo.GetProductUnitByIDModel(id)
	if err != nil {
		return nil, errors.New("product unit not found")
	}

	// If product ID is being changed, check if new product exists
	if productID != 0 && productID != oldBatch.ProductID {
		productExists, err := s.productUnitRepo.CheckProductExists(productID)
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

	// Check if new name conflicts with existing product units for this product
	if name != nil && *name != "" {
		currentName := ""
		if oldBatch.Name != nil {
			currentName = *oldBatch.Name
		}
		if *name != currentName || productID != oldBatch.ProductID {
			exists, err := s.productUnitRepo.CheckProductUnitExists(productID, *name, id)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, errors.New("product unit name already in use for this product")
			}
		}
	}

	// Check if new name conflicts with existing product units for this product
	if barcode != nil && *barcode != "" {
		currentbarcode := ""
		if oldBatch.Barcode != nil {
			currentbarcode = *oldBatch.Barcode
		}
		if *barcode != currentbarcode || productID != oldBatch.ProductID {
			exists, err := s.productUnitRepo.CheckBarcodeProductUnitExists(productID, *barcode, 0)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, errors.New("product unit barcode already in use for this product")
			}
		}
	}

	// Basic validation
	if quantity != nil && *quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}

	if UnitPrice != nil && *UnitPrice < 0 {
		return nil, errors.New("UnitPrice cannot be negative")
	}

	// Update the model fields
	updateData := make(map[string]interface{})
	if productID != oldBatch.ProductID {
		updateData["product_id"] = productID
	}
	if name != nil {
		updateData["name"] = name
	}
	if quantity != nil {
		updateData["quantity"] = quantity
	}
	if UnitPrice != nil {
		updateData["unit_price"] = UnitPrice
	}
	if barcode != nil {
		updateData["barcode"] = barcode
	}
	if description != nil {
		updateData["description"] = description
	}
	// Always set the user who updated
	updateData["user_updt"] = userID

	err = s.productUnitRepo.UpdateProductUnit(id, updateData)
	if err != nil {
		return nil, err
	}

	// Create tracking record for update using actual changes
	err = s.trackUnitService.TrackUpdateFromChangesProductUnit(updateData, oldBatch, userID)
	if err != nil {
		return nil, err
	}

	updatedProductUnit, err := s.productUnitRepo.GetProductUnitByID(id)
	if err != nil {
		return nil, err
	}
	return updatedProductUnit, nil
}

func (s *ProductUnitService) DeleteProductUnit(id uint, userID uint) error {
	if id == 0 {
		return errors.New("invalid product unit ID")
	}

	if userID == 0 {
		return errors.New("user ID is required for audit trail")
	}

	// Check if product unit exists
	_, err := s.productUnitRepo.GetProductUnitByIDModel(id)
	if err != nil {
		return errors.New("product unit not found")
	}

	return s.productUnitRepo.DeleteProductUnitWithAudit(id, userID)
}
