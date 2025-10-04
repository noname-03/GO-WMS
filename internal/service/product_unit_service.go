package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
	"strings"
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

func (s *ProductUnitService) CreateProductUnit(productID uint, locationID uint, productBatchID uint, name *string, quantity *float64, unitPrice *float64, unitPriceRetail *float64, barcode *string, description *string, userID uint) (interface{}, error) {
	if productID == 0 {
		return nil, errors.New("product ID is required")
	}

	if locationID == 0 {
		return nil, errors.New("location ID is required")
	}

	if productBatchID == 0 {
		return nil, errors.New("product batch ID is required")
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

	// Check if product batch exists
	batchExists, err := s.productUnitRepo.CheckProductBatchExists(productBatchID)
	if err != nil {
		return nil, err
	}
	if !batchExists {
		return nil, errors.New("product batch not found")
	}

	// Check if product unit with same name exists for this product and location (if name is provided)
	if name != nil && *name != "" {
		exists, err := s.productUnitRepo.CheckProductUnitExists(productID, locationID, *name, 0)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("product unit with this name already exists for this product and location")
		}
	}

	// Check if barcode exists for this product and location (if barcode is provided)
	// if barcode != nil && *barcode != "" {
	// 	exists, err := s.productUnitRepo.CheckBarcodeProductUnitExists(productID, locationID, *barcode)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if exists {
	// 		return nil, errors.New("barcode already exists for this product and location")
	// 	}
	// }

	if barcode != nil && *barcode != "" {
		unit, err := s.productUnitRepo.GetProductUnitByBarcode(*barcode)
		if err == nil {
			// Barcode sudah ada, cek apakah product_id sama
			if unit.ProductId != productID {
				return nil, errors.New("barcode already exists for another product")
			}

			// Product_id sama, cek apakah kombinasi name dan location_id sama
			unitName := ""
			if unit.Name != nil {
				unitName = strings.ToLower(*unit.Name)
			}

			currentName := ""
			if name != nil {
				currentName = strings.ToLower(*name)
			}

			if unit.LocationId == locationID && unitName == currentName {
				return nil, errors.New("product unit with same name and location already exists")
			}
		}
		// Jika err != nil berarti barcode tidak ditemukan, jadi aman untuk create
	}

	productUnit := &model.ProductUnit{
		ProductID:       productID,
		LocationID:      locationID,
		ProductBatchID:  productBatchID,
		Name:            name,
		Quantity:        quantity,
		UnitPrice:       unitPrice,
		UnitPriceRetail: unitPriceRetail,
		Barcode:         barcode,
		Description:     description,
		UserIns:         &userID,
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

	res, err := s.productUnitRepo.GetProductUnitByID(productUnit.ID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ProductUnitService) UpdateProductUnit(id uint, productID uint, locationID uint, productBatchID uint, name *string, quantity *float64, UnitPrice *float64, unitPriceRetail *float64, barcode *string, description *string, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid product unit ID")
	}

	if locationID == 0 {
		return nil, errors.New("invalid location ID")
	}

	if productBatchID == 0 {
		return nil, errors.New("product batch ID is required")
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

	// If product batch ID is being changed, check if new product batch exists
	if productBatchID != 0 && productBatchID != oldBatch.ProductBatchID {
		batchExists, err := s.productUnitRepo.CheckProductBatchExists(productBatchID)
		if err != nil {
			return nil, err
		}
		if !batchExists {
			return nil, errors.New("product batch not found")
		}
	}

	// Use existing product if not provided
	if productID == 0 {
		productID = oldBatch.ProductID
	}

	if locationID == 0 {
		locationID = oldBatch.LocationID
	}

	if productBatchID == 0 {
		productBatchID = oldBatch.ProductBatchID
	}

	// Check if new name conflicts with existing product units for this product
	if name != nil && *name != "" {
		currentName := ""
		if oldBatch.Name != nil {
			currentName = *oldBatch.Name
		}
		if *name != currentName || productID != oldBatch.ProductID || (locationID != 0 && (oldBatch.LocationID == 0 || locationID != oldBatch.LocationID)) {
			exists, err := s.productUnitRepo.CheckProductUnitExists(productID, locationID, *name, id)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, errors.New("product unit name already in use for this product and location")
			}
		}
	}

	// Check if new name conflicts with existing product units for this product
	// if barcode != nil && *barcode != "" {
	// 	currentbarcode := ""
	// 	if oldBatch.Barcode != nil {
	// 		currentbarcode = *oldBatch.Barcode
	// 	}
	// 	if *barcode != currentbarcode || productID != oldBatch.ProductID {
	// 		exists, err := s.productUnitRepo.CheckBarcodeProductUnitExists(productID, locationID, *barcode)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		if exists {
	// 			return nil, errors.New("product unit barcode already in use for this product and location")
	// 		}
	// 	}
	// }

	if barcode != nil && *barcode != "" {
		currentbarcode := ""
		if oldBatch.Barcode != nil {
			currentbarcode = *oldBatch.Barcode
		}
		if *barcode != currentbarcode || productID != oldBatch.ProductID {
			unit, err := s.productUnitRepo.GetProductUnitByBarcode(*barcode)
			if err == nil {
				// Barcode sudah ada, cek apakah product_id sama
				if unit.ProductId != productID {
					return nil, errors.New("barcode already exists for another product")
				}

				// Product_id sama, cek apakah kombinasi name dan location_id sama
				unitName := ""
				if unit.Name != nil {
					unitName = strings.ToLower(*unit.Name)
				}

				currentName := ""
				if name != nil {
					currentName = strings.ToLower(*name)
				}

				if unit.LocationId == locationID && unitName == currentName {
					return nil, errors.New("product unit with same name and location already exists")
				}
			}
		}
		// Jika err != nil berarti barcode tidak ditemukan, jadi aman untuk create
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
	if locationID != oldBatch.LocationID {
		updateData["location_id"] = locationID
	}
	if productBatchID != oldBatch.ProductBatchID {
		updateData["product_batch_id"] = productBatchID
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
	if unitPriceRetail != nil {
		updateData["unit_price_retail"] = unitPriceRetail
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
	unitDelete, err := s.productUnitRepo.GetProductUnitByIDModel(id)
	if err != nil {
		return errors.New("product unit not found")
	}

	// Create tracking record for deletion (before actual deletion)
	err = s.trackUnitService.TrackDeleteProductUnit(unitDelete, userID)
	if err != nil {
		// return nil, err
	}

	return s.productUnitRepo.DeleteProductUnitWithAudit(id, userID)
}

// GetDeletedProductUnits returns all soft deleted product units
func (s *ProductUnitService) GetDeletedProductUnits() (interface{}, error) {
	return s.productUnitRepo.GetDeletedProductUnits()
}

// RestoreProductUnit restores a soft deleted product unit
func (s *ProductUnitService) RestoreProductUnit(id uint, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid product unit ID")
	}
	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	err := s.productUnitRepo.RestoreProductUnit(id, userID)
	if err != nil {
		return nil, err
	}

	restoredUnit, err := s.productUnitRepo.GetProductUnitByID(id)
	if err != nil {
		return nil, err
	}
	return restoredUnit, nil
}
