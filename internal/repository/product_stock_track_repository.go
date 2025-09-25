package repository

import (
	"myapp/database"
	"myapp/internal/model"
	"time"
)

type ProductStockTrackRepository struct{}

// productStockTrackResponse struct untuk response dengan relasi detail
type productStockTrackResponse struct {
	ID               uint      `json:"id"`
	ProductStockID   uint      `json:"productStockId"`
	ProductID        uint      `json:"productId"`
	ProductName      string    `json:"productName"`
	ProductBatchID   uint      `json:"productBatchId"`
	ProductBatchCode string    `json:"productBatchCode"`
	DateTrack        time.Time `json:"dateTrack"`
	Quantity         *float64  `json:"quantity"`
	Operation        string    `json:"operation"`
	Stock            *float64  `json:"stock"`
}

func NewProductStockTrackRepository() *ProductStockTrackRepository {
	return &ProductStockTrackRepository{}
}

func (r *ProductStockTrackRepository) GetAllProductStockTracks() ([]productStockTrackResponse, error) {
	var tracks []productStockTrackResponse

	result := database.DB.Table("product_stock_tracks pst").
		Select("pst.id, pst.product_stock_id, pst.product_id, p.name as product_name, pst.product_batch_id, pb.code_batch as product_batch_code, pst.date_track, pst.quantity, pst.operation, pst.stock").
		Joins("INNER JOIN products p ON pst.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pst.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Where("pst.deleted_at IS NULL").
		Order("pst.date_track DESC").
		Find(&tracks)

	return tracks, result.Error
}

func (r *ProductStockTrackRepository) GetProductStockTracksByStock(stockID uint) ([]productStockTrackResponse, error) {
	var tracks []productStockTrackResponse

	result := database.DB.Table("product_stock_tracks pst").
		Select("pst.id, pst.product_stock_id, pst.product_id, p.name as product_name, pst.product_batch_id, pb.code_batch as product_batch_code, pst.date_track, pst.quantity, pst.operation, pst.stock").
		Joins("INNER JOIN products p ON pst.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pst.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Where("pst.deleted_at IS NULL AND pst.product_stock_id = ?", stockID).
		Order("pst.date_track DESC").
		Find(&tracks)

	return tracks, result.Error
}

func (r *ProductStockTrackRepository) GetProductStockTracksByProduct(productID uint) ([]productStockTrackResponse, error) {
	var tracks []productStockTrackResponse

	result := database.DB.Table("product_stock_tracks pst").
		Select("pst.id, pst.product_stock_id, pst.product_id, p.name as product_name, pst.product_batch_id, pb.code_batch as product_batch_code, pst.date_track, pst.quantity, pst.operation, pst.stock").
		Joins("INNER JOIN products p ON pst.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pst.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Where("pst.deleted_at IS NULL AND pst.product_id = ?", productID).
		Order("pst.date_track DESC").
		Find(&tracks)

	return tracks, result.Error
}

func (r *ProductStockTrackRepository) GetProductStockTracksByDateRange(startDate, endDate time.Time) ([]productStockTrackResponse, error) {
	var tracks []productStockTrackResponse

	result := database.DB.Table("product_stock_tracks pst").
		Select("pst.id, pst.product_stock_id, pst.product_id, p.name as product_name, pst.product_batch_id, pb.code_batch as product_batch_code, pst.date_track, pst.quantity, pst.operation, pst.stock").
		Joins("INNER JOIN products p ON pst.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pst.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Where("pst.deleted_at IS NULL AND pst.date_track BETWEEN ? AND ?", startDate, endDate).
		Order("pst.date_track DESC").
		Find(&tracks)

	return tracks, result.Error
}

func (r *ProductStockTrackRepository) GetProductStockTrackByID(id uint) (productStockTrackResponse, error) {
	var track productStockTrackResponse

	result := database.DB.Table("product_stock_tracks pst").
		Select("pst.id, pst.product_stock_id, pst.product_id, p.name as product_name, pst.product_batch_id, pb.code_batch as product_batch_code, pst.date_track, pst.quantity, pst.operation, pst.stock").
		Joins("INNER JOIN products p ON pst.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN product_batches pb ON pst.product_batch_id = pb.id AND pb.deleted_at IS NULL").
		Where("pst.deleted_at IS NULL AND pst.id = ?", id).
		First(&track)

	return track, result.Error
}

// GetProductStockTrackModelByID returns model.ProductStockTrack for service operations
func (r *ProductStockTrackRepository) GetProductStockTrackModelByID(id uint) (model.ProductStockTrack, error) {
	var track model.ProductStockTrack
	result := database.DB.Where("id = ?", id).First(&track)
	return track, result.Error
}

func (r *ProductStockTrackRepository) CreateProductStockTrack(track *model.ProductStockTrack) error {
	return database.DB.Create(track).Error
}

func (r *ProductStockTrackRepository) UpdateProductStockTrack(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.ProductStockTrack{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *ProductStockTrackRepository) DeleteProductStockTrackWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the track
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"updated_at": time.Now(),
	}

	// Update the audit field first
	err := database.DB.Model(&model.ProductStockTrack{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.ProductStockTrack{}, id).Error
}

func (r *ProductStockTrackRepository) CheckProductStockExists(stockID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.ProductStock{}).Where("id = ?", stockID).Count(&count)
	return count > 0, result.Error
}