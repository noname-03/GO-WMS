package repository

import (
	"myapp/database"
	"myapp/internal/model"
)

type ProductBatchRepository struct{}

func NewProductBatchRepository() *ProductBatchRepository {
	return &ProductBatchRepository{}
}

func (r *ProductBatchRepository) GetAllProductBatches() ([]model.ProductBatch, error) {
	var batches []model.ProductBatch
	result := database.DB.Preload("Product").Preload("Product.Category").Preload("Product.Category.Brand").Find(&batches)
	return batches, result.Error
}

func (r *ProductBatchRepository) GetProductBatchesByProduct(productID uint) ([]model.ProductBatch, error) {
	var batches []model.ProductBatch
	result := database.DB.Preload("Product").Preload("Product.Category").Preload("Product.Category.Brand").Where("product_id = ?", productID).Find(&batches)
	return batches, result.Error
}

func (r *ProductBatchRepository) GetProductBatchByID(id uint) (model.ProductBatch, error) {
	var batch model.ProductBatch
	result := database.DB.Preload("Product").Preload("Product.Category").Preload("Product.Category.Brand").First(&batch, id)
	return batch, result.Error
}

func (r *ProductBatchRepository) CreateProductBatch(batch *model.ProductBatch) error {
	return database.DB.Create(batch).Error
}

func (r *ProductBatchRepository) UpdateProductBatch(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.ProductBatch{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *ProductBatchRepository) DeleteProductBatchWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the batch
	updateData := map[string]interface{}{
		"user_updt": userID,
	}

	// Update the audit field first
	err := database.DB.Model(&model.ProductBatch{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.ProductBatch{}, id).Error
}

func (r *ProductBatchRepository) CheckProductExists(productID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.Product{}).Where("id = ?", productID).Count(&count)
	return count > 0, result.Error
}
