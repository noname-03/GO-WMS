package repository

import (
	"myapp/database"
	"myapp/internal/model"
	"time"
)

type ProductStockRepository struct{}

// productStockResponse struct untuk response dengan product, batch, dan location name
type productStockResponse struct {
	ID               uint     `json:"id"`
	ProductBatchID   uint     `json:"productBatchId"`
	ProductBatchCode string   `json:"productBatchCode"`
	ProductID        uint     `json:"productId"`
	ProductName      string   `json:"productName"`
	LocationID       *uint    `json:"locationId"`
	LocationName     *string  `json:"locationName"`
	Quantity         *float64 `json:"quantity"`
}

func NewProductStockRepository() *ProductStockRepository {
	return &ProductStockRepository{}
}

func (r *ProductStockRepository) GetAllProductStocks() ([]productStockResponse, error) {
	var stocks []productStockResponse

	result := database.DB.Table("product_stocks ps").
		Select("ps.id, ps.product_batch_id, pb.code_batch as product_batch_code, ps.product_id, p.name as product_name, ps.location_id, l.name as location_name, ps.quantity").
		Joins("INNER JOIN product_batches pb ON ps.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Joins("INNER JOIN products p ON ps.product_id = p.id AND p.deleted_at IS NULL").
		Joins("LEFT JOIN locations l ON ps.location_id = l.id AND l.deleted_at IS NULL").
		Where("ps.deleted_at IS NULL").
		Order("ps.created_at DESC").
		Find(&stocks)

	return stocks, result.Error
}

func (r *ProductStockRepository) GetProductStocksByProduct(productID uint) ([]productStockResponse, error) {
	var stocks []productStockResponse

	result := database.DB.Table("product_stocks ps").
		Select("ps.id, ps.product_batch_id, pb.code_batch as product_batch_code, ps.product_id, p.name as product_name, ps.location_id, l.name as location_name, ps.quantity").
		Joins("INNER JOIN product_batches pb ON ps.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Joins("INNER JOIN products p ON ps.product_id = p.id AND p.deleted_at IS NULL").
		Joins("LEFT JOIN locations l ON ps.location_id = l.id AND l.deleted_at IS NULL").
		Where("ps.deleted_at IS NULL AND ps.product_id = ?", productID).
		Order("ps.created_at DESC").
		Find(&stocks)

	return stocks, result.Error
}

func (r *ProductStockRepository) GetProductStockByID(id uint) (productStockResponse, error) {
	var stock productStockResponse

	result := database.DB.Table("product_stocks ps").
		Select("ps.id, ps.product_batch_id, pb.code_batch as product_batch_code, ps.product_id, p.name as product_name, ps.location_id, l.name as location_name, ps.quantity").
		Joins("INNER JOIN product_batches pb ON ps.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Joins("INNER JOIN products p ON ps.product_id = p.id AND p.deleted_at IS NULL").
		Joins("LEFT JOIN locations l ON ps.location_id = l.id AND l.deleted_at IS NULL").
		Where("ps.deleted_at IS NULL AND ps.id = ?", id).
		First(&stock)

	return stock, result.Error
}

// GetProductStockModelByID returns model.ProductStock for service operations
func (r *ProductStockRepository) GetProductStockModelByID(id uint) (model.ProductStock, error) {
	var stock model.ProductStock
	result := database.DB.Where("id = ?", id).First(&stock)
	return stock, result.Error
}

func (r *ProductStockRepository) CreateProductStock(stock *model.ProductStock) error {
	return database.DB.Create(stock).Error
}

func (r *ProductStockRepository) UpdateProductStock(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.ProductStock{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *ProductStockRepository) DeleteProductStockWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the stock
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"updated_at": time.Now(),
	}

	// Update the audit field first
	err := database.DB.Model(&model.ProductStock{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.ProductStock{}, id).Error
}

func (r *ProductStockRepository) CheckProductExists(productID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.Product{}).Where("id = ?", productID).Count(&count)
	return count > 0, result.Error
}

func (r *ProductStockRepository) CheckProductBatchExists(batchID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.ProductBatch{}).Where("id = ?", batchID).Count(&count)
	return count > 0, result.Error
}
