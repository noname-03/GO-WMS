package repository

import (
	"myapp/database"
	"myapp/internal/model"
	"time"
)

type ProductItemRepository struct{}

// productItemResponse struct untuk response dengan relasi detail
type productItemResponse struct {
	ID               uint     `json:"id"`
	ProductStockID   uint     `json:"productStockId"`
	ProductID        uint     `json:"productId"`
	ProductName      string   `json:"productName"`
	ProductBatchID   uint     `json:"productBatchId"`
	ProductBatchCode string   `json:"productBatchCode"`
	LocationID       *uint    `json:"locationId"`
	LocationName     *string  `json:"locationName"`
	StockIn          *float64 `json:"stockIn"`
	StockOut         *float64 `json:"stockOut"`
	Quantity         *float64 `json:"quantity"`
}

func NewProductItemRepository() *ProductItemRepository {
	return &ProductItemRepository{}
}

func (r *ProductItemRepository) GetAllProductItems() ([]productItemResponse, error) {
	var items []productItemResponse

	result := database.DB.Table("product_items pi").
		Select("pi.id, pi.product_stock_id, pi.product_id, p.name as product_name, pi.product_batch_id, pb.code_batch as product_batch_code, ps.location_id, l.name as location_name, pi.stock_in, pi.stock_out, pi.quantity").
		Joins("INNER JOIN product_stocks ps ON pi.product_stock_id = ps.id AND ps.deleted_at IS NULL").
		Joins("INNER JOIN products p ON pi.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pi.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Joins("LEFT JOIN locations l ON ps.location_id = l.id AND l.deleted_at IS NULL").
		Where("pi.deleted_at IS NULL").
		Order("pi.created_at DESC").
		Find(&items)

	return items, result.Error
}

func (r *ProductItemRepository) GetProductItemsByStock(stockID uint) ([]productItemResponse, error) {
	var items []productItemResponse

	result := database.DB.Table("product_items pi").
		Select("pi.id, pi.product_stock_id, pi.product_id, p.name as product_name, pi.product_batch_id, pb.code_batch as product_batch_code, ps.location_id, l.name as location_name, pi.stock_in, pi.stock_out, pi.quantity").
		Joins("INNER JOIN product_stocks ps ON pi.product_stock_id = ps.id AND ps.deleted_at IS NULL").
		Joins("INNER JOIN products p ON pi.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pi.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Joins("LEFT JOIN locations l ON ps.location_id = l.id AND l.deleted_at IS NULL").
		Where("pi.deleted_at IS NULL AND pi.product_stock_id = ?", stockID).
		Order("pi.created_at DESC").
		Find(&items)

	return items, result.Error
}

func (r *ProductItemRepository) GetProductItemsByProduct(productID uint) ([]productItemResponse, error) {
	var items []productItemResponse

	result := database.DB.Table("product_items pi").
		Select("pi.id, pi.product_stock_id, pi.product_id, p.name as product_name, pi.product_batch_id, pb.code_batch as product_batch_code, ps.location_id, l.name as location_name, pi.stock_in, pi.stock_out, pi.quantity").
		Joins("INNER JOIN product_stocks ps ON pi.product_stock_id = ps.id AND ps.deleted_at IS NULL").
		Joins("INNER JOIN products p ON pi.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pi.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Joins("LEFT JOIN locations l ON ps.location_id = l.id AND l.deleted_at IS NULL").
		Where("pi.deleted_at IS NULL AND pi.product_id = ?", productID).
		Order("pi.created_at DESC").
		Find(&items)

	return items, result.Error
}

func (r *ProductItemRepository) GetProductItemsByLocation(locationID uint) ([]productItemResponse, error) {
	var items []productItemResponse

	result := database.DB.Table("product_items pi").
		Select("pi.id, pi.product_stock_id, pi.product_id, p.name as product_name, pi.product_batch_id, pb.code_batch as product_batch_code, ps.location_id, l.name as location_name, pi.stock_in, pi.stock_out, pi.quantity").
		Joins("INNER JOIN product_stocks ps ON pi.product_stock_id = ps.id AND ps.deleted_at IS NULL").
		Joins("INNER JOIN products p ON pi.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pi.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Joins("LEFT JOIN locations l ON ps.location_id = l.id AND l.deleted_at IS NULL").
		Where("pi.deleted_at IS NULL AND pi.location_id = ?", locationID).
		Order("pi.created_at DESC").
		Find(&items)

	return items, result.Error
}

func (r *ProductItemRepository) GetProductItemByID(id uint) (productItemResponse, error) {
	var item productItemResponse

	result := database.DB.Table("product_items pi").
		Select("pi.id, pi.product_stock_id, pi.product_id, p.name as product_name, pi.product_batch_id, pb.code_batch as product_batch_code, ps.location_id, l.name as location_name, pi.stock_in, pi.stock_out, pi.quantity").
		Joins("INNER JOIN product_stocks ps ON pi.product_stock_id = ps.id AND ps.deleted_at IS NULL").
		Joins("INNER JOIN products p ON pi.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pi.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Joins("LEFT JOIN locations l ON ps.location_id = l.id AND l.deleted_at IS NULL").
		Where("pi.deleted_at IS NULL AND pi.id = ?", id).
		First(&item)

	return item, result.Error
}

// GetProductItemModelByID returns model.ProductItem for service operations
func (r *ProductItemRepository) GetProductItemModelByID(id uint) (model.ProductItem, error) {
	var item model.ProductItem
	result := database.DB.Where("id = ?", id).First(&item)
	return item, result.Error
}

func (r *ProductItemRepository) CreateProductItem(item *model.ProductItem) error {
	return database.DB.Create(item).Error
}

func (r *ProductItemRepository) UpdateProductItem(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.ProductItem{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *ProductItemRepository) DeleteProductItemWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the item
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"updated_at": time.Now(),
	}

	// Update the audit field first
	err := database.DB.Model(&model.ProductItem{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.ProductItem{}, id).Error
}

func (r *ProductItemRepository) CheckProductStockExists(stockID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.ProductStock{}).Where("id = ?", stockID).Count(&count)
	return count > 0, result.Error
}

// GetItemsSummaryByProduct returns summary of items grouped by product
func (r *ProductItemRepository) GetItemsSummaryByProduct() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	rows, err := database.DB.Table("product_items pi").
		Select("pi.product_id, p.name as product_name, SUM(pi.stock_in) as total_stock_in, SUM(pi.stock_out) as total_stock_out, SUM(pi.quantity) as total_quantity").
		Joins("INNER JOIN products p ON pi.product_id = p.id AND p.deleted_at IS NULL").
		Where("pi.deleted_at IS NULL").
		Group("pi.product_id, p.name").
		Order("p.name ASC").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var productID uint
		var productName string
		var totalStockIn, totalStockOut, totalQuantity *float64

		err := rows.Scan(&productID, &productName, &totalStockIn, &totalStockOut, &totalQuantity)
		if err != nil {
			return nil, err
		}

		result := map[string]interface{}{
			"product_id":      productID,
			"product_name":    productName,
			"total_stock_in":  totalStockIn,
			"total_stock_out": totalStockOut,
			"total_quantity":  totalQuantity,
		}
		results = append(results, result)
	}

	return results, nil
}
