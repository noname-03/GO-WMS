package repository

import (
	"myapp/database"
	"myapp/internal/model"
	"time"
)

type ProductUnitRepository struct{}

// productUnitResponse struct untuk response product unit
type productUnitResponse struct {
	ID          uint     `json:"id"`
	ProductID   uint     `json:"productId"`
	ProductName string   `json:"productName"`
	Name        *string  `json:"name"`
	Quantity    *float64 `json:"quantity"`
	UnitPrice   *string  `json:"UnitPrice"`
	Barcode     *string  `json:"barcode"`
	Description *string  `json:"description"`
}

func NewProductUnitRepository() *ProductUnitRepository {
	return &ProductUnitRepository{}
}

func (r *ProductUnitRepository) GetAllProductUnits() ([]productUnitResponse, error) {
	var units []productUnitResponse
	result := database.DB.Table("product_units pu").
		Select("pu.id, pu.product_id, p.name as product_name, pu.name, pu.quantity, pu.unit_price, pu.barcode, pu.description").
		Joins("LEFT JOIN products p ON pu.product_id = p.id AND p.deleted_at IS NULL").
		Where("pu.deleted_at IS NULL").
		Order("pu.created_at DESC").
		Find(&units)
	return units, result.Error
}

func (r *ProductUnitRepository) GetProductUnitsByProduct(productID uint) ([]productUnitResponse, error) {
	var units []productUnitResponse
	result := database.DB.Table("product_units pu").
		Select("pu.id, pu.product_id, p.name as product_name, pu.name, pu.quantity, pu.unit_price, pu.barcode, pu.description").
		Joins("LEFT JOIN products p ON pu.product_id = p.id AND p.deleted_at IS NULL").
		Where("pu.deleted_at IS NULL AND pu.product_id = ?", productID).
		Order("pu.created_at DESC").
		Find(&units)
	return units, result.Error
}

func (r *ProductUnitRepository) GetProductUnitByID(id uint) (productUnitResponse, error) {
	var unit productUnitResponse
	result := database.DB.Table("product_units pu").
		Select("pu.id, pu.product_id, p.name as product_name, pu.name, pu.quantity, pu.unit_price, pu.barcode, pu.description").
		Joins("LEFT JOIN products p ON pu.product_id = p.id AND p.deleted_at IS NULL").
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

func (r *ProductUnitRepository) CheckProductUnitExists(productID uint, name string, excludeID uint) (bool, error) {
	var count int64
	query := database.DB.Model(&model.ProductUnit{}).Where("product_id = ? AND name = ? AND deleted_at IS NULL", productID, name)
	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}
	result := query.Count(&count)
	return count > 0, result.Error
}

func (r *ProductUnitRepository) CheckBarcodeProductUnitExists(productID uint, barcode string, excludeID uint) (bool, error) {
	var count int64
	query := database.DB.Model(&model.ProductUnit{}).Where("product_id = ? AND barcode = ? AND deleted_at IS NULL", productID, barcode)
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

func (r *ProductUnitRepository) CheckBarcodeExists(barcode string, excludeID uint) (bool, error) {
	var count int64
	query := database.DB.Model(&model.ProductUnit{}).Where("barcode = ? AND deleted_at IS NULL", barcode)
	if excludeID != 0 {
		query = query.Where("id != ?", excludeID)
	}
	result := query.Count(&count)
	return count > 0, result.Error
}
