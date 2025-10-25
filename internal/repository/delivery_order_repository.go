package repository

import (
	"myapp/database"
	"myapp/internal/model"
	"time"
)

type DeliveryOrderRepository struct{}

// deliveryOrderWithDetailsResponse struct untuk response dengan purchase order details
type deliveryOrderWithDetailsResponse struct {
	ID              uint      `json:"id"`
	DONumber        string    `json:"doNumber"`
	PurchaseOrderID uint      `json:"purchaseOrderId"`
	PONumber        string    `json:"poNumber"`
	DeliveryDate    time.Time `json:"deliveryDate"`
	Status          string    `json:"status"`
	Description     *string   `json:"description"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func NewDeliveryOrderRepository() *DeliveryOrderRepository {
	return &DeliveryOrderRepository{}
}

func (r *DeliveryOrderRepository) GetAllDeliveryOrders() ([]deliveryOrderWithDetailsResponse, error) {
	var orders []deliveryOrderWithDetailsResponse

	result := database.DB.Table("delivery_orders dord").
		Select("dord.id, dord.do_number, dord.purchase_order_id, po.po_number, dord.delivery_date, dord.status, dord.description, dord.created_at, dord.updated_at").
		Joins("LEFT JOIN purchase_orders po ON dord.purchase_order_id = po.id AND po.deleted_at IS NULL").
		Where("dord.deleted_at IS NULL").
		Order("dord.created_at DESC").
		Find(&orders)

	return orders, result.Error
}

// GetFilteredDeliveryOrders returns delivery orders with optional filters
func (r *DeliveryOrderRepository) GetFilteredDeliveryOrders(purchaseOrderID *uint, deliveryDateFrom *time.Time, deliveryDateTo *time.Time, status *string) ([]deliveryOrderWithDetailsResponse, error) {
	var orders []deliveryOrderWithDetailsResponse

	query := database.DB.Table("delivery_orders dord").
		Select("dord.id, dord.do_number, dord.purchase_order_id, po.po_number, dord.delivery_date, dord.status, dord.description, dord.created_at, dord.updated_at").
		Joins("LEFT JOIN purchase_orders po ON dord.purchase_order_id = po.id AND po.deleted_at IS NULL").
		Where("dord.deleted_at IS NULL")

	// Apply filters if provided
	if purchaseOrderID != nil {
		query = query.Where("dord.purchase_order_id = ?", *purchaseOrderID)
	}

	if deliveryDateFrom != nil {
		query = query.Where("dord.delivery_date >= ?", *deliveryDateFrom)
	}

	if deliveryDateTo != nil {
		query = query.Where("dord.delivery_date <= ?", *deliveryDateTo)
	}

	if status != nil && *status != "" {
		query = query.Where("dord.status = ?", *status)
	}

	result := query.Order("dord.created_at DESC").Find(&orders)

	return orders, result.Error
}

func (r *DeliveryOrderRepository) GetDeliveryOrderByID(id uint) (deliveryOrderWithDetailsResponse, error) {
	var order deliveryOrderWithDetailsResponse

	result := database.DB.Table("delivery_orders dord").
		Select("dord.id, dord.do_number, dord.purchase_order_id, po.po_number, dord.delivery_date, dord.status, dord.description, dord.created_at, dord.updated_at").
		Joins("INNER JOIN purchase_orders po ON dord.purchase_order_id = po.id AND po.deleted_at IS NULL").
		Where("dord.id = ? AND dord.deleted_at IS NULL", id).
		First(&order)

	return order, result.Error
}

// GetDeliveryOrderModelByID returns model.DeliveryOrder for service operations
func (r *DeliveryOrderRepository) GetDeliveryOrderModelByID(id uint) (model.DeliveryOrder, error) {
	var order model.DeliveryOrder
	result := database.DB.Where("id = ?", id).First(&order)
	return order, result.Error
}

// deliveryOrderWithItemsResponse struct untuk response dengan items
type deliveryOrderWithItemsResponse struct {
	ID              uint                                `json:"id"`
	DONumber        string                              `json:"doNumber"`
	PurchaseOrderID uint                                `json:"purchaseOrderId"`
	PONumber        string                              `json:"poNumber"`
	DeliveryDate    time.Time                           `json:"deliveryDate"`
	Status          string                              `json:"status"`
	Description     *string                             `json:"description"`
	CreatedAt       time.Time                           `json:"createdAt"`
	UpdatedAt       time.Time                           `json:"updatedAt"`
	Items           []doItemDetailsForWithItemsResponse `gorm:"-" json:"items"`
}

type doItemDetailsForWithItemsResponse struct {
	ID           uint     `json:"id"`
	ProductID    uint     `json:"productId"`
	ProductName  string   `json:"productName"`
	QtyDelivered *float64 `json:"qtyDelivered"`
	Description  *string  `json:"description"`
}

// GetDeliveryOrderWithItems returns delivery order with all its items
func (r *DeliveryOrderRepository) GetDeliveryOrderWithItems(id uint) (deliveryOrderWithItemsResponse, error) {
	var order deliveryOrderWithItemsResponse

	// Get delivery order
	result := database.DB.Table("delivery_orders dord").
		Select("dord.id, dord.do_number, dord.purchase_order_id, po.po_number, dord.delivery_date, dord.status, dord.description, dord.created_at, dord.updated_at").
		Joins("INNER JOIN purchase_orders po ON dord.purchase_order_id = po.id AND po.deleted_at IS NULL").
		Where("dord.id = ? AND dord.deleted_at IS NULL", id).
		First(&order)

	if result.Error != nil {
		return order, result.Error
	}

	// Get items
	var items []doItemDetailsForWithItemsResponse
	result = database.DB.Table("delivery_order_items doi").
		Select("doi.id, doi.product_id, p.name as product_name, doi.qty_delivered, doi.description").
		Joins("INNER JOIN products p ON doi.product_id = p.id AND p.deleted_at IS NULL").
		Where("doi.delivery_order_id = ? AND doi.deleted_at IS NULL", id).
		Order("doi.id ASC").
		Find(&items)

	if result.Error != nil {
		return order, result.Error
	}

	order.Items = items
	return order, nil
}

func (r *DeliveryOrderRepository) CreateDeliveryOrder(order *model.DeliveryOrder) error {
	return database.DB.Create(order).Error
}

// CreateDeliveryOrderWithItems creates a delivery order and its items in a transaction
func (r *DeliveryOrderRepository) CreateDeliveryOrderWithItems(order *model.DeliveryOrder, items []model.DeliveryOrderItem) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Create delivery order
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set delivery_order_id for all items
	for i := range items {
		items[i].DeliveryOrderID = order.ID
	}

	// Create all items
	if len(items) > 0 {
		if err := tx.Create(&items).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// UpdateDeliveryOrderWithItems updates a delivery order and replaces all its items in a transaction
func (r *DeliveryOrderRepository) UpdateDeliveryOrderWithItems(id uint, order map[string]interface{}, items []model.DeliveryOrderItem) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Update delivery order
	if err := tx.Model(&model.DeliveryOrder{}).Where("id = ?", id).Updates(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete existing items (soft delete)
	if err := tx.Where("delivery_order_id = ?", id).Delete(&model.DeliveryOrderItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set delivery_order_id for all new items
	for i := range items {
		items[i].DeliveryOrderID = id
	}

	// Create new items
	if len(items) > 0 {
		if err := tx.Create(&items).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *DeliveryOrderRepository) UpdateDeliveryOrder(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.DeliveryOrder{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *DeliveryOrderRepository) DeleteDeliveryOrderWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the order
	updateData := map[string]interface{}{
		"user_updt": userID,
	}

	// Update the audit field first
	err := database.DB.Model(&model.DeliveryOrder{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.DeliveryOrder{}, id).Error
}

func (r *DeliveryOrderRepository) CheckPurchaseOrderExists(purchaseOrderID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.PurchaseOrder{}).Where("id = ?", purchaseOrderID).Count(&count)
	return count > 0, result.Error
}

func (r *DeliveryOrderRepository) CheckDONumberExists(doNumber string) (bool, error) {
	var count int64
	result := database.DB.Model(&model.DeliveryOrder{}).Where("do_number = ?", doNumber).Count(&count)
	return count > 0, result.Error
}

// GetDeletedDeliveryOrders returns all soft deleted delivery orders
func (r *DeliveryOrderRepository) GetDeletedDeliveryOrders() ([]deliveryOrderWithDetailsResponse, error) {
	var orders []deliveryOrderWithDetailsResponse

	result := database.DB.Table("delivery_orders dord").
		Select("dord.id, dord.do_number, dord.purchase_order_id, po.po_number, dord.delivery_date, dord.status, dord.description, dord.created_at, dord.updated_at").
		Joins("LEFT JOIN purchase_orders po ON dord.purchase_order_id = po.id").
		Where("dord.deleted_at IS NOT NULL").
		Order("dord.deleted_at DESC").
		Find(&orders)

	return orders, result.Error
}

// RestoreDeliveryOrder restores a soft deleted delivery order
func (r *DeliveryOrderRepository) RestoreDeliveryOrder(id uint, userID uint) error {
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"deleted_at": nil,
	}
	return database.DB.Unscoped().Model(&model.DeliveryOrder{}).Where("id = ?", id).Updates(updateData).Error
}
