package seeder

import (
	"log"
	"myapp/internal/model"

	"gorm.io/gorm"
)

type ProductUnitSeeder struct{}

func NewProductUnitSeeder() SeederInterface {
	return &ProductUnitSeeder{}
}

func (s *ProductUnitSeeder) GetName() string {
	return "ProductUnitSeeder"
}

func (s *ProductUnitSeeder) Seed(db *gorm.DB) error {
	log.Println("🌱 Running ProductUnitSeeder...")

	// Check if data already exists
	var count int64
	db.Model(&model.ProductUnit{}).Count(&count)
	if count > 0 {
		log.Println("ProductUnit data already exists, skipping...")
		return nil
	}

	// Get some products for reference
	var products []model.Product
	if err := db.Limit(5).Find(&products).Error; err != nil {
		return err
	}

	if len(products) == 0 {
		log.Println("No products found, skipping ProductUnit seeding...")
		return nil
	}

	// Get some locations for reference
	var locations []model.Location
	if err := db.Limit(3).Find(&locations).Error; err != nil {
		log.Printf("Error finding locations: %v", err)
		return err
	}

	if len(locations) == 0 {
		log.Println("No locations found, ProductUnit seeding requires locations. Please seed locations first.")
		return nil
	}

	// Get some product batch for reference
	var productBatches []model.ProductBatch
	if err := db.Limit(3).Find(&productBatches).Error; err != nil {
		log.Printf("Error finding product batches: %v", err)
		return err
	}

	if len(productBatches) == 0 {
		log.Println("No product batches found, ProductUnit seeding requires product batches. Please seed product batches first.")
		return nil
	} // Get admin user for audit trail
	var adminUser model.User
	if err := db.Where("email = ?", "admin@wms.com").First(&adminUser).Error; err != nil {
		log.Println("Admin user not found, using nil for audit trail")
	}

	var userID *uint
	if adminUser.ID != 0 {
		userID = &adminUser.ID
	}

	productUnits := []model.ProductUnit{
		{
			ProductID:       products[0].ID,
			LocationID:      getLocationID(locations, 0),
			ProductBatchID:  getProductBatchID(productBatches, 0),
			Name:            stringPtr("Piece"),
			Quantity:        float64Ptr(1),
			UnitPrice:       float64Ptr(15000),
			UnitPriceRetail: float64Ptr(18000),
			Barcode:         stringPtr("8991234567890"),
			Description:     stringPtr("Single unit piece"),
			UserIns:         userID,
		},
		{
			ProductID:       products[0].ID,
			LocationID:      getLocationID(locations, 1),
			ProductBatchID:  getProductBatchID(productBatches, 0),
			Name:            stringPtr("Box"),
			Quantity:        float64Ptr(12),
			UnitPrice:       float64Ptr(170000),
			UnitPriceRetail: float64Ptr(210000),
			Barcode:         stringPtr("8991234567891"),
			Description:     stringPtr("Box of 12 pieces"),
			UserIns:         userID,
		},
		{
			ProductID:       products[1].ID,
			LocationID:      getLocationID(locations, 0),
			ProductBatchID:  getProductBatchID(productBatches, 0),
			Name:            stringPtr("Liter"),
			Quantity:        float64Ptr(1),
			UnitPrice:       float64Ptr(25000),
			UnitPriceRetail: float64Ptr(30000),
			Barcode:         stringPtr("8991234567892"),
			Description:     stringPtr("1 liter bottle"),
			UserIns:         userID,
		},
		{
			ProductID:       products[1].ID,
			LocationID:      getLocationID(locations, 2),
			ProductBatchID:  getProductBatchID(productBatches, 1),
			Name:            stringPtr("Gallon"),
			Quantity:        float64Ptr(4),
			UnitPrice:       float64Ptr(95000),
			UnitPriceRetail: float64Ptr(115000),
			Barcode:         stringPtr("8991234567893"),
			Description:     stringPtr("4 liter gallon"),
			UserIns:         userID,
		},
		{
			ProductID:       products[2].ID,
			LocationID:      getLocationID(locations, 0),
			ProductBatchID:  getProductBatchID(productBatches, 2),
			Name:            stringPtr("Kilogram"),
			Quantity:        float64Ptr(1),
			UnitPrice:       float64Ptr(50000),
			UnitPriceRetail: float64Ptr(62000),
			Barcode:         stringPtr("8991234567894"),
			Description:     stringPtr("1 kilogram pack"),
			UserIns:         userID,
		},
	}

	// Additional units for more products if available
	if len(products) > 3 {
		additionalUnits := []model.ProductUnit{
			{
				ProductID:       products[3].ID,
				LocationID:      getLocationID(locations, 1),
				ProductBatchID:  getProductBatchID(productBatches, 0),
				Name:            stringPtr("Unit"),
				Quantity:        float64Ptr(1),
				UnitPrice:       float64Ptr(75000),
				UnitPriceRetail: float64Ptr(90000),
				Barcode:         stringPtr("8991234567895"),
				Description:     stringPtr("Single unit"),
				UserIns:         userID,
			},
			{
				ProductID:       products[3].ID,
				LocationID:      getLocationID(locations, 2),
				ProductBatchID:  getProductBatchID(productBatches, 1),
				Name:            stringPtr("Dozen"),
				Quantity:        float64Ptr(12),
				UnitPrice:       float64Ptr(850000),
				UnitPriceRetail: float64Ptr(1020000),
				Barcode:         stringPtr("8991234567896"),
				Description:     stringPtr("Dozen pack"),
				UserIns:         userID,
			},
		}
		productUnits = append(productUnits, additionalUnits...)
	}

	if len(products) > 4 {
		moreUnits := []model.ProductUnit{
			{
				ProductID:       products[4].ID,
				LocationID:      getLocationID(locations, 0),
				ProductBatchID:  getProductBatchID(productBatches, 2),
				Name:            stringPtr("Gram"),
				Quantity:        float64Ptr(0.1),
				UnitPrice:       float64Ptr(5000),
				UnitPriceRetail: float64Ptr(6000),
				Barcode:         stringPtr("8991234567897"),
				Description:     stringPtr("100 gram pack"),
				UserIns:         userID,
			},
		}
		productUnits = append(productUnits, moreUnits...)
	}

	// Insert product units
	for _, unit := range productUnits {
		if err := db.Create(&unit).Error; err != nil {
			log.Printf("Failed to create product unit: %v", err)
			return err
		}
	}

	log.Printf("✅ Created %d product units", len(productUnits))
	return nil
}

func getLocationID(locations []model.Location, index int) uint {
	if len(locations) == 0 {
		log.Printf("⚠️  Warning: No locations available, returning 1 as default")
		return 1 // Default to first location ID (assuming location with ID 1 exists)
	}
	if index >= len(locations) {
		index = 0 // Default to first location if index out of bounds
	}
	return locations[index].ID
}

func getProductBatchID(productBatches []model.ProductBatch, index int) uint {
	if len(productBatches) == 0 {
		log.Printf("⚠️  Warning: No product batches available, returning 1 as default")
		return 1 // Default to first batch ID (assuming batch with ID 1 exists)
	}
	if index >= len(productBatches) {
		index = 0 // Default to first batch if index out of bounds
	}
	return productBatches[index].ID
}
