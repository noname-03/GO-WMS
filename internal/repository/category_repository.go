package repository

import (
	"myapp/database"
	"myapp/internal/model"
)

type CategoryRepository struct{}

// categoryWithBrandResponse struct untuk response dengan brand name
type categoryWithBrandResponse struct {
	ID          uint    `json:"id"`
	BrandID     uint    `json:"brandId"`
	BrandName   string  `json:"brandName"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (r *CategoryRepository) GetAllCategories() ([]categoryWithBrandResponse, error) {
	var categories []categoryWithBrandResponse

	result := database.DB.Debug().Table("categories c").
		Select("c.id, c.brand_id, c.name, c.description, b.name as brand_name").
		Joins("LEFT JOIN brands b ON c.brand_id = b.id AND b.deleted_at IS NULL").
		Where("c.deleted_at IS NULL").
		Order("c.name ASC").
		Find(&categories)

	return categories, result.Error
}

func (r *CategoryRepository) GetCategoriesByBrand(brandID uint) ([]categoryWithBrandResponse, error) {
	var categories []categoryWithBrandResponse

	result := database.DB.Table("categories c").
		Select("c.id, c.brand_id, c.name, c.description, b.name as brand_name").
		Joins("INNER JOIN brands b ON c.brand_id = b.id AND b.deleted_at IS NULL").
		Where("c.brand_id = ? AND c.deleted_at IS NULL", brandID).
		Order("c.name ASC").
		Find(&categories)

	return categories, result.Error
}

func (r *CategoryRepository) GetCategoryByID(id uint) (categoryWithBrandResponse, error) {
	var category categoryWithBrandResponse

	result := database.DB.Table("categories c").
		Select("c.id, c.brand_id, c.name, c.description, b.name as brand_name").
		Joins("INNER JOIN brands b ON c.brand_id = b.id AND b.deleted_at IS NULL").
		Where("c.id = ? AND c.deleted_at IS NULL", id).
		First(&category)

	return category, result.Error
}

// GetCategoryModelByID returns model.Category for service operations
func (r *CategoryRepository) GetCategoryModelByID(id uint) (model.Category, error) {
	var category model.Category
	result := database.DB.Preload("Brand").Where("id = ?", id).First(&category)
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
