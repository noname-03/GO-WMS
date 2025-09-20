package repository

import (
	"myapp/database"
	"myapp/internal/model"
)

type ProductBatchTrackRepository struct{}

func NewProductBatchTrackRepository() *ProductBatchTrackRepository {
	return &ProductBatchTrackRepository{}
}

// GetAllTracks retrieves all product batch tracking records with relationships
func (r *ProductBatchTrackRepository) GetAllTracks() ([]model.ProductBatchTrack, error) {
	var tracks []model.ProductBatchTrack
	err := database.DB.Preload("ProductBatch").Preload("ProductBatch.Product").Preload("ProductBatch.Product.Category").Preload("ProductBatch.Product.Category.Brand").Preload("Creator").Preload("Updater").Find(&tracks).Error
	return tracks, err
}

// GetTracksByProductBatchID retrieves all tracking records for a specific product batch
func (r *ProductBatchTrackRepository) GetTracksByProductBatchID(productBatchID uint) ([]model.ProductBatchTrack, error) {
	var tracks []model.ProductBatchTrack
	err := database.DB.Where("product_batch_id = ?", productBatchID).Preload("ProductBatch").Preload("ProductBatch.Product").Preload("ProductBatch.Product.Category").Preload("ProductBatch.Product.Category.Brand").Preload("Creator").Preload("Updater").Order("created_at DESC").Find(&tracks).Error
	return tracks, err
}

// GetTrackByID retrieves a specific tracking record by ID
func (r *ProductBatchTrackRepository) GetTrackByID(id uint) (model.ProductBatchTrack, error) {
	var track model.ProductBatchTrack
	err := database.DB.Preload("ProductBatch").Preload("ProductBatch.Product").Preload("ProductBatch.Product.Category").Preload("ProductBatch.Product.Category.Brand").Preload("Creator").Preload("Updater").First(&track, id).Error
	return track, err
}

// CreateTrack creates a new tracking record
func (r *ProductBatchTrackRepository) CreateTrack(track *model.ProductBatchTrack) error {
	return database.DB.Create(track).Error
}

// GetTracksByUserID retrieves all tracking records made by a specific user
func (r *ProductBatchTrackRepository) GetTracksByUserID(userID uint) ([]model.ProductBatchTrack, error) {
	var tracks []model.ProductBatchTrack
	err := database.DB.Where("user_inst = ? OR user_updt = ?", userID, userID).Preload("ProductBatch").Preload("ProductBatch.Product").Preload("ProductBatch.Product.Category").Preload("ProductBatch.Product.Category.Brand").Preload("Creator").Preload("Updater").Order("created_at DESC").Find(&tracks).Error
	return tracks, err
}

// GetLatestTrackForProductBatch retrieves the most recent tracking record for a product batch
func (r *ProductBatchTrackRepository) GetLatestTrackForProductBatch(productBatchID uint) (model.ProductBatchTrack, error) {
	var track model.ProductBatchTrack
	err := database.DB.Where("product_batch_id = ?", productBatchID).Preload("ProductBatch").Preload("ProductBatch.Product").Preload("ProductBatch.Product.Category").Preload("ProductBatch.Product.Category.Brand").Preload("Creator").Preload("Updater").Order("created_at DESC").First(&track).Error
	return track, err
}

// GetTracksForMultipleProductBatches retrieves tracking records for multiple product batches
func (r *ProductBatchTrackRepository) GetTracksForMultipleProductBatches(productBatchIDs []uint) ([]model.ProductBatchTrack, error) {
	var tracks []model.ProductBatchTrack
	err := database.DB.Where("product_batch_id IN ?", productBatchIDs).Preload("ProductBatch").Preload("ProductBatch.Product").Preload("ProductBatch.Product.Category").Preload("ProductBatch.Product.Category.Brand").Preload("Creator").Preload("Updater").Order("created_at DESC").Find(&tracks).Error
	return tracks, err
}

// DeleteTrack soft deletes a tracking record
func (r *ProductBatchTrackRepository) DeleteTrack(id uint) error {
	return database.DB.Delete(&model.ProductBatchTrack{}, id).Error
}

// CountTracksByProductBatchID counts tracking records for a specific product batch
func (r *ProductBatchTrackRepository) CountTracksByProductBatchID(productBatchID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&model.ProductBatchTrack{}).Where("product_batch_id = ?", productBatchID).Count(&count).Error
	return count, err
}
