package seeder

import (
	"myapp/internal/model"
	"time"

	"gorm.io/gorm"
)

type DeliveryOrderSeeder struct{}

func NewDeliveryOrderSeeder() SeederInterface {
	return &DeliveryOrderSeeder{}
}

func (s *DeliveryOrderSeeder) GetName() string {
	return "DeliveryOrderSeeder"
}

func (s *DeliveryOrderSeeder) Seed(db *gorm.DB) error {
	deliveryOrders := []model.DeliveryOrder{
		{
			DONumber:        "DO-001",
			PurchaseOrderID: 1,
			DeliveryDate:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Status:          "received",
			Description:     stringPtr("First delivery - all items received"),
			UserIns:         uintPtr(1),
		},
		{
			DONumber:        "DO-002",
			PurchaseOrderID: 2,
			DeliveryDate:    time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			Status:          "shipped",
			Description:     stringPtr("Second delivery - in transit"),
			UserIns:         uintPtr(1),
		},
		{
			DONumber:        "DO-003",
			PurchaseOrderID: 3,
			DeliveryDate:    time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC),
			Status:          "received",
			Description:     stringPtr("Third delivery - completed"),
			UserIns:         uintPtr(2),
		},
		{
			DONumber:        "DO-004",
			PurchaseOrderID: 4,
			DeliveryDate:    time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			Status:          "draft",
			Description:     stringPtr("Fourth delivery - being prepared"),
			UserIns:         uintPtr(2),
		},
		{
			DONumber:        "DO-005",
			PurchaseOrderID: 5,
			DeliveryDate:    time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC),
			Status:          "closed",
			Description:     stringPtr("Fifth delivery - closed and archived"),
			UserIns:         uintPtr(1),
		},
	}

	for _, do := range deliveryOrders {
		// Check if DO number already exists
		var existing model.DeliveryOrder
		err := db.Where("do_number = ?", do.DONumber).First(&existing).Error
		if err == nil {
			// DO already exists, skip
			continue
		} else if err != gorm.ErrRecordNotFound {
			// Unexpected error
			return err
		}

		// Create delivery order if not exists
		if err := db.Create(&do).Error; err != nil {
			return err
		}
	}

	return nil
}
