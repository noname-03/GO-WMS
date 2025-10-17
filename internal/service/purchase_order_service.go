package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
	"time"
)

type PurchaseOrderService struct {
	poRepo *repository.PurchaseOrderRepository
}

func NewPurchaseOrderService() *PurchaseOrderService {
	return &PurchaseOrderService{
		poRepo: repository.NewPurchaseOrderRepository(),
	}
}

// Business logic methods
func (s *PurchaseOrderService) GetAllPurchaseOrders() (interface{}, error) {
	return s.poRepo.GetAllPurchaseOrders()
}

// GetFilteredPurchaseOrders returns purchase orders with optional filters
func (s *PurchaseOrderService) GetFilteredPurchaseOrders(userID *uint, orderDateFrom *time.Time, orderDateTo *time.Time, status *string) (interface{}, error) {
	return s.poRepo.GetFilteredPurchaseOrders(userID, orderDateFrom, orderDateTo, status)
}

func (s *PurchaseOrderService) GetPurchaseOrderByID(id uint) (interface{}, error) {
	order, err := s.poRepo.GetPurchaseOrderByID(id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *PurchaseOrderService) CreatePurchaseOrder(poNumber string, orderUserID uint, orderDate time.Time, status string, totalAmount *float64, description *string, userID uint) (interface{}, error) {
	if poNumber == "" {
		return nil, errors.New("PO number is required")
	}

	if orderUserID == 0 {
		return nil, errors.New("user ID is required")
	}

	if orderDate.IsZero() {
		return nil, errors.New("order date is required")
	}

	if status == "" {
		status = "draft" // Default status
	}

	// Validate status
	validStatuses := map[string]bool{
		"draft":     true,
		"submitted": true,
		"approved":  true,
		"received":  true,
		"closed":    true,
	}
	if !validStatuses[status] {
		return nil, errors.New("invalid status: must be one of draft, submitted, approved, received, closed")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if user exists
	userExists, err := s.poRepo.CheckUserExists(orderUserID)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, errors.New("user not found")
	}

	// Check if PO number already exists
	poExists, err := s.poRepo.CheckPONumberExists(poNumber)
	if err != nil {
		return nil, err
	}
	if poExists {
		return nil, errors.New("PO number already exists")
	}

	// Set default total amount if not provided
	if totalAmount == nil {
		zero := 0.0
		totalAmount = &zero
	}

	order := &model.PurchaseOrder{
		PONumber:    poNumber,
		UserID:      orderUserID,
		OrderDate:   orderDate,
		Status:      status,
		TotalAmount: totalAmount,
		Description: description,
		UserIns:     &userID,
	}

	err = s.poRepo.CreatePurchaseOrder(order)
	if err != nil {
		return nil, err
	}

	// Fetch the created order with relationships
	createdOrder, err := s.poRepo.GetPurchaseOrderByID(order.ID)
	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}

func (s *PurchaseOrderService) UpdatePurchaseOrder(id uint, poNumber string, orderUserID uint, orderDate time.Time, status string, totalAmount *float64, description *string, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid purchase order ID")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if purchase order exists
	oldOrder, err := s.poRepo.GetPurchaseOrderModelByID(id)
	if err != nil {
		return nil, errors.New("purchase order not found")
	}

	// Validate status if provided
	if status != "" {
		validStatuses := map[string]bool{
			"draft":     true,
			"submitted": true,
			"approved":  true,
			"received":  true,
			"closed":    true,
		}
		if !validStatuses[status] {
			return nil, errors.New("invalid status: must be one of draft, submitted, approved, received, closed")
		}
	}

	// If user ID is being changed, check if new user exists
	if orderUserID != 0 && orderUserID != oldOrder.UserID {
		userExists, err := s.poRepo.CheckUserExists(orderUserID)
		if err != nil {
			return nil, err
		}
		if !userExists {
			return nil, errors.New("user not found")
		}
	}

	// Prepare update data with audit trail
	updateData := make(map[string]interface{})
	if poNumber != "" && poNumber != oldOrder.PONumber {
		// Check if new PO number already exists
		poExists, err := s.poRepo.CheckPONumberExists(poNumber)
		if err != nil {
			return nil, err
		}
		if poExists {
			return nil, errors.New("PO number already exists")
		}
		updateData["po_number"] = poNumber
	}
	if orderUserID != 0 && orderUserID != oldOrder.UserID {
		updateData["user_id"] = orderUserID
	}
	if !orderDate.IsZero() {
		updateData["order_date"] = orderDate
	}
	if status != "" {
		updateData["status"] = status
	}
	if totalAmount != nil {
		updateData["total_amount"] = totalAmount
	}
	if description != nil {
		updateData["description"] = description
	}
	// Always set the user who updated
	updateData["user_updt"] = userID

	err = s.poRepo.UpdatePurchaseOrder(id, updateData)
	if err != nil {
		return nil, err
	}

	updatedOrder, err := s.poRepo.GetPurchaseOrderByID(id)
	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

func (s *PurchaseOrderService) DeletePurchaseOrder(id uint, userID uint) error {
	if id == 0 {
		return errors.New("invalid purchase order ID")
	}

	if userID == 0 {
		return errors.New("user ID is required for audit trail")
	}

	// Check if purchase order exists
	_, err := s.poRepo.GetPurchaseOrderModelByID(id)
	if err != nil {
		return errors.New("purchase order not found")
	}

	return s.poRepo.DeletePurchaseOrderWithAudit(id, userID)
}

// GetDeletedPurchaseOrders returns all soft deleted purchase orders
func (s *PurchaseOrderService) GetDeletedPurchaseOrders() (interface{}, error) {
	return s.poRepo.GetDeletedPurchaseOrders()
}

// RestorePurchaseOrder restores a soft deleted purchase order
func (s *PurchaseOrderService) RestorePurchaseOrder(id uint, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid purchase order ID")
	}
	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	err := s.poRepo.RestorePurchaseOrder(id, userID)
	if err != nil {
		return nil, err
	}

	restoredOrder, err := s.poRepo.GetPurchaseOrderByID(id)
	if err != nil {
		return nil, err
	}
	return restoredOrder, nil
}
