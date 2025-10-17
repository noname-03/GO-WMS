package repository

import (
	"myapp/database"
	"myapp/internal/model"
)

type DeliveryOrderItemRepository struct{}

// deliveryOrderItemWithDetailsResponse struct untuk response dengan product name
type deliveryOrderItemWithDetailsResponse struct {
	ID              uint     `json:"id"`
	DeliveryOrderID uint     `json:"deliveryOrderId"`
	ProductID       uint     `json:"productId"`
	ProductName     string   `json:"productName"`
	QtyDelivered    *float64 `json:"qtyDelivered"`
	Description     *string  `json:"description"`
}

func NewDeliveryOrderItemRepository() *DeliveryOrderItemRepository {
	return &DeliveryOrderItemRepository{}
}

func (r *DeliveryOrderItemRepository) GetAllDeliveryOrderItems() ([]deliveryOrderItemWithDetailsResponse, error) {
	var items []deliveryOrderItemWithDetailsResponse

	result := database.DB.Table("delivery_order_items doi").
		Select("doi.id, doi.delivery_order_id, doi.product_id, p.name as product_name, doi.qty_delivered, doi.description").
		Joins("LEFT JOIN products p ON doi.product_id = p.id AND p.deleted_at IS NULL").
		Where("doi.deleted_at IS NULL").
		Order("doi.id ASC").
		Find(&items)

	return items, result.Error
}

func (r *DeliveryOrderItemRepository) GetDeliveryOrderItemsByDeliveryOrder(deliveryOrderID uint) ([]deliveryOrderItemWithDetailsResponse, error) {
	var items []deliveryOrderItemWithDetailsResponse

	result := database.DB.Table("delivery_order_items doi").
		Select("doi.id, doi.delivery_order_id, doi.product_id, p.name as product_name, doi.qty_delivered, doi.description").
		Joins("INNER JOIN products p ON doi.product_id = p.id AND p.deleted_at IS NULL").
		Where("doi.delivery_order_id = ? AND doi.deleted_at IS NULL", deliveryOrderID).
		Order("doi.id ASC").
		Find(&items)

	return items, result.Error
}

func (r *DeliveryOrderItemRepository) GetDeliveryOrderItemByID(id uint) (deliveryOrderItemWithDetailsResponse, error) {
	var item deliveryOrderItemWithDetailsResponse

	result := database.DB.Table("delivery_order_items doi").
		Select("doi.id, doi.delivery_order_id, doi.product_id, p.name as product_name, doi.qty_delivered, doi.description").
		Joins("INNER JOIN products p ON doi.product_id = p.id AND p.deleted_at IS NULL").
		Where("doi.id = ? AND doi.deleted_at IS NULL", id).
		First(&item)

	return item, result.Error
}

// GetDeliveryOrderItemModelByID returns model.DeliveryOrderItem for service operations
func (r *DeliveryOrderItemRepository) GetDeliveryOrderItemModelByID(id uint) (model.DeliveryOrderItem, error) {
	var item model.DeliveryOrderItem
	result := database.DB.Where("id = ?", id).First(&item)
	return item, result.Error
}

func (r *DeliveryOrderItemRepository) CreateDeliveryOrderItem(item *model.DeliveryOrderItem) error {
	return database.DB.Create(item).Error
}

func (r *DeliveryOrderItemRepository) UpdateDeliveryOrderItem(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.DeliveryOrderItem{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *DeliveryOrderItemRepository) DeleteDeliveryOrderItemWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the item
	updateData := map[string]interface{}{
		"user_updt": userID,
	}

	// Update the audit field first
	err := database.DB.Model(&model.DeliveryOrderItem{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.DeliveryOrderItem{}, id).Error
}

func (r *DeliveryOrderItemRepository) CheckDeliveryOrderExists(deliveryOrderID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.DeliveryOrder{}).Where("id = ?", deliveryOrderID).Count(&count)
	return count > 0, result.Error
}

func (r *DeliveryOrderItemRepository) CheckProductExists(productID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.Product{}).Where("id = ?", productID).Count(&count)
	return count > 0, result.Error
}

// GetDeletedDeliveryOrderItems returns all soft deleted delivery order items
func (r *DeliveryOrderItemRepository) GetDeletedDeliveryOrderItems() ([]deliveryOrderItemWithDetailsResponse, error) {
	var items []deliveryOrderItemWithDetailsResponse

	result := database.DB.Table("delivery_order_items doi").
		Select("doi.id, doi.delivery_order_id, doi.product_id, p.name as product_name, doi.qty_delivered, doi.description").
		Joins("LEFT JOIN products p ON doi.product_id = p.id").
		Where("doi.deleted_at IS NOT NULL").
		Order("doi.deleted_at DESC").
		Find(&items)

	return items, result.Error
}

// RestoreDeliveryOrderItem restores a soft deleted delivery order item
func (r *DeliveryOrderItemRepository) RestoreDeliveryOrderItem(id uint, userID uint) error {
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"deleted_at": nil,
	}
	return database.DB.Unscoped().Model(&model.DeliveryOrderItem{}).Where("id = ?", id).Updates(updateData).Error
}
