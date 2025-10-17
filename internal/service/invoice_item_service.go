package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
)

type InvoiceItemService struct {
	repo *repository.InvoiceItemRepository
}

func NewInvoiceItemService() *InvoiceItemService {
	return &InvoiceItemService{
		repo: repository.NewInvoiceItemRepository(),
	}
}

func (s *InvoiceItemService) GetAllInvoiceItems() (interface{}, error) {
	return s.repo.GetAllInvoiceItems()
}

func (s *InvoiceItemService) GetInvoiceItemByID(id uint) (interface{}, error) {
	item, err := s.repo.GetInvoiceItemByID(id)
	if err != nil {
		return nil, errors.New("invoice item not found")
	}
	return item, nil
}

func (s *InvoiceItemService) GetInvoiceItemsByInvoice(invoiceID uint) (interface{}, error) {
	// Check if invoice exists
	exists, err := s.repo.CheckInvoiceExists(invoiceID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("invoice not found")
	}

	return s.repo.GetInvoiceItemsByInvoice(invoiceID)
}

func (s *InvoiceItemService) CreateInvoiceItem(
	invoiceID uint,
	productID uint,
	qtyInvoiced float64,
	unitPrice float64,
	totalPrice float64,
	description *string,
	userID uint,
) (interface{}, error) {
	// Validate invoice exists
	exists, err := s.repo.CheckInvoiceExists(invoiceID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("invoice not found")
	}

	// Validate product exists
	exists, err = s.repo.CheckProductExists(productID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("product not found")
	}

	item := &model.InvoiceItem{
		InvoiceID:   invoiceID,
		ProductID:   productID,
		QtyInvoiced: qtyInvoiced,
		UnitPrice:   unitPrice,
		TotalPrice:  totalPrice,
		Description: description,
		UserIns:     userID,
		UserUpdt:    userID,
	}

	err = s.repo.CreateInvoiceItem(item)
	if err != nil {
		return nil, err
	}

	// Return the created item with details
	return s.repo.GetInvoiceItemByID(item.ID)
}

func (s *InvoiceItemService) UpdateInvoiceItem(
	id uint,
	invoiceID uint,
	productID uint,
	qtyInvoiced float64,
	unitPrice float64,
	totalPrice float64,
	description *string,
	userID uint,
) (interface{}, error) {
	// Check if item exists
	_, err := s.repo.GetInvoiceItemModelByID(id)
	if err != nil {
		return nil, errors.New("invoice item not found")
	}

	// Validate invoice exists
	exists, err := s.repo.CheckInvoiceExists(invoiceID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("invoice not found")
	}

	// Validate product exists
	exists, err = s.repo.CheckProductExists(productID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("product not found")
	}

	updateData := map[string]interface{}{
		"invoice_id":   invoiceID,
		"product_id":   productID,
		"qty_invoiced": qtyInvoiced,
		"unit_price":   unitPrice,
		"total_price":  totalPrice,
		"description":  description,
		"user_updt":    userID,
	}

	err = s.repo.UpdateInvoiceItem(id, updateData)
	if err != nil {
		return nil, err
	}

	return s.repo.GetInvoiceItemByID(id)
}

func (s *InvoiceItemService) DeleteInvoiceItem(id uint, userID uint) error {
	// Check if item exists
	_, err := s.repo.GetInvoiceItemModelByID(id)
	if err != nil {
		return errors.New("invoice item not found")
	}

	return s.repo.DeleteInvoiceItemWithAudit(id, userID)
}

func (s *InvoiceItemService) GetDeletedInvoiceItems() (interface{}, error) {
	return s.repo.GetDeletedInvoiceItems()
}

func (s *InvoiceItemService) RestoreInvoiceItem(id uint, userID uint) (interface{}, error) {
	err := s.repo.RestoreInvoiceItem(id, userID)
	if err != nil {
		return nil, errors.New("invoice item not found")
	}

	return s.repo.GetInvoiceItemByID(id)
}
