package repository

import (
	"myapp/database"
	"myapp/internal/model"
	"time"
)

type InvoiceRepository struct{}

// invoiceWithDetailsResponse struct untuk response dengan user (reseller) name, PO, DO details
type invoiceWithDetailsResponse struct {
	ID              uint      `json:"id"`
	InvoiceNumber   string    `json:"invoiceNumber"`
	UserID          uint      `json:"userId"`
	UserName        string    `json:"userName"`
	PurchaseOrderID *uint     `json:"purchaseOrderId"`
	PONumber        *string   `json:"poNumber"`
	DeliveryOrderID *uint     `json:"deliveryOrderId"`
	DONumber        *string   `json:"doNumber"`
	InvoiceDate     time.Time `json:"invoiceDate"`
	Status          string    `json:"status"`
	TotalAmount     float64   `json:"totalAmount"`
	Description     *string   `json:"description"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func NewInvoiceRepository() *InvoiceRepository {
	return &InvoiceRepository{}
}

func (r *InvoiceRepository) GetAllInvoices() ([]invoiceWithDetailsResponse, error) {
	var invoices []invoiceWithDetailsResponse

	result := database.DB.Table("invoices inv").
		Select("inv.id, inv.invoice_number, inv.user_id, u.name as user_name, inv.purchase_order_id, po.po_number, inv.delivery_order_id, dord.do_number, inv.invoice_date, inv.status, inv.total_amount, inv.description, inv.created_at, inv.updated_at").
		Joins("LEFT JOIN users u ON inv.user_id = u.id AND u.deleted_at IS NULL").
		Joins("LEFT JOIN purchase_orders po ON inv.purchase_order_id = po.id AND po.deleted_at IS NULL").
		Joins("LEFT JOIN delivery_orders dord ON inv.delivery_order_id = dord.id AND dord.deleted_at IS NULL").
		Where("inv.deleted_at IS NULL").
		Order("inv.created_at DESC").
		Find(&invoices)

	return invoices, result.Error
}

// GetFilteredInvoices returns invoices with optional filters
func (r *InvoiceRepository) GetFilteredInvoices(userID *uint, purchaseOrderID *uint, deliveryOrderID *uint, invoiceDateFrom *time.Time, invoiceDateTo *time.Time, status *string) ([]invoiceWithDetailsResponse, error) {
	var invoices []invoiceWithDetailsResponse

	query := database.DB.Table("invoices inv").
		Select("inv.id, inv.invoice_number, inv.user_id, u.name as user_name, inv.purchase_order_id, po.po_number, inv.delivery_order_id, dord.do_number, inv.invoice_date, inv.status, inv.total_amount, inv.description, inv.created_at, inv.updated_at").
		Joins("LEFT JOIN users u ON inv.user_id = u.id AND u.deleted_at IS NULL").
		Joins("LEFT JOIN purchase_orders po ON inv.purchase_order_id = po.id AND po.deleted_at IS NULL").
		Joins("LEFT JOIN delivery_orders dord ON inv.delivery_order_id = dord.id AND dord.deleted_at IS NULL").
		Where("inv.deleted_at IS NULL")

	// Apply filters if provided
	if userID != nil {
		query = query.Where("inv.user_id = ?", *userID)
	}

	if purchaseOrderID != nil {
		query = query.Where("inv.purchase_order_id = ?", *purchaseOrderID)
	}

	if deliveryOrderID != nil {
		query = query.Where("inv.delivery_order_id = ?", *deliveryOrderID)
	}

	if invoiceDateFrom != nil {
		query = query.Where("inv.invoice_date >= ?", *invoiceDateFrom)
	}

	if invoiceDateTo != nil {
		query = query.Where("inv.invoice_date <= ?", *invoiceDateTo)
	}

	if status != nil && *status != "" {
		query = query.Where("inv.status = ?", *status)
	}

	result := query.Order("inv.created_at DESC").Find(&invoices)

	return invoices, result.Error
}

func (r *InvoiceRepository) GetInvoiceByID(id uint) (invoiceWithDetailsResponse, error) {
	var invoice invoiceWithDetailsResponse

	result := database.DB.Table("invoices inv").
		Select("inv.id, inv.invoice_number, inv.user_id, u.name as user_name, inv.purchase_order_id, po.po_number, inv.delivery_order_id, dord.do_number, inv.invoice_date, inv.status, inv.total_amount, inv.description, inv.created_at, inv.updated_at").
		Joins("INNER JOIN users u ON inv.user_id = u.id AND u.deleted_at IS NULL").
		Joins("LEFT JOIN purchase_orders po ON inv.purchase_order_id = po.id AND po.deleted_at IS NULL").
		Joins("LEFT JOIN delivery_orders dord ON inv.delivery_order_id = dord.id AND dord.deleted_at IS NULL").
		Where("inv.id = ? AND inv.deleted_at IS NULL", id).
		First(&invoice)

	return invoice, result.Error
}

// GetInvoiceModelByID returns model.Invoice for service operations
func (r *InvoiceRepository) GetInvoiceModelByID(id uint) (model.Invoice, error) {
	var invoice model.Invoice
	result := database.DB.Where("id = ?", id).First(&invoice)
	return invoice, result.Error
}

func (r *InvoiceRepository) CreateInvoice(invoice *model.Invoice) error {
	return database.DB.Create(invoice).Error
}

func (r *InvoiceRepository) UpdateInvoice(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.Invoice{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *InvoiceRepository) DeleteInvoiceWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the invoice
	updateData := map[string]interface{}{
		"user_updt": userID,
	}

	// Update the audit field first
	err := database.DB.Model(&model.Invoice{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.Invoice{}, id).Error
}

func (r *InvoiceRepository) CheckUserExists(userID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.User{}).Where("id = ?", userID).Count(&count)
	return count > 0, result.Error
}

func (r *InvoiceRepository) CheckPurchaseOrderExists(purchaseOrderID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.PurchaseOrder{}).Where("id = ?", purchaseOrderID).Count(&count)
	return count > 0, result.Error
}

func (r *InvoiceRepository) CheckDeliveryOrderExists(deliveryOrderID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.DeliveryOrder{}).Where("id = ?", deliveryOrderID).Count(&count)
	return count > 0, result.Error
}

func (r *InvoiceRepository) CheckInvoiceNumberExists(invoiceNumber string) (bool, error) {
	var count int64
	result := database.DB.Model(&model.Invoice{}).Where("invoice_number = ?", invoiceNumber).Count(&count)
	return count > 0, result.Error
}

// GetDeletedInvoices returns all soft deleted invoices
func (r *InvoiceRepository) GetDeletedInvoices() ([]invoiceWithDetailsResponse, error) {
	var invoices []invoiceWithDetailsResponse

	result := database.DB.Table("invoices inv").
		Select("inv.id, inv.invoice_number, inv.user_id, u.name as user_name, inv.purchase_order_id, po.po_number, inv.delivery_order_id, dord.do_number, inv.invoice_date, inv.status, inv.total_amount, inv.description, inv.created_at, inv.updated_at").
		Joins("LEFT JOIN users u ON inv.user_id = u.id").
		Joins("LEFT JOIN purchase_orders po ON inv.purchase_order_id = po.id").
		Joins("LEFT JOIN delivery_orders dord ON inv.delivery_order_id = dord.id").
		Where("inv.deleted_at IS NOT NULL").
		Order("inv.deleted_at DESC").
		Find(&invoices)

	return invoices, result.Error
}

// RestoreInvoice restores a soft deleted invoice
func (r *InvoiceRepository) RestoreInvoice(id uint, userID uint) error {
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"deleted_at": nil,
	}
	return database.DB.Unscoped().Model(&model.Invoice{}).Where("id = ?", id).Updates(updateData).Error
}
