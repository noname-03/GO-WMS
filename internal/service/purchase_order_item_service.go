package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
)

type PurchaseOrderItemService struct {
	poiRepo *repository.PurchaseOrderItemRepository
}

func NewPurchaseOrderItemService() *PurchaseOrderItemService {
	return &PurchaseOrderItemService{
		poiRepo: repository.NewPurchaseOrderItemRepository(),
	}
}

// Business logic methods
func (s *PurchaseOrderItemService) GetAllPurchaseOrderItems() (interface{}, error) {
	return s.poiRepo.GetAllPurchaseOrderItems()
}

func (s *PurchaseOrderItemService) GetPurchaseOrderItemsByPurchaseOrder(purchaseOrderID uint) (interface{}, error) {
	if purchaseOrderID == 0 {
		return nil, errors.New("invalid purchase order ID")
	}

	// Check if purchase order exists
	poExists, err := s.poiRepo.CheckPurchaseOrderExists(purchaseOrderID)
	if err != nil {
		return nil, err
	}
	if !poExists {
		return nil, errors.New("purchase order not found")
	}

	return s.poiRepo.GetPurchaseOrderItemsByPurchaseOrder(purchaseOrderID)
}

func (s *PurchaseOrderItemService) GetPurchaseOrderItemByID(id uint) (interface{}, error) {
	item, err := s.poiRepo.GetPurchaseOrderItemByID(id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (s *PurchaseOrderItemService) CreatePurchaseOrderItem(purchaseOrderID uint, productID uint, qtyOrdered *float64, unitPrice *float64, discount *float64, totalPrice *float64, description *string, userID uint) (interface{}, error) {
	if purchaseOrderID == 0 {
		return nil, errors.New("purchase order ID is required")
	}

	if productID == 0 {
		return nil, errors.New("product ID is required")
	}

	if qtyOrdered == nil || *qtyOrdered <= 0 {
		return nil, errors.New("quantity ordered must be greater than 0")
	}

	if unitPrice == nil || *unitPrice < 0 {
		return nil, errors.New("unit price must be greater than or equal to 0")
	}

	if totalPrice == nil || *totalPrice < 0 {
		return nil, errors.New("total price must be greater than or equal to 0")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if purchase order exists
	poExists, err := s.poiRepo.CheckPurchaseOrderExists(purchaseOrderID)
	if err != nil {
		return nil, err
	}
	if !poExists {
		return nil, errors.New("purchase order not found")
	}

	// Check if product exists
	productExists, err := s.poiRepo.CheckProductExists(productID)
	if err != nil {
		return nil, err
	}
	if !productExists {
		return nil, errors.New("product not found")
	}

	// Set default discount if not provided
	if discount == nil {
		zero := 0.0
		discount = &zero
	}

	item := &model.PurchaseOrderItem{
		PurchaseOrderID: purchaseOrderID,
		ProductID:       productID,
		QtyOrdered:      qtyOrdered,
		UnitPrice:       unitPrice,
		Discount:        discount,
		TotalPrice:      totalPrice,
		Description:     description,
		UserIns:         &userID,
	}

	err = s.poiRepo.CreatePurchaseOrderItem(item)
	if err != nil {
		return nil, err
	}

	// Fetch the created item with relationships
	createdItem, err := s.poiRepo.GetPurchaseOrderItemByID(item.ID)
	if err != nil {
		return nil, err
	}

	return createdItem, nil
}

func (s *PurchaseOrderItemService) UpdatePurchaseOrderItem(id uint, purchaseOrderID uint, productID uint, qtyOrdered *float64, unitPrice *float64, discount *float64, totalPrice *float64, description *string, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid purchase order item ID")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if item exists
	oldItem, err := s.poiRepo.GetPurchaseOrderItemModelByID(id)
	if err != nil {
		return nil, errors.New("purchase order item not found")
	}

	// Validate quantities and prices if provided
	if qtyOrdered != nil && *qtyOrdered <= 0 {
		return nil, errors.New("quantity ordered must be greater than 0")
	}

	if unitPrice != nil && *unitPrice < 0 {
		return nil, errors.New("unit price must be greater than or equal to 0")
	}

	if totalPrice != nil && *totalPrice < 0 {
		return nil, errors.New("total price must be greater than or equal to 0")
	}

	// If purchase order ID is being changed, check if new purchase order exists
	if purchaseOrderID != 0 && purchaseOrderID != oldItem.PurchaseOrderID {
		poExists, err := s.poiRepo.CheckPurchaseOrderExists(purchaseOrderID)
		if err != nil {
			return nil, err
		}
		if !poExists {
			return nil, errors.New("purchase order not found")
		}
	}

	// If product ID is being changed, check if new product exists
	if productID != 0 && productID != oldItem.ProductID {
		productExists, err := s.poiRepo.CheckProductExists(productID)
		if err != nil {
			return nil, err
		}
		if !productExists {
			return nil, errors.New("product not found")
		}
	}

	// Prepare update data with audit trail
	updateData := make(map[string]interface{})
	if purchaseOrderID != 0 && purchaseOrderID != oldItem.PurchaseOrderID {
		updateData["purchase_order_id"] = purchaseOrderID
	}
	if productID != 0 && productID != oldItem.ProductID {
		updateData["product_id"] = productID
	}
	if qtyOrdered != nil {
		updateData["qty_ordered"] = qtyOrdered
	}
	if unitPrice != nil {
		updateData["unit_price"] = unitPrice
	}
	if discount != nil {
		updateData["discount"] = discount
	}
	if totalPrice != nil {
		updateData["total_price"] = totalPrice
	}
	if description != nil {
		updateData["description"] = description
	}
	// Always set the user who updated
	updateData["user_updt"] = userID

	err = s.poiRepo.UpdatePurchaseOrderItem(id, updateData)
	if err != nil {
		return nil, err
	}

	updatedItem, err := s.poiRepo.GetPurchaseOrderItemByID(id)
	if err != nil {
		return nil, err
	}

	return updatedItem, nil
}

func (s *PurchaseOrderItemService) DeletePurchaseOrderItem(id uint, userID uint) error {
	if id == 0 {
		return errors.New("invalid purchase order item ID")
	}

	if userID == 0 {
		return errors.New("user ID is required for audit trail")
	}

	// Check if item exists
	_, err := s.poiRepo.GetPurchaseOrderItemModelByID(id)
	if err != nil {
		return errors.New("purchase order item not found")
	}

	return s.poiRepo.DeletePurchaseOrderItemWithAudit(id, userID)
}

// GetDeletedPurchaseOrderItems returns all soft deleted purchase order items
func (s *PurchaseOrderItemService) GetDeletedPurchaseOrderItems() (interface{}, error) {
	return s.poiRepo.GetDeletedPurchaseOrderItems()
}

// RestorePurchaseOrderItem restores a soft deleted purchase order item
func (s *PurchaseOrderItemService) RestorePurchaseOrderItem(id uint, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid purchase order item ID")
	}
	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	err := s.poiRepo.RestorePurchaseOrderItem(id, userID)
	if err != nil {
		return nil, err
	}

	restoredItem, err := s.poiRepo.GetPurchaseOrderItemByID(id)
	if err != nil {
		return nil, err
	}
	return restoredItem, nil
}
