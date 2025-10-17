package seeder

import (
	"myapp/internal/model"

	"gorm.io/gorm"
)

type DeliveryOrderItemSeeder struct{}

func NewDeliveryOrderItemSeeder() SeederInterface {
	return &DeliveryOrderItemSeeder{}
}

func (s *DeliveryOrderItemSeeder) GetName() string {
	return "DeliveryOrderItemSeeder"
}

func (s *DeliveryOrderItemSeeder) Seed(db *gorm.DB) error {
	deliveryOrderItems := []model.DeliveryOrderItem{
		// DO-001 items (from PO-001)
		{
			DeliveryOrderID: 1,
			ProductID:       1,
			QtyDelivered:    float64Ptr(100),
			Description:     stringPtr("Full quantity delivered"),
			UserIns:         uintPtr(1),
		},
		{
			DeliveryOrderID: 1,
			ProductID:       2,
			QtyDelivered:    float64Ptr(50),
			Description:     stringPtr("Full quantity delivered"),
			UserIns:         uintPtr(1),
		},
		// DO-002 items (from PO-002)
		{
			DeliveryOrderID: 2,
			ProductID:       3,
			QtyDelivered:    float64Ptr(75),
			Description:     stringPtr("Partial delivery - in transit"),
			UserIns:         uintPtr(1),
		},
		{
			DeliveryOrderID: 2,
			ProductID:       4,
			QtyDelivered:    float64Ptr(30),
			Description:     stringPtr("Partial delivery - in transit"),
			UserIns:         uintPtr(1),
		},
		// DO-003 items (from PO-003)
		{
			DeliveryOrderID: 3,
			ProductID:       5,
			QtyDelivered:    float64Ptr(200),
			Description:     stringPtr("Full quantity received"),
			UserIns:         uintPtr(2),
		},
		{
			DeliveryOrderID: 3,
			ProductID:       1,
			QtyDelivered:    float64Ptr(25),
			Description:     stringPtr("Full quantity received"),
			UserIns:         uintPtr(2),
		},
		// DO-004 items (from PO-004)
		{
			DeliveryOrderID: 4,
			ProductID:       2,
			QtyDelivered:    float64Ptr(150),
			Description:     stringPtr("Being prepared for shipment"),
			UserIns:         uintPtr(2),
		},
		{
			DeliveryOrderID: 4,
			ProductID:       3,
			QtyDelivered:    float64Ptr(60),
			Description:     stringPtr("Being prepared for shipment"),
			UserIns:         uintPtr(2),
		},
		// DO-005 items (from PO-005)
		{
			DeliveryOrderID: 5,
			ProductID:       4,
			QtyDelivered:    float64Ptr(300),
			Description:     stringPtr("Delivered and closed"),
			UserIns:         uintPtr(1),
		},
		{
			DeliveryOrderID: 5,
			ProductID:       5,
			QtyDelivered:    float64Ptr(40),
			Description:     stringPtr("Delivered and closed"),
			UserIns:         uintPtr(1),
		},
	}

	for _, item := range deliveryOrderItems {
		// Use FirstOrCreate to avoid duplicates
		var existing model.DeliveryOrderItem
		err := db.Where("delivery_order_id = ? AND product_id = ?", item.DeliveryOrderID, item.ProductID).
			First(&existing).Error

		if err == nil {
			// Item already exists, skip
			continue
		} else if err != gorm.ErrRecordNotFound {
			// Unexpected error
			return err
		}

		// Create item if not exists
		if err := db.Create(&item).Error; err != nil {
			return err
		}
	}

	return nil
}
