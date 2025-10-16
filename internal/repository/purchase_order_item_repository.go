package repository

import (
	"myapp/database"
	"myapp/internal/model"
)

type PurchaseOrderItemRepository struct{}

// purchaseOrderItemWithDetailsResponse struct untuk response dengan product name
type purchaseOrderItemWithDetailsResponse struct {
	ID              uint     `json:"id"`
	PurchaseOrderID uint     `json:"purchaseOrderId"`
	ProductID       uint     `json:"productId"`
	ProductName     string   `json:"productName"`
	QtyOrdered      *float64 `json:"qtyOrdered"`
	UnitPrice       *float64 `json:"unitPrice"`
	Discount        *float64 `json:"discount"`
	TotalPrice      *float64 `json:"totalPrice"`
	Description     *string  `json:"description"`
}

func NewPurchaseOrderItemRepository() *PurchaseOrderItemRepository {
	return &PurchaseOrderItemRepository{}
}

func (r *PurchaseOrderItemRepository) GetAllPurchaseOrderItems() ([]purchaseOrderItemWithDetailsResponse, error) {
	var items []purchaseOrderItemWithDetailsResponse

	result := database.DB.Table("purchase_order_items poi").
		Select("poi.id, poi.purchase_order_id, poi.product_id, p.name as product_name, poi.qty_ordered, poi.unit_price, poi.discount, poi.total_price, poi.description").
		Joins("LEFT JOIN products p ON poi.product_id = p.id AND p.deleted_at IS NULL").
		Where("poi.deleted_at IS NULL").
		Order("poi.id ASC").
		Find(&items)

	return items, result.Error
}

func (r *PurchaseOrderItemRepository) GetPurchaseOrderItemsByPurchaseOrder(purchaseOrderID uint) ([]purchaseOrderItemWithDetailsResponse, error) {
	var items []purchaseOrderItemWithDetailsResponse

	result := database.DB.Table("purchase_order_items poi").
		Select("poi.id, poi.purchase_order_id, poi.product_id, p.name as product_name, poi.qty_ordered, poi.unit_price, poi.discount, poi.total_price, poi.description").
		Joins("INNER JOIN products p ON poi.product_id = p.id AND p.deleted_at IS NULL").
		Where("poi.purchase_order_id = ? AND poi.deleted_at IS NULL", purchaseOrderID).
		Order("poi.id ASC").
		Find(&items)

	return items, result.Error
}

func (r *PurchaseOrderItemRepository) GetPurchaseOrderItemByID(id uint) (purchaseOrderItemWithDetailsResponse, error) {
	var item purchaseOrderItemWithDetailsResponse

	result := database.DB.Table("purchase_order_items poi").
		Select("poi.id, poi.purchase_order_id, poi.product_id, p.name as product_name, poi.qty_ordered, poi.unit_price, poi.discount, poi.total_price, poi.description").
		Joins("INNER JOIN products p ON poi.product_id = p.id AND p.deleted_at IS NULL").
		Where("poi.id = ? AND poi.deleted_at IS NULL", id).
		First(&item)

	return item, result.Error
}

// GetPurchaseOrderItemModelByID returns model.PurchaseOrderItem for service operations
func (r *PurchaseOrderItemRepository) GetPurchaseOrderItemModelByID(id uint) (model.PurchaseOrderItem, error) {
	var item model.PurchaseOrderItem
	result := database.DB.Where("id = ?", id).First(&item)
	return item, result.Error
}

func (r *PurchaseOrderItemRepository) CreatePurchaseOrderItem(item *model.PurchaseOrderItem) error {
	return database.DB.Create(item).Error
}

func (r *PurchaseOrderItemRepository) UpdatePurchaseOrderItem(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.PurchaseOrderItem{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *PurchaseOrderItemRepository) DeletePurchaseOrderItemWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the item
	updateData := map[string]interface{}{
		"user_updt": userID,
	}

	// Update the audit field first
	err := database.DB.Model(&model.PurchaseOrderItem{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.PurchaseOrderItem{}, id).Error
}

func (r *PurchaseOrderItemRepository) CheckPurchaseOrderExists(purchaseOrderID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.PurchaseOrder{}).Where("id = ?", purchaseOrderID).Count(&count)
	return count > 0, result.Error
}

func (r *PurchaseOrderItemRepository) CheckProductExists(productID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.Product{}).Where("id = ?", productID).Count(&count)
	return count > 0, result.Error
}

// GetDeletedPurchaseOrderItems returns all soft deleted purchase order items
func (r *PurchaseOrderItemRepository) GetDeletedPurchaseOrderItems() ([]purchaseOrderItemWithDetailsResponse, error) {
	var items []purchaseOrderItemWithDetailsResponse

	result := database.DB.Table("purchase_order_items poi").
		Select("poi.id, poi.purchase_order_id, poi.product_id, p.name as product_name, poi.qty_ordered, poi.unit_price, poi.discount, poi.total_price, poi.description").
		Joins("LEFT JOIN products p ON poi.product_id = p.id").
		Where("poi.deleted_at IS NOT NULL").
		Order("poi.deleted_at DESC").
		Find(&items)

	return items, result.Error
}

// RestorePurchaseOrderItem restores a soft deleted purchase order item
func (r *PurchaseOrderItemRepository) RestorePurchaseOrderItem(id uint, userID uint) error {
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"deleted_at": nil,
	}
	return database.DB.Unscoped().Model(&model.PurchaseOrderItem{}).Where("id = ?", id).Updates(updateData).Error
}
