package seeder

import (
	"fmt"
	"myapp/internal/model"

	"gorm.io/gorm"
)

type InvoiceItemSeeder struct{}

func NewInvoiceItemSeeder() SeederInterface {
	return &InvoiceItemSeeder{}
}

func (s *InvoiceItemSeeder) GetName() string {
	return "InvoiceItemSeeder"
}

func (s *InvoiceItemSeeder) Seed(db *gorm.DB) error {
	fmt.Printf("🌱 Running %s...\n", s.GetName())
	invoiceItems := []model.InvoiceItem{
		// Items for Invoice 1 (INV-2024-001)
		{
			InvoiceID:   1,
			ProductID:   1,
			QtyInvoiced: 100.00,
			UnitPrice:   25000.00,
			TotalPrice:  2500000.00,
			Description: stringPtr("Laptop units for invoice 1"),
			UserIns:     1,
			UserUpdt:    1,
		},
		{
			InvoiceID:   1,
			ProductID:   2,
			QtyInvoiced: 50.00,
			UnitPrice:   50000.00,
			TotalPrice:  2500000.00,
			Description: stringPtr("Phone units for invoice 1"),
			UserIns:     1,
			UserUpdt:    1,
		},
		// Items for Invoice 2 (INV-2024-002)
		{
			InvoiceID:   2,
			ProductID:   3,
			QtyInvoiced: 200.00,
			UnitPrice:   15000.00,
			TotalPrice:  3000000.00,
			Description: stringPtr("Headphone units for invoice 2"),
			UserIns:     1,
			UserUpdt:    1,
		},
		{
			InvoiceID:   2,
			ProductID:   4,
			QtyInvoiced: 150.00,
			UnitPrice:   30000.00,
			TotalPrice:  4500000.00,
			Description: stringPtr("Keyboard units for invoice 2"),
			UserIns:     1,
			UserUpdt:    1,
		},
		// Items for Invoice 3 (INV-2024-003)
		{
			InvoiceID:   3,
			ProductID:   1,
			QtyInvoiced: 80.00,
			UnitPrice:   25000.00,
			TotalPrice:  2000000.00,
			Description: stringPtr("Laptop units for invoice 3"),
			UserIns:     1,
			UserUpdt:    1,
		},
		{
			InvoiceID:   3,
			ProductID:   5,
			QtyInvoiced: 100.00,
			UnitPrice:   80000.00,
			TotalPrice:  8000000.00,
			Description: stringPtr("Monitor units for invoice 3"),
			UserIns:     1,
			UserUpdt:    1,
		},
		// Items for Invoice 4 (INV-2024-004)
		{
			InvoiceID:   4,
			ProductID:   2,
			QtyInvoiced: 70.00,
			UnitPrice:   50000.00,
			TotalPrice:  3500000.00,
			Description: stringPtr("Phone units for invoice 4"),
			UserIns:     1,
			UserUpdt:    1,
		},
		// Items for Invoice 5 (INV-2024-005)
		{
			InvoiceID:   5,
			ProductID:   1,
			QtyInvoiced: 150.00,
			UnitPrice:   25000.00,
			TotalPrice:  3750000.00,
			Description: stringPtr("Laptop units for invoice 5"),
			UserIns:     1,
			UserUpdt:    1,
		},
		{
			InvoiceID:   5,
			ProductID:   3,
			QtyInvoiced: 250.00,
			UnitPrice:   15000.00,
			TotalPrice:  3750000.00,
			Description: stringPtr("Headphone units for invoice 5"),
			UserIns:     1,
			UserUpdt:    1,
		},
		{
			InvoiceID:   5,
			ProductID:   5,
			QtyInvoiced: 50.00,
			UnitPrice:   100000.00,
			TotalPrice:  5000000.00,
			Description: stringPtr("Monitor units for invoice 5"),
			UserIns:     1,
			UserUpdt:    1,
		},
	}

	for i, item := range invoiceItems {
		var existing model.InvoiceItem
		result := db.Where("invoice_id = ? AND product_id = ? AND qty_invoiced = ?", item.InvoiceID, item.ProductID, item.QtyInvoiced).First(&existing)
		if result.Error != nil {
			// Item not found, create it
			if err := db.Create(&item).Error; err != nil {
				fmt.Printf("❌ Error creating invoice item %d: %v\n", i+1, err)
				return err
			} else {
				fmt.Printf("✅ Invoice item %d created successfully\n", i+1)
			}
		} else {
			fmt.Printf("⏭️  Invoice item %d already exists, skipping...\n", i+1)
		}
	}

	fmt.Printf("✅ %s completed\n", s.GetName())
	return nil
}
