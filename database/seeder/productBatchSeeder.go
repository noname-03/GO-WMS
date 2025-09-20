package seeder

import (
	"log"
	"myapp/internal/model"
	"time"

	"gorm.io/gorm"
)

type ProductBatchSeeder struct{}

func NewProductBatchSeeder() SeederInterface {
	return &ProductBatchSeeder{}
}

func (s *ProductBatchSeeder) GetName() string {
	return "ProductBatchSeeder"
}

func (s *ProductBatchSeeder) Seed(db *gorm.DB) error {
	log.Println("🌱 Running ProductBatchSeeder...")

	// Get products first to create relationships
	var camry, galaxy, airZoom, iphone15 model.Product
	db.Where("name LIKE ?", "%Camry%").First(&camry)
	db.Where("name LIKE ?", "%Galaxy S24%").First(&galaxy)
	db.Where("name LIKE ?", "%Air Zoom%").First(&airZoom)
	db.Where("name LIKE ?", "%iPhone 15 Pro%").First(&iphone15)

	// Code batch and descriptions as pointers for nullable fields
	batch1 := "BATCH-CAM-2024-001"
	batch2 := "BATCH-GAL-2024-002"
	batch3 := "BATCH-NIK-2024-003"
	batch4 := "BATCH-IPH-2024-004"

	desc1 := "Toyota Camry batch Q1 2024"
	desc2 := "Samsung Galaxy S24 batch Q1 2024"
	desc3 := "Nike Air Zoom batch Q1 2024"

	price1 := 35000.00
	price2 := 12000000.00
	price3 := 2500000.00
	price4 := 18000000.00

	// Calculate expiry dates
	now := time.Now()
	futureDate1 := now.AddDate(2, 0, 0) // 2 years from now
	futureDate2 := now.AddDate(1, 6, 0) // 1.5 years from now
	futureDate3 := now.AddDate(3, 0, 0) // 3 years from now
	futureDate4 := now.AddDate(1, 0, 0) // 1 year from now

	productBatches := []model.ProductBatch{
		// Toyota Camry batches
		{ProductID: camry.ID, CodeBatch: &batch1, UnitPrice: &price1, ExpDate: futureDate1, Description: &desc1},
		{ProductID: camry.ID, CodeBatch: nil, UnitPrice: &price1, ExpDate: futureDate1, Description: nil},

		// Samsung Galaxy batches
		{ProductID: galaxy.ID, CodeBatch: &batch2, UnitPrice: &price2, ExpDate: futureDate2, Description: &desc2},
		{ProductID: galaxy.ID, CodeBatch: nil, UnitPrice: &price2, ExpDate: futureDate2, Description: nil},

		// Nike Air Zoom batches
		{ProductID: airZoom.ID, CodeBatch: &batch3, UnitPrice: &price3, ExpDate: futureDate3, Description: &desc3},
		{ProductID: airZoom.ID, CodeBatch: nil, UnitPrice: &price3, ExpDate: futureDate3, Description: nil},

		// iPhone 15 Pro batches
		{ProductID: iphone15.ID, CodeBatch: &batch4, UnitPrice: &price4, ExpDate: futureDate4, Description: nil},
		{ProductID: iphone15.ID, CodeBatch: nil, UnitPrice: nil, ExpDate: futureDate4, Description: nil},
	}

	for _, batch := range productBatches {
		var existing model.ProductBatch
		// Check by product_id and exp_date to avoid exact duplicates
		result := db.Where("product_id = ? AND exp_date = ?", batch.ProductID, batch.ExpDate).First(&existing)
		if result.Error != nil {
			// ProductBatch doesn't exist, create it
			if err := db.Create(&batch).Error; err != nil {
				log.Printf("❌ Failed to seed product batch for product ID %d: %v", batch.ProductID, err)
				return err
			}
			log.Printf("✅ Product batch for product ID %d created successfully", batch.ProductID)
		} else {
			log.Printf("✅ ProductBatchSeeder: Product batch for product ID %d already exists, skipping...", batch.ProductID)
		}
	}

	return nil
}
