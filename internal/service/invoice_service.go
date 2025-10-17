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

func (s *InvoiceService) GetInvoiceByID(id uint) (interface{}, error) {
	invoice, err := s.repo.GetInvoiceByID(id)
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
