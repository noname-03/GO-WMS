package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
	"time"
)

type InvoiceService struct {
	repo *repository.InvoiceRepository
}

func NewInvoiceService() *InvoiceService {
	return &InvoiceService{
		repo: repository.NewInvoiceRepository(),
	}
}

func (s *InvoiceService) GetAllInvoices() (interface{}, error) {
	return s.repo.GetAllInvoices()
}

// GetFilteredInvoices returns invoices with optional filters
func (s *InvoiceService) GetFilteredInvoices(userID *uint, purchaseOrderID *uint, deliveryOrderID *uint, invoiceDateFrom *time.Time, invoiceDateTo *time.Time, status *string) (interface{}, error) {
	return s.repo.GetFilteredInvoices(userID, purchaseOrderID, deliveryOrderID, invoiceDateFrom, invoiceDateTo, status)
}

func (s *InvoiceService) GetInvoiceByID(id uint) (interface{}, error) {
	invoice, err := s.repo.GetInvoiceByID(id)
	if err != nil {
		return nil, errors.New("invoice not found")
	}
	return invoice, nil
}

// GetInvoiceWithItems returns invoice with all its items
func (s *InvoiceService) GetInvoiceWithItems(id uint) (interface{}, error) {
	invoice, err := s.repo.GetInvoiceWithItems(id)
	if err != nil {
		return nil, errors.New("invoice not found")
	}
	return invoice, nil
}

func (s *InvoiceService) CreateInvoice(
	invoiceNumber string,
	userID uint,
	purchaseOrderID *uint,
	deliveryOrderID *uint,
	invoiceDate time.Time,
	status string,
	totalAmount float64,
	description *string,
	userInsID uint,
) (interface{}, error) {
	// Validate user exists
	exists, err := s.repo.CheckUserExists(userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("user not found")
	}

	// Validate purchase order if provided
	if purchaseOrderID != nil {
		exists, err := s.repo.CheckPurchaseOrderExists(*purchaseOrderID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("purchase order not found")
		}
	}

	// Validate delivery order if provided
	if deliveryOrderID != nil {
		exists, err := s.repo.CheckDeliveryOrderExists(*deliveryOrderID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("delivery order not found")
		}
	}

	// Check if invoice number already exists
	exists, err = s.repo.CheckInvoiceNumberExists(invoiceNumber)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("invoice number already exists")
	}

	invoice := &model.Invoice{
		InvoiceNumber:   invoiceNumber,
		UserID:          userID,
		PurchaseOrderID: purchaseOrderID,
		DeliveryOrderID: deliveryOrderID,
		InvoiceDate:     invoiceDate,
		Status:          status,
		TotalAmount:     totalAmount,
		Description:     description,
		UserIns:         userInsID,
		UserUpdt:        userInsID,
	}

	err = s.repo.CreateInvoice(invoice)
	if err != nil {
		return nil, err
	}

	// Return the created invoice with details
	return s.repo.GetInvoiceByID(invoice.ID)
}

// CreateInvoiceWithItems creates an invoice with multiple items in a transaction
func (s *InvoiceService) CreateInvoiceWithItems(
	invoiceNumber string,
	userID uint,
	purchaseOrderID *uint,
	deliveryOrderID *uint,
	invoiceDate time.Time,
	status string,
	totalAmount float64,
	description *string,
	items []model.InvoiceItem,
	userInsID uint,
) (interface{}, error) {
	// Validate user exists
	exists, err := s.repo.CheckUserExists(userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("user not found")
	}

	// Validate purchase order if provided
	if purchaseOrderID != nil {
		exists, err := s.repo.CheckPurchaseOrderExists(*purchaseOrderID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("purchase order not found")
		}
	}

	// Validate delivery order if provided
	if deliveryOrderID != nil {
		exists, err := s.repo.CheckDeliveryOrderExists(*deliveryOrderID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("delivery order not found")
		}
	}

	// Check if invoice number already exists
	exists, err = s.repo.CheckInvoiceNumberExists(invoiceNumber)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("invoice number already exists")
	}

	// Validate items
	if len(items) == 0 {
		return nil, errors.New("at least one item is required")
	}

	// Set user_ins for all items
	for i := range items {
		items[i].UserIns = userInsID
		items[i].UserUpdt = userInsID
	}

	invoice := &model.Invoice{
		InvoiceNumber:   invoiceNumber,
		UserID:          userID,
		PurchaseOrderID: purchaseOrderID,
		DeliveryOrderID: deliveryOrderID,
		InvoiceDate:     invoiceDate,
		Status:          status,
		TotalAmount:     totalAmount,
		Description:     description,
		UserIns:         userInsID,
		UserUpdt:        userInsID,
	}

	err = s.repo.CreateInvoiceWithItems(invoice, items)
	if err != nil {
		return nil, err
	}

	// Return the created invoice with details
	return s.repo.GetInvoiceByID(invoice.ID)
}

// UpdateInvoiceWithItems updates an invoice and replaces all its items in a transaction
func (s *InvoiceService) UpdateInvoiceWithItems(
	id uint,
	invoiceNumber string,
	userID uint,
	purchaseOrderID *uint,
	deliveryOrderID *uint,
	invoiceDate time.Time,
	status string,
	totalAmount float64,
	description *string,
	items []model.InvoiceItem,
	userUpdtID uint,
) (interface{}, error) {
	// Check if invoice exists
	oldInvoice, err := s.repo.GetInvoiceModelByID(id)
	if err != nil {
		return nil, errors.New("invoice not found")
	}

	// Use existing values if not provided
	if invoiceNumber == "" {
		invoiceNumber = oldInvoice.InvoiceNumber
	}
	if userID == 0 {
		userID = oldInvoice.UserID
	}
	if status == "" {
		status = oldInvoice.Status
	}

	// Validate user exists
	exists, err := s.repo.CheckUserExists(userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("user not found")
	}

	// Validate purchase order if provided
	if purchaseOrderID != nil {
		exists, err := s.repo.CheckPurchaseOrderExists(*purchaseOrderID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("purchase order not found")
		}
	}

	// Validate delivery order if provided
	if deliveryOrderID != nil {
		exists, err := s.repo.CheckDeliveryOrderExists(*deliveryOrderID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("delivery order not found")
		}
	}

	// Validate items
	if len(items) == 0 {
		return nil, errors.New("at least one item is required")
	}

	// Set user_ins and user_updt for all items
	for i := range items {
		items[i].UserIns = userUpdtID
		items[i].UserUpdt = userUpdtID
	}

	updateData := map[string]interface{}{
		"invoice_number":    invoiceNumber,
		"user_id":           userID,
		"purchase_order_id": purchaseOrderID,
		"delivery_order_id": deliveryOrderID,
		"invoice_date":      invoiceDate,
		"status":            status,
		"total_amount":      totalAmount,
		"description":       description,
		"user_updt":         userUpdtID,
	}

	err = s.repo.UpdateInvoiceWithItems(id, updateData, items)
	if err != nil {
		return nil, err
	}

	return s.repo.GetInvoiceByID(id)
}

func (s *InvoiceService) UpdateInvoice(
	id uint,
	invoiceNumber string,
	userID uint,
	purchaseOrderID *uint,
	deliveryOrderID *uint,
	invoiceDate time.Time,
	status string,
	totalAmount float64,
	description *string,
	userUpdtID uint,
) (interface{}, error) {
	// Check if invoice exists
	_, err := s.repo.GetInvoiceModelByID(id)
	if err != nil {
		return nil, errors.New("invoice not found")
	}

	// Validate user exists
	exists, err := s.repo.CheckUserExists(userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("user not found")
	}

	// Validate purchase order if provided
	if purchaseOrderID != nil {
		exists, err := s.repo.CheckPurchaseOrderExists(*purchaseOrderID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("purchase order not found")
		}
	}

	// Validate delivery order if provided
	if deliveryOrderID != nil {
		exists, err := s.repo.CheckDeliveryOrderExists(*deliveryOrderID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("delivery order not found")
		}
	}

	updateData := map[string]interface{}{
		"invoice_number":    invoiceNumber,
		"user_id":           userID,
		"purchase_order_id": purchaseOrderID,
		"delivery_order_id": deliveryOrderID,
		"invoice_date":      invoiceDate,
		"status":            status,
		"total_amount":      totalAmount,
		"description":       description,
		"user_updt":         userUpdtID,
	}

	err = s.repo.UpdateInvoice(id, updateData)
	if err != nil {
		return nil, err
	}

	return s.repo.GetInvoiceByID(id)
}

func (s *InvoiceService) DeleteInvoice(id uint, userID uint) error {
	// Check if invoice exists
	_, err := s.repo.GetInvoiceModelByID(id)
	if err != nil {
		return errors.New("invoice not found")
	}

	return s.repo.DeleteInvoiceWithAudit(id, userID)
}

func (s *InvoiceService) GetDeletedInvoices() (interface{}, error) {
	return s.repo.GetDeletedInvoices()
}

func (s *InvoiceService) RestoreInvoice(id uint, userID uint) (interface{}, error) {
	err := s.repo.RestoreInvoice(id, userID)
	if err != nil {
		return nil, errors.New("invoice not found")
	}

	return s.repo.GetInvoiceByID(id)
}
