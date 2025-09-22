package repository

import (
	"myapp/database"
	"myapp/internal/model"
)

type CategoryRepository struct{}

// categoryBasicResponse struct untuk GetCategoriesByBrand tanpa relasi Brand
type categoryBasicResponse struct {
	ID          uint    `json:"id"`
	BrandID     uint    `json:"brandId"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (r *CategoryRepository) GetAllCategories() ([]categoryBasicResponse, error) {
	var categories []categoryBasicResponse

	result := database.DB.Table("categories").
		Select("id, brand_id, name, description").
		Where("deleted_at IS NULL").
		Order("name ASC").
		Find(&categories)

	return categories, result.Error
}

func (r *CategoryRepository) GetCategoriesByBrand(brandID uint) ([]categoryBasicResponse, error) {
	var categories []categoryBasicResponse

	result := database.DB.Table("categories").
		Select("id, brand_id, name, description").
		Where("brand_id = ? AND deleted_at IS NULL", brandID).
		Find(&categories)

	return categories, result.Error
}

func (r *CategoryRepository) GetCategoryByID(id uint) (model.Category, error) {
	var category model.Category
	result := database.DB.Preload("Brand").First(&category, id)
	return category, result.Error
}

func (r *CategoryRepository) CreateCategory(category *model.Category) error {
	return database.DB.Create(category).Error
}

func (r *CategoryRepository) UpdateCategory(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.Category{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *CategoryRepository) DeleteCategoryWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the category
	updateData := map[string]interface{}{
		"user_updt": userID,
	}

	// Update the audit field first
	err := database.DB.Model(&model.Category{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.Category{}, id).Error
}

func (r *CategoryRepository) CheckCategoryExists(name string, brandID uint) (bool, error) {
	var count int64
	query := database.DB.Model(&model.Category{}).Unscoped().Where("name ILIKE ? AND brand_id = ?", name, brandID)

	result := query.Count(&count)
	return count > 0, result.Error
}

func (r *CategoryRepository) CheckBrandExists(brandID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.Brand{}).Where("id = ?", brandID).Count(&count)
	return count > 0, result.Error
}
