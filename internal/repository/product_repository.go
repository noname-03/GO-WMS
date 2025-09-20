package repository

import (
	"myapp/database"
	"myapp/internal/model"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (r *ProductRepository) GetAllProducts() ([]model.Product, error) {
	var products []model.Product
	result := database.DB.Preload("Category").Preload("Category.Brand").Find(&products)
	return products, result.Error
}

func (r *ProductRepository) GetProductsByCategory(categoryID uint) ([]model.Product, error) {
	var products []model.Product
	result := database.DB.Preload("Category").Preload("Category.Brand").Where("category_id = ?", categoryID).Find(&products)
	return products, result.Error
}

func (r *ProductRepository) GetProductByID(id uint) (model.Product, error) {
	var product model.Product
	result := database.DB.Preload("Category").Preload("Category.Brand").Preload("ProductBatches").First(&product, id)
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
