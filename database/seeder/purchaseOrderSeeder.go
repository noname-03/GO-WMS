package seeder

import (
	"log"
	"myapp/internal/model"
	"time"

	"gorm.io/gorm"
)

type PurchaseOrderSeeder struct{}

func NewPurchaseOrderSeeder() SeederInterface {
	return &PurchaseOrderSeeder{}
}

func (s *PurchaseOrderSeeder) GetName() string {
	return "PurchaseOrderSeeder"
}

func (s *PurchaseOrderSeeder) Seed(db *gorm.DB) error {
	log.Println("Seeding purchase orders...")

	// Sample purchase orders
	purchaseOrders := []model.PurchaseOrder{
		{
			PONumber:    "PO-2025-001",
			UserID:      1, // Assuming user ID 1 exists
			OrderDate:   time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
			Status:      "draft",
			TotalAmount: float64Ptr(1500000),
			Description: stringPtr("Purchase order for product restock"),
			UserIns:     uintPtr(1),
		},
		{
			PONumber:    "PO-2025-002",
			UserID:      2, // Assuming user ID 2 exists
			OrderDate:   time.Date(2025, 10, 5, 0, 0, 0, 0, time.UTC),
			Status:      "submitted",
			TotalAmount: float64Ptr(2500000),
			Description: stringPtr("Bulk order for reseller"),
			UserIns:     uintPtr(1),
		},
		{
			PONumber:    "PO-2025-003",
			UserID:      1,
			OrderDate:   time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC),
			Status:      "approved",
			TotalAmount: float64Ptr(3000000),
			Description: stringPtr("Approved purchase order"),
			UserIns:     uintPtr(1),
		},
		{
			PONumber:    "PO-2025-004",
			UserID:      2,
			OrderDate:   time.Date(2025, 10, 12, 0, 0, 0, 0, time.UTC),
			Status:      "received",
			TotalAmount: float64Ptr(1800000),
			Description: stringPtr("Purchase order received"),
			UserIns:     uintPtr(1),
		},
		{
			PONumber:    "PO-2025-005",
			UserID:      1,
			OrderDate:   time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC),
			Status:      "closed",
			TotalAmount: float64Ptr(2200000),
			Description: stringPtr("Completed purchase order"),
			UserIns:     uintPtr(1),
		},
	}

	for _, po := range purchaseOrders {
		// Use FirstOrCreate to skip duplicates
		var existingPO model.PurchaseOrder
		result := db.Where("po_number = ?", po.PONumber).FirstOrCreate(&existingPO, &po)
		if result.Error != nil {
			log.Printf("Error seeding purchase order %s: %v", po.PONumber, result.Error)
			return result.Error
		}
		if result.RowsAffected > 0 {
			log.Printf("Created purchase order: %s", po.PONumber)
		} else {
			log.Printf("Purchase order already exists, skipped: %s", po.PONumber)
		}
	}

	log.Println("Purchase order seeding completed!")
	return nil
}
