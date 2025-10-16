package seeder

import (
	"log"
	"myapp/internal/model"

	"gorm.io/gorm"
)

type PurchaseOrderItemSeeder struct{}

func NewPurchaseOrderItemSeeder() SeederInterface {
	return &PurchaseOrderItemSeeder{}
}

func (s *PurchaseOrderItemSeeder) GetName() string {
	return "PurchaseOrderItemSeeder"
}

func (s *PurchaseOrderItemSeeder) Seed(db *gorm.DB) error {
	log.Println("Seeding purchase order items...")

	// Sample purchase order items
	purchaseOrderItems := []model.PurchaseOrderItem{
		// Items for PO-2025-001
		{
			PurchaseOrderID: 1, // Assuming PO ID 1 exists
			ProductID:       1, // Assuming product ID 1 exists
			QtyOrdered:      float64Ptr(100),
			UnitPrice:       float64Ptr(10000),
			Discount:        float64Ptr(0),
			TotalPrice:      float64Ptr(1000000),
			Description:     stringPtr("Product A - 100 units"),
			UserIns:         uintPtr(1),
		},
		{
			PurchaseOrderID: 1,
			ProductID:       2,
			QtyOrdered:      float64Ptr(50),
			UnitPrice:       float64Ptr(10000),
			Discount:        float64Ptr(0),
			TotalPrice:      float64Ptr(500000),
			Description:     stringPtr("Product B - 50 units"),
			UserIns:         uintPtr(1),
		},
		// Items for PO-2025-002
		{
			PurchaseOrderID: 2,
			ProductID:       1,
			QtyOrdered:      float64Ptr(150),
			UnitPrice:       float64Ptr(10000),
			Discount:        float64Ptr(50000),
			TotalPrice:      float64Ptr(1450000),
			Description:     stringPtr("Product A - 150 units with discount"),
			UserIns:         uintPtr(1),
		},
		{
			PurchaseOrderID: 2,
			ProductID:       3,
			QtyOrdered:      float64Ptr(70),
			UnitPrice:       float64Ptr(15000),
			Discount:        float64Ptr(0),
			TotalPrice:      float64Ptr(1050000),
			Description:     stringPtr("Product C - 70 units"),
			UserIns:         uintPtr(1),
		},
		// Items for PO-2025-003
		{
			PurchaseOrderID: 3,
			ProductID:       2,
			QtyOrdered:      float64Ptr(200),
			UnitPrice:       float64Ptr(10000),
			Discount:        float64Ptr(0),
			TotalPrice:      float64Ptr(2000000),
			Description:     stringPtr("Product B - 200 units"),
			UserIns:         uintPtr(1),
		},
		{
			PurchaseOrderID: 3,
			ProductID:       4,
			QtyOrdered:      float64Ptr(50),
			UnitPrice:       float64Ptr(20000),
			Discount:        float64Ptr(0),
			TotalPrice:      float64Ptr(1000000),
			Description:     stringPtr("Product D - 50 units"),
			UserIns:         uintPtr(1),
		},
		// Items for PO-2025-004
		{
			PurchaseOrderID: 4,
			ProductID:       1,
			QtyOrdered:      float64Ptr(120),
			UnitPrice:       float64Ptr(10000),
			Discount:        float64Ptr(20000),
			TotalPrice:      float64Ptr(1180000),
			Description:     stringPtr("Product A - 120 units with discount"),
			UserIns:         uintPtr(1),
		},
		{
			PurchaseOrderID: 4,
			ProductID:       3,
			QtyOrdered:      float64Ptr(40),
			UnitPrice:       float64Ptr(15000),
			Discount:        float64Ptr(0),
			TotalPrice:      float64Ptr(600000),
			Description:     stringPtr("Product C - 40 units"),
			UserIns:         uintPtr(1),
		},
		// Items for PO-2025-005
		{
			PurchaseOrderID: 5,
			ProductID:       2,
			QtyOrdered:      float64Ptr(100),
			UnitPrice:       float64Ptr(12000),
			Discount:        float64Ptr(0),
			TotalPrice:      float64Ptr(1200000),
			Description:     stringPtr("Product B - 100 units at new price"),
			UserIns:         uintPtr(1),
		},
		{
			PurchaseOrderID: 5,
			ProductID:       4,
			QtyOrdered:      float64Ptr(50),
			UnitPrice:       float64Ptr(20000),
			Discount:        float64Ptr(0),
			TotalPrice:      float64Ptr(1000000),
			Description:     stringPtr("Product D - 50 units"),
			UserIns:         uintPtr(1),
		},
	}

	for i, item := range purchaseOrderItems {
		// Use FirstOrCreate to skip duplicates based on purchase_order_id and product_id
		var existingItem model.PurchaseOrderItem
		result := db.Where("purchase_order_id = ? AND product_id = ?", item.PurchaseOrderID, item.ProductID).
			FirstOrCreate(&existingItem, &item)
		if result.Error != nil {
			log.Printf("Error seeding purchase order item #%d: %v", i+1, result.Error)
			return result.Error
		}
		if result.RowsAffected > 0 {
			log.Printf("Created purchase order item: PO ID %d, Product ID %d", item.PurchaseOrderID, item.ProductID)
		} else {
			log.Printf("Purchase order item already exists, skipped: PO ID %d, Product ID %d", item.PurchaseOrderID, item.ProductID)
		}
	}

	log.Println("Purchase order item seeding completed!")
	return nil
}
