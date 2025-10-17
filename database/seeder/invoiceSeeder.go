package seeder

import (
	"fmt"
	"myapp/internal/model"
	"time"

	"gorm.io/gorm"
)

type InvoiceSeeder struct{}

func NewInvoiceSeeder() SeederInterface {
	return &InvoiceSeeder{}
}

func (s *InvoiceSeeder) GetName() string {
	return "InvoiceSeeder"
}

func (s *InvoiceSeeder) Seed(db *gorm.DB) error {
	fmt.Printf("🌱 Running %s...\n", s.GetName())
	invoices := []model.Invoice{
		{
			InvoiceNumber:   "INV-2024-001",
			UserID:          1,
			PurchaseOrderID: uintPtr(1),
			DeliveryOrderID: uintPtr(1),
			InvoiceDate:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Status:          "draft",
			TotalAmount:     5000000.00,
			Description:     stringPtr("Invoice for first purchase order"),
			UserIns:         1,
			UserUpdt:        1,
		},
		{
			InvoiceNumber:   "INV-2024-002",
			UserID:          2,
			PurchaseOrderID: uintPtr(2),
			DeliveryOrderID: uintPtr(2),
			InvoiceDate:     time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			Status:          "sent",
			TotalAmount:     7500000.00,
			Description:     stringPtr("Invoice for second purchase order"),
			UserIns:         1,
			UserUpdt:        1,
		},
		{
			InvoiceNumber:   "INV-2024-003",
			UserID:          1,
			PurchaseOrderID: uintPtr(3),
			DeliveryOrderID: uintPtr(3),
			InvoiceDate:     time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			Status:          "paid",
			TotalAmount:     10000000.00,
			Description:     stringPtr("Invoice for third purchase order"),
			UserIns:         1,
			UserUpdt:        1,
		},
		{
			InvoiceNumber:   "INV-2024-004",
			UserID:          2,
			PurchaseOrderID: uintPtr(4),
			DeliveryOrderID: nil,
			InvoiceDate:     time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
			Status:          "sent",
			TotalAmount:     3500000.00,
			Description:     stringPtr("Invoice without delivery order"),
			UserIns:         1,
			UserUpdt:        1,
		},
		{
			InvoiceNumber:   "INV-2024-005",
			UserID:          1,
			PurchaseOrderID: uintPtr(5),
			DeliveryOrderID: uintPtr(4),
			InvoiceDate:     time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC),
			Status:          "closed",
			TotalAmount:     12500000.00,
			Description:     stringPtr("Invoice completed and closed"),
			UserIns:         1,
			UserUpdt:        1,
		},
	}

	for _, invoice := range invoices {
		var existing model.Invoice
		result := db.Where("invoice_number = ?", invoice.InvoiceNumber).First(&existing)
		if result.Error != nil {
			// Invoice not found, create it
			if err := db.Create(&invoice).Error; err != nil {
				fmt.Printf("❌ Error creating invoice %s: %v\n", invoice.InvoiceNumber, err)
				return err
			} else {
				fmt.Printf("✅ Invoice %s created successfully\n", invoice.InvoiceNumber)
			}
		} else {
			fmt.Printf("⏭️  Invoice %s already exists, skipping...\n", invoice.InvoiceNumber)
		}
	}

	fmt.Printf("✅ %s completed\n", s.GetName())
	return nil
}
