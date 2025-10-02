package repository

import (
	"myapp/database"
	"myapp/internal/model"
	"time"
)

type ProductUnitRepository struct{}

// productUnitResponse struct untuk response product unit
type productUnitResponse struct {
	ID               uint     `json:"id"`
	ProductID        uint     `json:"productId"`
	ProductName      string   `json:"productName"`
	LocationID       uint     `json:"locationId"`
	LocationName     string   `json:"locationName"`
	ProductBatchID   uint     `json:"productBatchId"`
	ProductBatchName string   `json:"productBatchCode"`
	Name             *string  `json:"name"`
	Quantity         *float64 `json:"quantity"`
	UnitPrice        *string  `json:"unitPrice"`
	Barcode          *string  `json:"barcode"`
	Description      *string  `json:"description"`
}

type productUnitOriResponse struct {
	ProductId      uint     `json:"product_id"`
	LocationId     uint     `json:"location_id"`
	ProductBatchId uint     `json:"product_batch_id"`
	Name           *string  `json:"name"`
	Quantity       *float64 `json:"quantity"`
	UnitPrice      *float64 `json:"unit_price"`
	Barcode        *string  `json:"barcode"`
	Description    *string  `json:"description"`
}

func NewProductUnitRepository() *ProductUnitRepository {
	return &ProductUnitRepository{}
}

func (r *ProductUnitRepository) GetAllProductUnits() ([]productUnitResponse, error) {
	var units []productUnitResponse
	result := database.DB.Table("product_units pu").
		Select("pu.id, pu.product_id, p.name as product_name, pu.location_id, l.name as location_name, pu.product_batch_id, pb.code_batch as product_batch_name, pu.name, pu.quantity, pu.unit_price, pu.barcode, pu.description").
		Joins("LEFT JOIN products p ON pu.product_id = p.id AND p.deleted_at IS NULL").
		Joins("LEFT JOIN locations l ON pu.location_id = l.id AND l.deleted_at IS NULL").
		Joins("LEFT JOIN product_batches pb ON pu.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Where("pu.deleted_at IS NULL").
		Order("pu.created_at DESC").
		Find(&units)
	return units, result.Error
}

func (r *ProductUnitRepository) GetProductUnitsByProduct(productID uint) ([]productUnitResponse, error) {
	var units []productUnitResponse
	result := database.DB.Table("product_units pu").
		Select("pu.id, pu.product_id, p.name as product_name, pu.location_id, l.name as location_name, pu.product_batch_id, pb.code_batch as product_batch_name, pu.name, pu.quantity, pu.unit_price, pu.barcode, pu.description").
		Joins("LEFT JOIN products p ON pu.product_id = p.id AND p.deleted_at IS NULL").
		Joins("LEFT JOIN locations l ON pu.location_id = l.id AND l.deleted_at IS NULL").
		Joins("LEFT JOIN product_batches pb ON pu.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Where("pu.deleted_at IS NULL AND pu.product_id = ?", productID).
		Order("pu.created_at DESC").
		Find(&units)
	return units, result.Error
}

func (r *ProductUnitRepository) GetProductUnitByID(id uint) (productUnitResponse, error) {
	var unit productUnitResponse
	result := database.DB.Table("product_units pu").
		Select("pu.id, pu.product_id, p.name as product_name, pu.location_id, l.name as location_name, pu.product_batch_id, pb.code_batch as product_batch_name, pu.name, pu.quantity, pu.unit_price, pu.barcode, pu.description").
		Joins("LEFT JOIN products p ON pu.product_id = p.id AND p.deleted_at IS NULL").
		Joins("LEFT JOIN locations l ON pu.location_id = l.id AND l.deleted_at IS NULL").
		Joins("LEFT JOIN product_batches pb ON pu.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Where("pu.deleted_at IS NULL AND pu.id = ?", id).
		First(&unit)
	return unit, result.Error
}

// GetProductUnitById return model.ProductUnit for internal use (not for response)
func (r *ProductUnitRepository) GetProductUnitByIDModel(id uint) (model.ProductUnit, error) {
	var unit model.ProductUnit
	result := database.DB.Where("deleted_at IS NULL AND id = ?", id).First(&unit)
	return unit, result.Error
}

func (r *ProductUnitRepository) CheckProductExists(productID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.Product{}).Where("id = ? AND deleted_at IS NULL", productID).Count(&count)
	return count > 0, result.Error
}

func (r *ProductUnitRepository) CheckProductBatchExists(productBatchID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.ProductBatch{}).Where("id = ? AND deleted_at IS NULL", productBatchID).Count(&count)
	return count > 0, result.Error
}

func (r *ProductUnitRepository) CheckProductUnitExists(productID uint, locationID uint, name string, excludeID uint) (bool, error) {
	var count int64
	query := database.DB.Model(&model.ProductUnit{}).Where("product_id = ? AND name = ? AND deleted_at IS NULL", productID, name)

	if locationID != 0 {
		query = query.Where("location_id = ?", locationID)
	} else {
		query = query.Where("location_id IS NULL")
	}

	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}
	result := query.Count(&count)
	return count > 0, result.Error
}

func (r *ProductUnitRepository) CheckBarcodeProductUnitExists(productID uint, locationID uint, barcode string, excludeID uint) (bool, error) {
	var count int64
	query := database.DB.Model(&model.ProductUnit{}).Where("product_id = ? AND barcode = ? AND deleted_at IS NULL", productID, barcode)

	if locationID != 0 {
		query = query.Where("location_id = ?", locationID)
	} else {
		query = query.Where("location_id IS NULL")
	}

	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}
	result := query.Count(&count)
	return count > 0, result.Error
}

func (r *ProductUnitRepository) CreateProductUnit(unit *model.ProductUnit) error {
	return database.DB.Create(unit).Error
}

func (r *ProductUnitRepository) UpdateProductUnit(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.ProductUnit{}).Where("id = ?", id).Updates(updateData).Error
}
func (r *ProductUnitRepository) DeleteProductUnitWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the unit
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"updated_at": time.Now(),
	}
	// Update the audit field first
	err := database.DB.Model(&model.ProductUnit{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}
	// Then soft delete the unit
	return database.DB.Model(&model.ProductUnit{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}

func (r *ProductUnitRepository) CheckBarcodeExists(barcode string) (bool, error) {
	var count int64
	query := database.DB.Model(&model.ProductUnit{}).Where("barcode = ? AND deleted_at IS NULL", barcode)
	result := query.Count(&count)
	return count > 0, result.Error
}

// GetProductUnitByBarcode returns ProductUnit data based on barcode
func (r *ProductUnitRepository) GetProductUnitByBarcode(barcode string) (productUnitOriResponse, error) {
	var unit productUnitOriResponse
	result := database.DB.Table("product_units pu").
		Select("pu.id, pu.product_id, pu.location_id, pu.product_batch_id, pu.name, pu.quantity, pu.unit_price, pu.barcode, pu.description").
		Where("pu.barcode = ? AND pu.deleted_at IS NULL", barcode).
		First(&unit)
	return unit, result.Error
}

// GetDeletedProductUnits returns all soft deleted product units
func (r *ProductUnitRepository) GetDeletedProductUnits() ([]productUnitResponse, error) {
	var units []productUnitResponse

	result := database.DB.Table("product_units pu").
		Select("pu.id, pu.product_id, p.name as product_name, pu.location_id, l.name as location_name, pu.product_batch_id, pb.code_batch as product_batch_name, pu.name, pu.quantity, pu.unit_price, pu.barcode, pu.description").
		Joins("LEFT JOIN products p ON pu.product_id = p.id").
		Joins("LEFT JOIN locations l ON pu.location_id = l.id").
		Joins("LEFT JOIN product_batches pb ON pu.product_batch_id = pb.id").
		Where("pu.deleted_at IS NOT NULL").
		Order("pu.deleted_at DESC").
		Find(&units)

	return units, result.Error
}

// RestoreProductUnit restores a soft deleted product unit
func (r *ProductUnitRepository) RestoreProductUnit(id uint, userID uint) error {
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"deleted_at": nil,
	}
	return database.DB.Unscoped().Model(&model.ProductUnit{}).Where("id = ?", id).Updates(updateData).Error
}
