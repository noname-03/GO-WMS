package repository

import (
	"myapp/database"
	"myapp/internal/model"
	"time"
)

type ProductBatchRepository struct{}

// productBatchWithDetailsResponse struct untuk response dengan product, category, dan brand name
type productBatchWithDetailsResponse struct {
	ID           uint      `json:"id"`
	ProductID    uint      `json:"productId"`
	ProductName  string    `json:"productName"`
	CategoryID   uint      `json:"categoryId"`
	CategoryName string    `json:"categoryName"`
	BrandID      uint      `json:"brandId"`
	BrandName    string    `json:"brandName"`
	CodeBatch    string    `json:"codeBatch"`
	ExpDate      time.Time `json:"expDate"` // Hidden dari JSON
	Description  *string   `json:"description"`
}

func NewProductBatchRepository() *ProductBatchRepository {
	return &ProductBatchRepository{}
}

func (r *ProductBatchRepository) GetAllProductBatches() ([]productBatchWithDetailsResponse, error) {
	var batches []productBatchWithDetailsResponse

	result := database.DB.Table("product_batches pb").
		Select("pb.id, pb.product_id, p.name as product_name, c.id as category_id, c.name as category_name, b.id as brand_id, b.name as brand_name, pb.code_batch, pb.exp_date, pb.description").
		Joins("LEFT JOIN products p ON pb.product_id = p.id AND p.deleted_at IS NULL").
		Joins("LEFT JOIN categories c ON p.category_id = c.id AND c.deleted_at IS NULL").
		Joins("LEFT JOIN brands b ON c.brand_id = b.id AND b.deleted_at IS NULL").
		Where("pb.deleted_at IS NULL").
		Order("pb.code_batch ASC").
		Find(&batches)

	return batches, result.Error
}

func (r *ProductBatchRepository) GetProductBatchesByProduct(productID uint) ([]productBatchWithDetailsResponse, error) {
	var batches []productBatchWithDetailsResponse

	result := database.DB.Table("product_batches pb").
		Select("pb.id, pb.product_id, p.name as product_name, c.id as category_id, c.name as category_name, b.id as brand_id, b.name as brand_name, pb.code_batch, pb.exp_date, pb.description").
		Joins("INNER JOIN products p ON pb.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN categories c ON p.category_id = c.id AND c.deleted_at IS NULL").
		Joins("INNER JOIN brands b ON c.brand_id = b.id AND b.deleted_at IS NULL").
		Where("pb.product_id = ? AND pb.deleted_at IS NULL", productID).
		Order("pb.code_batch ASC").
		Find(&batches)

	return batches, result.Error
}

func (r *ProductBatchRepository) GetProductBatchByID(id uint) (productBatchWithDetailsResponse, error) {
	var batch productBatchWithDetailsResponse

	result := database.DB.Table("product_batches pb").
		Select("pb.id, pb.product_id, p.name as product_name, c.id as category_id, c.name as category_name, b.id as brand_id, b.name as brand_name, pb.code_batch, pb.exp_date, pb.description").
		Joins("INNER JOIN products p ON pb.product_id = p.id AND p.deleted_at IS NULL").
		Joins("INNER JOIN categories c ON p.category_id = c.id AND c.deleted_at IS NULL").
		Joins("INNER JOIN brands b ON c.brand_id = b.id AND b.deleted_at IS NULL").
		Where("pb.id = ? AND pb.deleted_at IS NULL", id).
		First(&batch)

	return batch, result.Error
}

// GetProductBatchModelByID returns model.ProductBatch for service operations
func (r *ProductBatchRepository) GetProductBatchModelByID(id uint) (model.ProductBatch, error) {
	var batch model.ProductBatch
	result := database.DB.Where("id = ?", id).First(&batch)
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
