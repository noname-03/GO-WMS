package repository

import (
	"myapp/database"
	"myapp/internal/model"
)

type ProductUnitTrackRepository struct{}

// productUnitTrackWithDetailsResponse struct untuk response dengan product unit dan product name
type productUnitTrackWithDetailsResponse struct {
	ID            uint    `json:"id"`
	ProductUnitID uint    `json:"productUnitId"`
	ProductID     uint    `json:"productId"`
	ProductName   string  `json:"productName"`
	UnitName      *string `json:"unitName"`
	UnitType      *string `json:"unitType"`
	Description   *string `json:"description"`
}

func NewProductUnitTrackRepository() *ProductUnitTrackRepository {
	return &ProductUnitTrackRepository{}
}

func (r *ProductUnitTrackRepository) GetAllProductUnitTracks() ([]productUnitTrackWithDetailsResponse, error) {
	var tracks []productUnitTrackWithDetailsResponse

	result := database.DB.Table("product_unit_tracks put").
		Select("put.id, put.product_unit_id, pu.product_id, p.name as product_name, pu.name as unit_name, pu.unit_price, pu.barcode, put.description").
		Joins("LEFT JOIN product_units pu ON put.product_unit_id = pu.id AND pu.deleted_at IS NULL").
		Joins("LEFT JOIN products p ON pu.product_id = p.id AND p.deleted_at IS NULL").
		Where("put.deleted_at IS NULL").
		Order("put.created_at DESC").
		Find(&tracks)

	return tracks, result.Error
}

func (r *ProductUnitTrackRepository) GetProductUnitTracksByProductUnit(productUnitID uint) ([]productUnitTrackWithDetailsResponse, error) {
	var tracks []productUnitTrackWithDetailsResponse

	result := database.DB.Table("product_unit_tracks put").
		Select("put.id, put.product_unit_id, pu.product_id, p.name as product_name, pu.name as unit_name, pu.unit_price, pu.barcode, put.description").
		Joins("INNER JOIN product_units pu ON put.product_unit_id = pu.id AND pu.deleted_at IS NULL").
		Joins("INNER JOIN products p ON pu.product_id = p.id AND p.deleted_at IS NULL").
		Where("put.product_unit_id = ? AND put.deleted_at IS NULL", productUnitID).
		Order("put.created_at DESC").
		Find(&tracks)

	return tracks, result.Error
}

func (r *ProductUnitTrackRepository) GetProductUnitTracksByProduct(productID uint) ([]productUnitTrackWithDetailsResponse, error) {
	var tracks []productUnitTrackWithDetailsResponse

	result := database.DB.Table("product_unit_tracks put").
		Select("put.id, put.product_unit_id, pu.product_id, p.name as product_name, pu.name as unit_name, pu.unit_price, pu.barcode, put.description").
		Joins("INNER JOIN product_units pu ON put.product_unit_id = pu.id AND pu.deleted_at IS NULL").
		Joins("INNER JOIN products p ON pu.product_id = p.id AND p.deleted_at IS NULL").
		Where("pu.product_id = ? AND put.deleted_at IS NULL", productID).
		Order("put.created_at DESC").
		Find(&tracks)

	return tracks, result.Error
}

func (r *ProductUnitTrackRepository) GetProductUnitTrackByID(id uint) (productUnitTrackWithDetailsResponse, error) {
	var track productUnitTrackWithDetailsResponse

	result := database.DB.Table("product_unit_tracks put").
		Select("put.id, put.product_unit_id, pu.product_id, p.name as product_name, pu.name as unit_name, pu.unit_price, pu.barcode, put.description").
		Joins("INNER JOIN product_units pu ON put.product_unit_id = pu.id AND pu.deleted_at IS NULL").
		Joins("INNER JOIN products p ON pu.product_id = p.id AND p.deleted_at IS NULL").
		Where("put.id = ? AND put.deleted_at IS NULL", id).
		First(&track)

	return track, result.Error
}

// GetProductUnitTrackModelByID returns model.ProductUnitTrack for service operations
func (r *ProductUnitTrackRepository) GetProductUnitTrackModelByID(id uint) (model.ProductUnitTrack, error) {
	var track model.ProductUnitTrack
	result := database.DB.Where("id = ?", id).First(&track)
	return track, result.Error
}

func (r *ProductUnitTrackRepository) CreateProductUnitTrack(track *model.ProductUnitTrack) error {
	return database.DB.Create(track).Error
}

func (r *ProductUnitTrackRepository) UpdateProductUnitTrack(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.ProductUnitTrack{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *ProductUnitTrackRepository) DeleteProductUnitTrackWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the track
	updateData := map[string]interface{}{
		"user_updt": userID,
	}

	// Update the audit field first
	err := database.DB.Model(&model.ProductUnitTrack{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.ProductUnitTrack{}, id).Error
}

func (r *ProductUnitTrackRepository) CheckProductUnitExists(productUnitID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.ProductUnit{}).Where("id = ?", productUnitID).Count(&count)
	return count > 0, result.Error
}
