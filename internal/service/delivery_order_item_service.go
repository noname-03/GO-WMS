package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
)

type DeliveryOrderItemService struct {
	doiRepo *repository.DeliveryOrderItemRepository
}

func NewDeliveryOrderItemService() *DeliveryOrderItemService {
	return &DeliveryOrderItemService{
		doiRepo: repository.NewDeliveryOrderItemRepository(),
	}
}

// Business logic methods
func (s *DeliveryOrderItemService) GetAllDeliveryOrderItems() (interface{}, error) {
	return s.doiRepo.GetAllDeliveryOrderItems()
}

func (s *DeliveryOrderItemService) GetDeliveryOrderItemByID(id uint) (interface{}, error) {
	item, err := s.doiRepo.GetDeliveryOrderItemByID(id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *DeliveryOrderItemService) GetDeliveryOrderItemsByDeliveryOrder(deliveryOrderID uint) (interface{}, error) {
	if deliveryOrderID == 0 {
		return nil, errors.New("delivery order ID is required")
	}
	return s.doiRepo.GetDeliveryOrderItemsByDeliveryOrder(deliveryOrderID)
}

func (s *DeliveryOrderItemService) CreateDeliveryOrderItem(deliveryOrderID uint, productID uint, qtyDelivered *float64, description *string, userID uint) (interface{}, error) {
	if deliveryOrderID == 0 {
		return nil, errors.New("delivery order ID is required")
	}

	if productID == 0 {
		return nil, errors.New("product ID is required")
	}

	if qtyDelivered == nil || *qtyDelivered <= 0 {
		return nil, errors.New("quantity delivered must be greater than 0")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if delivery order exists
	doExists, err := s.doiRepo.CheckDeliveryOrderExists(deliveryOrderID)
	if err != nil {
		return nil, err
	}
	if !doExists {
		return nil, errors.New("delivery order not found")
	}

	// Check if product exists
	productExists, err := s.doiRepo.CheckProductExists(productID)
	if err != nil {
		return nil, err
	}
	if !productExists {
		return nil, errors.New("product not found")
	}

	item := &model.DeliveryOrderItem{
		DeliveryOrderID: deliveryOrderID,
		ProductID:       productID,
		QtyDelivered:    qtyDelivered,
		Description:     description,
		UserIns:         &userID,
	}

	err = s.doiRepo.CreateDeliveryOrderItem(item)
	if err != nil {
		return nil, err
	}

	// Fetch the created item with relationships
	createdItem, err := s.doiRepo.GetDeliveryOrderItemByID(item.ID)
	if err != nil {
		return nil, err
	}

	return createdItem, nil
}

func (s *DeliveryOrderItemService) UpdateDeliveryOrderItem(id uint, deliveryOrderID uint, productID uint, qtyDelivered *float64, description *string, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid delivery order item ID")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if delivery order item exists
	oldItem, err := s.doiRepo.GetDeliveryOrderItemModelByID(id)
	if err != nil {
		return nil, errors.New("delivery order item not found")
	}

	// If delivery order ID is being changed, check if new delivery order exists
	if deliveryOrderID != 0 && deliveryOrderID != oldItem.DeliveryOrderID {
		doExists, err := s.doiRepo.CheckDeliveryOrderExists(deliveryOrderID)
		if err != nil {
			return nil, err
		}
		if !doExists {
			return nil, errors.New("delivery order not found")
		}
	}

	// If product ID is being changed, check if new product exists
	if productID != 0 && productID != oldItem.ProductID {
		productExists, err := s.doiRepo.CheckProductExists(productID)
		if err != nil {
			return nil, err
		}
		if !productExists {
			return nil, errors.New("product not found")
		}
	}

	// Prepare update data with audit trail
	updateData := make(map[string]interface{})
	if deliveryOrderID != 0 && deliveryOrderID != oldItem.DeliveryOrderID {
		updateData["delivery_order_id"] = deliveryOrderID
	}
	if productID != 0 && productID != oldItem.ProductID {
		updateData["product_id"] = productID
	}
	if qtyDelivered != nil && *qtyDelivered > 0 && (oldItem.QtyDelivered == nil || *qtyDelivered != *oldItem.QtyDelivered) {
		updateData["qty_delivered"] = qtyDelivered
	}
	if description != nil {
		updateData["description"] = description
	}
	// Always set the user who updated
	updateData["user_updt"] = userID

	err = s.doiRepo.UpdateDeliveryOrderItem(id, updateData)
	if err != nil {
		return nil, err
	}

	updatedItem, err := s.doiRepo.GetDeliveryOrderItemByID(id)
	if err != nil {
		return nil, err
	}

	return updatedItem, nil
}

func (s *DeliveryOrderItemService) DeleteDeliveryOrderItem(id uint, userID uint) error {
	if id == 0 {
		return errors.New("invalid delivery order item ID")
	}

	if userID == 0 {
		return errors.New("user ID is required for audit trail")
	}

	// Check if delivery order item exists
	_, err := s.doiRepo.GetDeliveryOrderItemModelByID(id)
	if err != nil {
		return errors.New("delivery order item not found")
	}

	return s.doiRepo.DeleteDeliveryOrderItemWithAudit(id, userID)
}

// GetDeletedDeliveryOrderItems returns all soft deleted delivery order items
func (s *DeliveryOrderItemService) GetDeletedDeliveryOrderItems() (interface{}, error) {
	return s.doiRepo.GetDeletedDeliveryOrderItems()
}

// RestoreDeliveryOrderItem restores a soft deleted delivery order item
func (s *DeliveryOrderItemService) RestoreDeliveryOrderItem(id uint, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid delivery order item ID")
	}
	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	err := s.doiRepo.RestoreDeliveryOrderItem(id, userID)
	if err != nil {
		return nil, err
	}

	restoredItem, err := s.doiRepo.GetDeliveryOrderItemByID(id)
	if err != nil {
		return nil, err
	}
	return restoredItem, nil
}
