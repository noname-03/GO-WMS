package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
	"time"
)

type DeliveryOrderService struct {
	doRepo *repository.DeliveryOrderRepository
}

func NewDeliveryOrderService() *DeliveryOrderService {
	return &DeliveryOrderService{
		doRepo: repository.NewDeliveryOrderRepository(),
	}
}

// Business logic methods
func (s *DeliveryOrderService) GetAllDeliveryOrders() (interface{}, error) {
	return s.doRepo.GetAllDeliveryOrders()
}

// GetFilteredDeliveryOrders returns delivery orders with optional filters
func (s *DeliveryOrderService) GetFilteredDeliveryOrders(purchaseOrderID *uint, deliveryDateFrom *time.Time, deliveryDateTo *time.Time, status *string) (interface{}, error) {
	return s.doRepo.GetFilteredDeliveryOrders(purchaseOrderID, deliveryDateFrom, deliveryDateTo, status)
}

func (s *DeliveryOrderService) GetDeliveryOrderByID(id uint) (interface{}, error) {
	order, err := s.doRepo.GetDeliveryOrderByID(id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// GetDeliveryOrderWithItems returns delivery order with all its items
func (s *DeliveryOrderService) GetDeliveryOrderWithItems(id uint) (interface{}, error) {
	order, err := s.doRepo.GetDeliveryOrderWithItems(id)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *DeliveryOrderService) CreateDeliveryOrder(doNumber string, purchaseOrderID uint, deliveryDate time.Time, status string, description *string, userID uint) (interface{}, error) {
	if doNumber == "" {
		return nil, errors.New("DO number is required")
	}

	if purchaseOrderID == 0 {
		return nil, errors.New("purchase order ID is required")
	}

	if deliveryDate.IsZero() {
		return nil, errors.New("delivery date is required")
	}

	if status == "" {
		status = "draft" // Default status
	}

	// Validate status
	validStatuses := map[string]bool{
		"draft":    true,
		"shipped":  true,
		"received": true,
		"closed":   true,
	}
	if !validStatuses[status] {
		return nil, errors.New("invalid status: must be one of draft, shipped, received, closed")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if purchase order exists
	poExists, err := s.doRepo.CheckPurchaseOrderExists(purchaseOrderID)
	if err != nil {
		return nil, err
	}
	if !poExists {
		return nil, errors.New("purchase order not found")
	}

	// Check if DO number already exists
	doExists, err := s.doRepo.CheckDONumberExists(doNumber)
	if err != nil {
		return nil, err
	}
	if doExists {
		return nil, errors.New("DO number already exists")
	}

	order := &model.DeliveryOrder{
		DONumber:        doNumber,
		PurchaseOrderID: purchaseOrderID,
		DeliveryDate:    deliveryDate,
		Status:          status,
		Description:     description,
		UserIns:         &userID,
	}

	err = s.doRepo.CreateDeliveryOrder(order)
	if err != nil {
		return nil, err
	}

	// Fetch the created order with relationships
	createdOrder, err := s.doRepo.GetDeliveryOrderByID(order.ID)
	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}

// CreateDeliveryOrderWithItems creates a delivery order with multiple items in a transaction
func (s *DeliveryOrderService) CreateDeliveryOrderWithItems(doNumber string, purchaseOrderID uint, deliveryDate time.Time, status string, description *string, items []model.DeliveryOrderItem, userID uint) (interface{}, error) {
	if doNumber == "" {
		return nil, errors.New("DO number is required")
	}

	if purchaseOrderID == 0 {
		return nil, errors.New("purchase order ID is required")
	}

	if deliveryDate.IsZero() {
		return nil, errors.New("delivery date is required")
	}

	if status == "" {
		status = "draft" // Default status
	}

	// Validate status
	validStatuses := map[string]bool{
		"draft":    true,
		"shipped":  true,
		"received": true,
		"closed":   true,
	}
	if !validStatuses[status] {
		return nil, errors.New("invalid status: must be one of draft, shipped, received, closed")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if purchase order exists
	poExists, err := s.doRepo.CheckPurchaseOrderExists(purchaseOrderID)
	if err != nil {
		return nil, err
	}
	if !poExists {
		return nil, errors.New("purchase order not found")
	}

	// Check if DO number already exists
	doExists, err := s.doRepo.CheckDONumberExists(doNumber)
	if err != nil {
		return nil, err
	}
	if doExists {
		return nil, errors.New("DO number already exists")
	}

	// Validate items
	if len(items) == 0 {
		return nil, errors.New("at least one item is required")
	}

	// Set user_ins for all items
	for i := range items {
		items[i].UserIns = &userID
	}

	order := &model.DeliveryOrder{
		DONumber:        doNumber,
		PurchaseOrderID: purchaseOrderID,
		DeliveryDate:    deliveryDate,
		Status:          status,
		Description:     description,
		UserIns:         &userID,
	}

	err = s.doRepo.CreateDeliveryOrderWithItems(order, items)
	if err != nil {
		return nil, err
	}

	// Fetch the created order with relationships
	createdOrder, err := s.doRepo.GetDeliveryOrderByID(order.ID)
	if err != nil {
		return nil, err
	}

	return createdOrder, nil
}

func (s *DeliveryOrderService) UpdateDeliveryOrder(id uint, doNumber string, purchaseOrderID uint, deliveryDate time.Time, status string, description *string, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid delivery order ID")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if delivery order exists
	oldOrder, err := s.doRepo.GetDeliveryOrderModelByID(id)
	if err != nil {
		return nil, errors.New("delivery order not found")
	}

	// Validate status if provided
	if status != "" {
		validStatuses := map[string]bool{
			"draft":    true,
			"shipped":  true,
			"received": true,
			"closed":   true,
		}
		if !validStatuses[status] {
			return nil, errors.New("invalid status: must be one of draft, shipped, received, closed")
		}
	}

	// If purchase order ID is being changed, check if new purchase order exists
	if purchaseOrderID != 0 && purchaseOrderID != oldOrder.PurchaseOrderID {
		poExists, err := s.doRepo.CheckPurchaseOrderExists(purchaseOrderID)
		if err != nil {
			return nil, err
		}
		if !poExists {
			return nil, errors.New("purchase order not found")
		}
	}

	// Prepare update data with audit trail
	updateData := make(map[string]interface{})
	if doNumber != "" && doNumber != oldOrder.DONumber {
		// Check if new DO number already exists
		doExists, err := s.doRepo.CheckDONumberExists(doNumber)
		if err != nil {
			return nil, err
		}
		if doExists {
			return nil, errors.New("DO number already exists")
		}
		updateData["do_number"] = doNumber
	}
	if purchaseOrderID != 0 && purchaseOrderID != oldOrder.PurchaseOrderID {
		updateData["purchase_order_id"] = purchaseOrderID
	}
	if !deliveryDate.IsZero() {
		updateData["delivery_date"] = deliveryDate
	}
	if status != "" {
		updateData["status"] = status
	}
	if description != nil {
		updateData["description"] = description
	}
	// Always set the user who updated
	updateData["user_updt"] = userID

	err = s.doRepo.UpdateDeliveryOrder(id, updateData)
	if err != nil {
		return nil, err
	}

	updatedOrder, err := s.doRepo.GetDeliveryOrderByID(id)
	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

func (s *DeliveryOrderService) DeleteDeliveryOrder(id uint, userID uint) error {
	if id == 0 {
		return errors.New("invalid delivery order ID")
	}

	if userID == 0 {
		return errors.New("user ID is required for audit trail")
	}

	// Check if delivery order exists
	_, err := s.doRepo.GetDeliveryOrderModelByID(id)
	if err != nil {
		return errors.New("delivery order not found")
	}

	return s.doRepo.DeleteDeliveryOrderWithAudit(id, userID)
}

// GetDeletedDeliveryOrders returns all soft deleted delivery orders
func (s *DeliveryOrderService) GetDeletedDeliveryOrders() (interface{}, error) {
	return s.doRepo.GetDeletedDeliveryOrders()
}

// RestoreDeliveryOrder restores a soft deleted delivery order
func (s *DeliveryOrderService) RestoreDeliveryOrder(id uint, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid delivery order ID")
	}
	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	err := s.doRepo.RestoreDeliveryOrder(id, userID)
	if err != nil {
		return nil, err
	}

	restoredOrder, err := s.doRepo.GetDeliveryOrderByID(id)
	if err != nil {
		return nil, err
	}
	return restoredOrder, nil
}
