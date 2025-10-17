package repository

import (
	"myapp/database"
	"myapp/internal/model"
	"time"
)

type InvoiceItemRepository struct{}

// invoiceItemWithDetailsResponse struct untuk response dengan invoice number dan product name
type invoiceItemWithDetailsResponse struct {
	ID            uint      `json:"id"`
	InvoiceID     uint      `json:"invoiceId"`
	InvoiceNumber string    `json:"invoiceNumber"`
	ProductID     uint      `json:"productId"`
	ProductName   string    `json:"productName"`
	QtyInvoiced   float64   `json:"qtyInvoiced"`
	UnitPrice     float64   `json:"unitPrice"`
	TotalPrice    float64   `json:"totalPrice"`
	Description   *string   `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func NewInvoiceItemRepository() *InvoiceItemRepository {
	return &InvoiceItemRepository{}
}

func (r *InvoiceItemRepository) GetAllInvoiceItems() ([]invoiceItemWithDetailsResponse, error) {
	var items []invoiceItemWithDetailsResponse

	result := database.DB.Table("invoice_items ii").
		Select("ii.id, ii.invoice_id, inv.invoice_number, ii.product_id, p.name as product_name, ii.qty_invoiced, ii.unit_price, ii.total_price, ii.description, ii.created_at, ii.updated_at").
		Joins("INNER JOIN invoices inv ON ii.invoice_id = inv.id AND inv.deleted_at IS NULL").
		Joins("INNER JOIN products p ON ii.product_id = p.id AND p.deleted_at IS NULL").
		Where("ii.deleted_at IS NULL").
		Order("ii.created_at DESC").
		Find(&items)

	return items, result.Error
}

func (r *InvoiceItemRepository) GetInvoiceItemByID(id uint) (invoiceItemWithDetailsResponse, error) {
	var item invoiceItemWithDetailsResponse

	result := database.DB.Table("invoice_items ii").
		Select("ii.id, ii.invoice_id, inv.invoice_number, ii.product_id, p.name as product_name, ii.qty_invoiced, ii.unit_price, ii.total_price, ii.description, ii.created_at, ii.updated_at").
		Joins("INNER JOIN invoices inv ON ii.invoice_id = inv.id AND inv.deleted_at IS NULL").
		Joins("INNER JOIN products p ON ii.product_id = p.id AND p.deleted_at IS NULL").
		Where("ii.id = ? AND ii.deleted_at IS NULL", id).
		First(&item)

	return item, result.Error
}

func (r *InvoiceItemRepository) GetInvoiceItemsByInvoice(invoiceID uint) ([]invoiceItemWithDetailsResponse, error) {
	var items []invoiceItemWithDetailsResponse

	result := database.DB.Table("invoice_items ii").
		Select("ii.id, ii.invoice_id, inv.invoice_number, ii.product_id, p.name as product_name, ii.qty_invoiced, ii.unit_price, ii.total_price, ii.description, ii.created_at, ii.updated_at").
		Joins("INNER JOIN invoices inv ON ii.invoice_id = inv.id AND inv.deleted_at IS NULL").
		Joins("INNER JOIN products p ON ii.product_id = p.id AND p.deleted_at IS NULL").
		Where("ii.invoice_id = ? AND ii.deleted_at IS NULL", invoiceID).
		Order("ii.created_at DESC").
		Find(&items)

	return items, result.Error
}

// GetInvoiceItemModelByID returns model.InvoiceItem for service operations
func (r *InvoiceItemRepository) GetInvoiceItemModelByID(id uint) (model.InvoiceItem, error) {
	var item model.InvoiceItem
	result := database.DB.Where("id = ?", id).First(&item)
	return item, result.Error
}

func (r *InvoiceItemRepository) CreateInvoiceItem(item *model.InvoiceItem) error {
	return database.DB.Create(item).Error
}

func (r *InvoiceItemRepository) UpdateInvoiceItem(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.InvoiceItem{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *InvoiceItemRepository) DeleteInvoiceItemWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the item
	updateData := map[string]interface{}{
		"user_updt": userID,
	}

	// Update the audit field first
	err := database.DB.Model(&model.InvoiceItem{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.InvoiceItem{}, id).Error
}

func (r *InvoiceItemRepository) CheckInvoiceExists(invoiceID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.Invoice{}).Where("id = ?", invoiceID).Count(&count)
	return count > 0, result.Error
}

func (r *InvoiceItemRepository) CheckProductExists(productID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.Product{}).Where("id = ?", productID).Count(&count)
	return count > 0, result.Error
}

// GetDeletedInvoiceItems returns all soft deleted invoice items
func (r *InvoiceItemRepository) GetDeletedInvoiceItems() ([]invoiceItemWithDetailsResponse, error) {
	var items []invoiceItemWithDetailsResponse

	result := database.DB.Table("invoice_items ii").
		Select("ii.id, ii.invoice_id, inv.invoice_number, ii.product_id, p.name as product_name, ii.qty_invoiced, ii.unit_price, ii.total_price, ii.description, ii.created_at, ii.updated_at").
		Joins("LEFT JOIN invoices inv ON ii.invoice_id = inv.id").
		Joins("LEFT JOIN products p ON ii.product_id = p.id").
		Where("ii.deleted_at IS NOT NULL").
		Order("ii.deleted_at DESC").
		Find(&items)

	return items, result.Error
}

// RestoreInvoiceItem restores a soft deleted invoice item
func (r *InvoiceItemRepository) RestoreInvoiceItem(id uint, userID uint) error {
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"deleted_at": nil,
	}
	return database.DB.Unscoped().Model(&model.InvoiceItem{}).Where("id = ?", id).Updates(updateData).Error
}
