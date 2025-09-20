package repository

import (
	"myapp/database"
	"myapp/internal/model"
)

type BrandRepository struct{}

func NewBrandRepository() *BrandRepository {
	return &BrandRepository{}
}

func (r *BrandRepository) GetAllBrands() ([]model.Brand, error) {
	var brands []model.Brand
	result := database.DB.Find(&brands)
	return brands, result.Error
}

func (r *BrandRepository) GetBrandByID(id uint) (model.Brand, error) {
	var brand model.Brand
	result := database.DB.First(&brand, id)
	return brand, result.Error
}

func (r *BrandRepository) CreateBrand(brand *model.Brand) error {
	return database.DB.Create(brand).Error
}

func (r *BrandRepository) UpdateBrand(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.Brand{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *BrandRepository) DeleteBrandWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the brand
	updateData := map[string]interface{}{
		"user_updt": userID,
	}

	// Update the audit field first
	err := database.DB.Model(&model.Brand{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.Brand{}, id).Error
}

func (r *BrandRepository) CheckBrandExists(name string) (bool, error) {
	var count int64
	query := database.DB.Model(&model.Brand{}).Unscoped().Where("name ILIKE ?", name)

	result := query.Count(&count)
	return count > 0, result.Error
}
