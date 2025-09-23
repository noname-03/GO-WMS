package repository

import (
	"myapp/database"
	"myapp/internal/model"
)

type ProductRepository struct{}

// productWithBrandCategoryResponse struct untuk response dengan brand dan category name
type productWithBrandCategoryResponse struct {
	ID           uint    `json:"id"`
	BrandID      uint    `json:"brandId"`
	CategoryID   uint    `json:"categoryId"`
	BrandName    string  `json:"brandName"`
	CategoryName string  `json:"categoryName"`
	Name         string  `json:"name"`
	Description  *string `json:"description"`
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (r *ProductRepository) GetAllProducts() ([]productWithBrandCategoryResponse, error) {
	var products []productWithBrandCategoryResponse

	result := database.DB.Table("products p").
		Select("p.id, c.brand_id, p.category_id, b.name as brand_name, c.name as category_name, p.name, p.description").
		Joins("LEFT JOIN categories c ON p.category_id = c.id AND c.deleted_at IS NULL").
		Joins("LEFT JOIN brands b ON c.brand_id = b.id AND b.deleted_at IS NULL").
		Where("p.deleted_at IS NULL").
		Order("p.name ASC").
		Find(&products)

	return products, result.Error
}

func (r *ProductRepository) GetProductsByCategory(categoryID uint) ([]productWithBrandCategoryResponse, error) {
	var products []productWithBrandCategoryResponse

	result := database.DB.Table("products p").
		Select("p.id, c.brand_id, p.category_id, b.name as brand_name, c.name as category_name, p.name, p.description").
		Joins("INNER JOIN categories c ON p.category_id = c.id AND c.deleted_at IS NULL").
		Joins("INNER JOIN brands b ON c.brand_id = b.id AND b.deleted_at IS NULL").
		Where("p.category_id = ? AND p.deleted_at IS NULL", categoryID).
		Order("p.name ASC").
		Find(&products)

	return products, result.Error
}

func (r *ProductRepository) GetProductByID(id uint) (productWithBrandCategoryResponse, error) {
	var product productWithBrandCategoryResponse

	result := database.DB.Table("products p").
		Select("p.id, c.brand_id, p.category_id, b.name as brand_name, c.name as category_name, p.name, p.description").
		Joins("INNER JOIN categories c ON p.category_id = c.id AND c.deleted_at IS NULL").
		Joins("INNER JOIN brands b ON c.brand_id = b.id AND b.deleted_at IS NULL").
		Where("p.id = ? AND p.deleted_at IS NULL", id).
		First(&product)

	return product, result.Error
}

// GetProductModelByID returns model.Product for service operations
func (r *ProductRepository) GetProductModelByID(id uint) (model.Product, error) {
	var product model.Product
	result := database.DB.Where("id = ?", id).First(&product)
	return product, result.Error
}

func (r *ProductRepository) CreateProduct(product *model.Product) error {
	return database.DB.Create(product).Error
}

func (r *ProductRepository) UpdateProduct(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.Product{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *ProductRepository) DeleteProductWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the product
	updateData := map[string]interface{}{
		"user_updt": userID,
	}

	// Update the audit field first
	err := database.DB.Model(&model.Product{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.Product{}, id).Error
}

func (r *ProductRepository) CheckProductExists(name string, categoryID uint) (bool, error) {
	var count int64
	query := database.DB.Model(&model.Product{}).Unscoped().Where("name ILIKE ? AND category_id = ?", name, categoryID)
	result := query.Count(&count)
	return count > 0, result.Error
}

func (r *ProductRepository) CheckCategoryExists(categoryID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.Category{}).Where("id = ?", categoryID).Count(&count)
	return count > 0, result.Error
}
