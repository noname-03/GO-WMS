package seeder

import (
	"log"
	"myapp/internal/model"

	"gorm.io/gorm"
)

type ProductSeeder struct{}

func NewProductSeeder() SeederInterface {
	return &ProductSeeder{}
}

func (s *ProductSeeder) GetName() string {
	return "ProductSeeder"
}

func (s *ProductSeeder) Seed(db *gorm.DB) error {
	log.Println("🌱 Running ProductSeeder...")

	// Get categories first to create relationships
	var sedan, smartphone, running, iphone model.Category
	db.Where("name = ?", "Sedan").First(&sedan)
	db.Where("name = ?", "Smartphone").First(&smartphone)
	db.Where("name = ?", "Running").First(&running)
	db.Where("name = ?", "iPhone").First(&iphone)

	// Description as pointer to string for nullable field
	desc1 := "Toyota sedan vehicle"
	desc2 := "Samsung smartphone device"
	desc3 := "Nike running shoes"
	desc4 := "Apple iPhone smartphone"

	products := []model.Product{
		// Toyota Sedan products
		{CategoryID: sedan.ID, Name: "Toyota Camry 2024", Description: &desc1},
		{CategoryID: sedan.ID, Name: "Toyota Corolla 2024", Description: &desc1},
		{CategoryID: sedan.ID, Name: "Toyota Avalon 2024", Description: nil},

		// Samsung Smartphone products
		{CategoryID: smartphone.ID, Name: "Samsung Galaxy S24", Description: &desc2},
		{CategoryID: smartphone.ID, Name: "Samsung Galaxy A54", Description: &desc2},
		{CategoryID: smartphone.ID, Name: "Samsung Galaxy Note 23", Description: nil},

		// Nike Running products
		{CategoryID: running.ID, Name: "Nike Air Zoom Pegasus", Description: &desc3},
		{CategoryID: running.ID, Name: "Nike React Infinity Run", Description: &desc3},
		{CategoryID: running.ID, Name: "Nike Air Max 270", Description: nil},

		// iPhone products
		{CategoryID: iphone.ID, Name: "iPhone 15 Pro", Description: &desc4},
		{CategoryID: iphone.ID, Name: "iPhone 15", Description: &desc4},
		{CategoryID: iphone.ID, Name: "iPhone 14", Description: nil},
	}

	for _, product := range products {
		var existing model.Product
		result := db.Where("name = ? AND category_id = ?", product.Name, product.CategoryID).First(&existing)
		if result.Error != nil {
			// Product doesn't exist, create it
			if err := db.Create(&product).Error; err != nil {
				log.Printf("❌ Failed to seed product %s for category ID %d: %v", product.Name, product.CategoryID, err)
				return err
			}
			log.Printf("✅ Product '%s' for category ID %d created successfully", product.Name, product.CategoryID)
		} else {
			log.Printf("✅ ProductSeeder: Product '%s' for category ID %d already exists, skipping...", product.Name, product.CategoryID)
		}
	}

	return nil
}
