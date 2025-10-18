package repository

import (
	"myapp/database"
	"myapp/internal/model"
	"time"
)

type PurchaseOrderRepository struct{}

// purchaseOrderWithDetailsResponse struct untuk response dengan user name dan items
type purchaseOrderWithDetailsResponse struct {
	ID          uint      `json:"id"`
	PONumber    string    `json:"poNumber"`
	UserID      uint      `json:"userId"`
	UserName    string    `json:"userName"`
	OrderDate   time.Time `json:"orderDate"`
	Status      string    `json:"status"`
	TotalAmount *float64  `json:"totalAmount"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewPurchaseOrderRepository() *PurchaseOrderRepository {
	return &PurchaseOrderRepository{}
}

func (r *PurchaseOrderRepository) GetAllPurchaseOrders() ([]purchaseOrderWithDetailsResponse, error) {
	var orders []purchaseOrderWithDetailsResponse

	result := database.DB.Table("purchase_orders po").
		Select("po.id, po.po_number, po.user_id, u.name as user_name, po.order_date, po.status, po.total_amount, po.description, po.created_at, po.updated_at").
		Joins("LEFT JOIN users u ON po.user_id = u.id AND u.deleted_at IS NULL").
		Where("po.deleted_at IS NULL").
		Order("po.created_at DESC").
		Find(&orders)

	return orders, result.Error
}

// GetFilteredPurchaseOrders returns purchase orders with optional filters
func (r *PurchaseOrderRepository) GetFilteredPurchaseOrders(userID *uint, orderDateFrom *time.Time, orderDateTo *time.Time, status *string) ([]purchaseOrderWithDetailsResponse, error) {
	var orders []purchaseOrderWithDetailsResponse

	query := database.DB.Table("purchase_orders po").
		Select("po.id, po.po_number, po.user_id, u.name as user_name, po.order_date, po.status, po.total_amount, po.description, po.created_at, po.updated_at").
		Joins("LEFT JOIN users u ON po.user_id = u.id AND u.deleted_at IS NULL").
		Where("po.deleted_at IS NULL")

	// Apply filters if provided
	if userID != nil {
		query = query.Where("po.user_id = ?", *userID)
	}

	if orderDateFrom != nil {
		query = query.Where("po.order_date >= ?", *orderDateFrom)
	}

	if orderDateTo != nil {
		query = query.Where("po.order_date <= ?", *orderDateTo)
	}

	if status != nil && *status != "" {
		query = query.Where("po.status = ?", *status)
	}

	result := query.Order("po.created_at DESC").Find(&orders)

	return orders, result.Error
}

func (r *PurchaseOrderRepository) GetPurchaseOrderByID(id uint) (purchaseOrderWithDetailsResponse, error) {
	var order purchaseOrderWithDetailsResponse

	result := database.DB.Table("purchase_orders po").
		Select("po.id, po.po_number, po.user_id, u.name as user_name, po.order_date, po.status, po.total_amount, po.description, po.created_at, po.updated_at").
		Joins("INNER JOIN users u ON po.user_id = u.id AND u.deleted_at IS NULL").
		Where("po.id = ? AND po.deleted_at IS NULL", id).
		First(&order)

	return order, result.Error
}

// GetPurchaseOrderModelByID returns model.PurchaseOrder for service operations
func (r *PurchaseOrderRepository) GetPurchaseOrderModelByID(id uint) (model.PurchaseOrder, error) {
	var order model.PurchaseOrder
	result := database.DB.Where("id = ?", id).First(&order)
	return order, result.Error
}

// purchaseOrderWithItemsResponse struct untuk response dengan items
type purchaseOrderWithItemsResponse struct {
	ID          uint                                `json:"id"`
	PONumber    string                              `json:"poNumber"`
	UserID      uint                                `json:"userId"`
	UserName    string                              `json:"userName"`
	OrderDate   time.Time                           `json:"orderDate"`
	Status      string                              `json:"status"`
	TotalAmount *float64                            `json:"totalAmount"`
	Description *string                             `json:"description"`
	CreatedAt   time.Time                           `json:"createdAt"`
	UpdatedAt   time.Time                           `json:"updatedAt"`
	Items       []poItemDetailsForWithItemsResponse `gorm:"-" json:"items"`
}

type poItemDetailsForWithItemsResponse struct {
	ID          uint     `json:"id"`
	ProductID   uint     `json:"productId"`
	ProductName string   `json:"productName"`
	QtyOrdered  *float64 `json:"qtyOrdered"`
	UnitPrice   *float64 `json:"unitPrice"`
	Discount    *float64 `json:"discount"`
	TotalPrice  *float64 `json:"totalPrice"`
	Description *string  `json:"description"`
}

// GetPurchaseOrderWithItems returns purchase order with all its items
func (r *PurchaseOrderRepository) GetPurchaseOrderWithItems(id uint) (purchaseOrderWithItemsResponse, error) {
	var order purchaseOrderWithItemsResponse

	// Get purchase order
	result := database.DB.Table("purchase_orders po").
		Select("po.id, po.po_number, po.user_id, u.name as user_name, po.order_date, po.status, po.total_amount, po.description, po.created_at, po.updated_at").
		Joins("INNER JOIN users u ON po.user_id = u.id AND u.deleted_at IS NULL").
		Where("po.id = ? AND po.deleted_at IS NULL", id).
		First(&order)

	if result.Error != nil {
		return order, result.Error
	}

	// Get items
	var items []poItemDetailsForWithItemsResponse
	result = database.DB.Table("purchase_order_items poi").
		Select("poi.id, poi.product_id, p.name as product_name, poi.qty_ordered, poi.unit_price, poi.discount, poi.total_price, poi.description").
		Joins("INNER JOIN products p ON poi.product_id = p.id AND p.deleted_at IS NULL").
		Where("poi.purchase_order_id = ? AND poi.deleted_at IS NULL", id).
		Order("poi.id ASC").
		Find(&items)

	if result.Error != nil {
		return order, result.Error
	}

	order.Items = items
	return order, nil
}

func (r *PurchaseOrderRepository) CreatePurchaseOrder(order *model.PurchaseOrder) error {
	return database.DB.Create(order).Error
}

// CreatePurchaseOrderWithItems creates a purchase order and its items in a transaction
func (r *PurchaseOrderRepository) CreatePurchaseOrderWithItems(order *model.PurchaseOrder, items []model.PurchaseOrderItem) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Create purchase order
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set purchase_order_id for all items
	for i := range items {
		items[i].PurchaseOrderID = order.ID
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

func (r *PurchaseOrderRepository) UpdatePurchaseOrder(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.PurchaseOrder{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *PurchaseOrderRepository) DeletePurchaseOrderWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the order
	updateData := map[string]interface{}{
		"user_updt": userID,
	}

	// Update the audit field first
	err := database.DB.Model(&model.PurchaseOrder{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.PurchaseOrder{}, id).Error
}

func (r *PurchaseOrderRepository) CheckUserExists(userID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.User{}).Where("id = ?", userID).Count(&count)
	return count > 0, result.Error
}

func (r *PurchaseOrderRepository) CheckPONumberExists(poNumber string) (bool, error) {
	var count int64
	result := database.DB.Model(&model.PurchaseOrder{}).Where("po_number = ?", poNumber).Count(&count)
	return count > 0, result.Error
}

// GetDeletedPurchaseOrders returns all soft deleted purchase orders
func (r *PurchaseOrderRepository) GetDeletedPurchaseOrders() ([]purchaseOrderWithDetailsResponse, error) {
	var orders []purchaseOrderWithDetailsResponse

	result := database.DB.Table("purchase_orders po").
		Select("po.id, po.po_number, po.user_id, u.name as user_name, po.order_date, po.status, po.total_amount, po.description, po.created_at, po.updated_at").
		Joins("LEFT JOIN users u ON po.user_id = u.id").
		Where("po.deleted_at IS NOT NULL").
		Order("po.deleted_at DESC").
		Find(&orders)

	return orders, result.Error
}

// RestorePurchaseOrder restores a soft deleted purchase order
func (r *PurchaseOrderRepository) RestorePurchaseOrder(id uint, userID uint) error {
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"deleted_at": nil,
	}
	return database.DB.Unscoped().Model(&model.PurchaseOrder{}).Where("id = ?", id).Updates(updateData).Error
}
