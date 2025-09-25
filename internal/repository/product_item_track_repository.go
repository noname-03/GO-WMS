package repository

import (
	"myapp/database"
	"myapp/internal/model"
	"time"
)

type ProductItemTrackRepository struct{}

// productItemTrackResponse struct untuk response dengan relasi detail
type productItemTrackResponse struct {
	ID               uint      `json:"id"`
	ProductItemID    uint      `json:"productItemId"`
	ProductStockID   uint      `json:"productStockId"`
	ProductID        uint      `json:"productId"`
	ProductName      string    `json:"productName"`
	ProductBatchID   uint      `json:"productBatchId"`
	ProductBatchCode string    `json:"productBatchCode"`
	DateTrack        time.Time `json:"dateTrack"`
	UnitPrice        *float64  `json:"unitPrice"`
	Quantity         *float64  `json:"quantity"`
	Operation        string    `json:"operation"`
	Stock            *float64  `json:"stock"`
}

func NewProductItemTrackRepository() *ProductItemTrackRepository {
	return &ProductItemTrackRepository{}
}

func (r *ProductItemTrackRepository) GetAllProductItemTracks() ([]productItemTrackResponse, error) {
	var tracks []productItemTrackResponse

	result := database.DB.Table("product_item_tracks pit").
		Select("pit.id, pit.product_item_id, pit.product_stock_id, pit.product_id, p.name as product_name, pit.product_batch_id, pb.code_batch as product_batch_code, pit.date_track, pit.unit_price, pit.quantity, pit.operation, pit.stock").
		Joins("INNER JOIN products p ON pit.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pit.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Where("pit.deleted_at IS NULL").
		Order("pit.date_track DESC").
		Find(&tracks)

	return tracks, result.Error
}

func (r *ProductItemTrackRepository) GetProductItemTracksByItem(itemID uint) ([]productItemTrackResponse, error) {
	var tracks []productItemTrackResponse

	result := database.DB.Table("product_item_tracks pit").
		Select("pit.id, pit.product_item_id, pit.product_stock_id, pit.product_id, p.name as product_name, pit.product_batch_id, pb.code_batch as product_batch_code, pit.date_track, pit.unit_price, pit.quantity, pit.operation, pit.stock").
		Joins("INNER JOIN products p ON pit.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pit.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Where("pit.deleted_at IS NULL AND pit.product_item_id = ?", itemID).
		Order("pit.date_track DESC").
		Find(&tracks)

	return tracks, result.Error
}

func (r *ProductItemTrackRepository) GetProductItemTracksByStock(stockID uint) ([]productItemTrackResponse, error) {
	var tracks []productItemTrackResponse

	result := database.DB.Table("product_item_tracks pit").
		Select("pit.id, pit.product_item_id, pit.product_stock_id, pit.product_id, p.name as product_name, pit.product_batch_id, pb.code_batch as product_batch_code, pit.date_track, pit.unit_price, pit.quantity, pit.operation, pit.stock").
		Joins("INNER JOIN products p ON pit.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pit.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Where("pit.deleted_at IS NULL AND pit.product_stock_id = ?", stockID).
		Order("pit.date_track DESC").
		Find(&tracks)

	return tracks, result.Error
}

func (r *ProductItemTrackRepository) GetProductItemTracksByProduct(productID uint) ([]productItemTrackResponse, error) {
	var tracks []productItemTrackResponse

	result := database.DB.Table("product_item_tracks pit").
		Select("pit.id, pit.product_item_id, pit.product_stock_id, pit.product_id, p.name as product_name, pit.product_batch_id, pb.code_batch as product_batch_code, pit.date_track, pit.unit_price, pit.quantity, pit.operation, pit.stock").
		Joins("INNER JOIN products p ON pit.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pit.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Where("pit.deleted_at IS NULL AND pit.product_id = ?", productID).
		Order("pit.date_track DESC").
		Find(&tracks)

	return tracks, result.Error
}

func (r *ProductItemTrackRepository) GetProductItemTracksByDateRange(startDate, endDate time.Time) ([]productItemTrackResponse, error) {
	var tracks []productItemTrackResponse

	result := database.DB.Table("product_item_tracks pit").
		Select("pit.id, pit.product_item_id, pit.product_stock_id, pit.product_id, p.name as product_name, pit.product_batch_id, pb.code_batch as product_batch_code, pit.date_track, pit.unit_price, pit.quantity, pit.operation, pit.stock").
		Joins("INNER JOIN products p ON pit.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pit.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Where("pit.deleted_at IS NULL AND pit.date_track BETWEEN ? AND ?", startDate, endDate).
		Order("pit.date_track DESC").
		Find(&tracks)

	return tracks, result.Error
}

func (r *ProductItemTrackRepository) GetProductItemTrackByID(id uint) (productItemTrackResponse, error) {
	var track productItemTrackResponse

	result := database.DB.Table("product_item_tracks pit").
		Select("pit.id, pit.product_item_id, pit.product_stock_id, pit.product_id, p.name as product_name, pit.product_batch_id, pb.code_batch as product_batch_code, pit.date_track, pit.unit_price, pit.quantity, pit.operation, pit.stock").
		Joins("INNER JOIN products p ON pit.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pit.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Where("pit.deleted_at IS NULL AND pit.id = ?", id).
		First(&track)

	return track, result.Error
}

// GetProductItemTrackModelByID returns model.ProductItemTrack for service operations
func (r *ProductItemTrackRepository) GetProductItemTrackModelByID(id uint) (model.ProductItemTrack, error) {
	var track model.ProductItemTrack
	result := database.DB.Where("id = ?", id).First(&track)
	return track, result.Error
}

func (r *ProductItemTrackRepository) CreateProductItemTrack(track *model.ProductItemTrack) error {
	return database.DB.Create(track).Error
}

func (r *ProductItemTrackRepository) UpdateProductItemTrack(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.ProductItemTrack{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *ProductItemTrackRepository) DeleteProductItemTrackWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the track
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"updated_at": time.Now(),
	}

	// Update the audit field first
	err := database.DB.Model(&model.ProductItemTrack{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.ProductItemTrack{}, id).Error
}

func (r *ProductItemTrackRepository) CheckProductItemExists(itemID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.ProductItem{}).Where("id = ?", itemID).Count(&count)
	return count > 0, result.Error
}


// GetTracksByOperation returns tracks filtered by operation type (Plus/Minus/In/Out)
func (r *ProductItemTrackRepository) GetTracksByOperation(operation string) ([]productItemTrackResponse, error) {
	var tracks []productItemTrackResponse

	result := database.DB.Table("product_item_tracks pit").
		Select("pit.id, pit.product_item_id, pit.product_stock_id, pit.product_id, p.name as product_name, pit.product_batch_id, pb.code_batch as product_batch_code, pit.date_track, pit.unit_price, pit.quantity, pit.operation, pit.stock").
		Joins("INNER JOIN products p ON pit.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pit.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Where("pit.deleted_at IS NULL AND pit.operation = ?", operation).
		Order("pit.date_track DESC").
		Find(&tracks)

	return tracks, result.Error
}

// GetValueReportByProduct returns value report grouped by product with total value calculations
func (r *ProductItemTrackRepository) GetValueReportByProduct() ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	rows, err := database.DB.Table("product_item_tracks pit").
		Select("pit.product_id, p.name as product_name, COUNT(*) as total_transactions, SUM(pit.quantity * pit.unit_price) as total_value, AVG(pit.unit_price) as avg_unit_price").
		Joins("INNER JOIN products p ON pit.product_id = p.id AND p.deleted_at IS NULL").
		Where("pit.deleted_at IS NULL AND pit.unit_price IS NOT NULL AND pit.quantity IS NOT NULL").
		Group("pit.product_id, p.name").
		Order("total_value DESC").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var productID uint
		var productName string
		var totalTransactions int64
		var totalValue, avgUnitPrice *float64

		err := rows.Scan(&productID, &productName, &totalTransactions, &totalValue, &avgUnitPrice)
		if err != nil {
			return nil, err
		}

		result := map[string]interface{}{
			"product_id":         productID,
			"product_name":       productName,
			"total_transactions": totalTransactions,
			"total_value":        totalValue,
			"avg_unit_price":     avgUnitPrice,
		}
		results = append(results, result)
	}

	return results, nil
}
