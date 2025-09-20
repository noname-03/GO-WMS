package seeder

import (
	"log"
	"myapp/internal/model"

	"gorm.io/gorm"
)

type CategorySeeder struct{}

func NewCategorySeeder() SeederInterface {
	return &CategorySeeder{}
}

func (s *CategorySeeder) GetName() string {
	return "CategorySeeder"
}

func (s *CategorySeeder) Seed(db *gorm.DB) error {
	log.Println("🌱 Running CategorySeeder...")

	// Get brands first to create relationships
	var toyota, samsung, nike, apple model.Brand
	db.Where("name = ?", "Toyota").First(&toyota)
	db.Where("name = ?", "Samsung").First(&samsung)
	db.Where("name = ?", "Nike").First(&nike)
	db.Where("name = ?", "Apple").First(&apple)

	// Description as pointer to string for nullable field
	desc1 := "Vehicle categories"
	desc2 := "Consumer electronics categories"
	desc3 := "Athletic wear categories"

	categories := []model.Category{
		// Toyota categories
		{BrandID: toyota.ID, Name: "Sedan", Description: &desc1},
		{BrandID: toyota.ID, Name: "SUV", Description: &desc1},
		{BrandID: toyota.ID, Name: "Hybrid", Description: &desc1},

		// Samsung categories
		{BrandID: samsung.ID, Name: "Smartphone", Description: &desc2},
		{BrandID: samsung.ID, Name: "TV", Description: &desc2},
		{BrandID: samsung.ID, Name: "Appliances", Description: nil},

		// Nike categories
		{BrandID: nike.ID, Name: "Running", Description: &desc3},
		{BrandID: nike.ID, Name: "Basketball", Description: &desc3},
		{BrandID: nike.ID, Name: "Lifestyle", Description: nil},

		// Apple categories
		{BrandID: apple.ID, Name: "iPhone", Description: nil},
		{BrandID: apple.ID, Name: "Mac", Description: nil},
		{BrandID: apple.ID, Name: "iPad", Description: nil},
	}

	for _, category := range categories {
		var existing model.Category
		result := db.Where("name = ? AND brand_id = ?", category.Name, category.BrandID).First(&existing)
		if result.Error != nil {
			// Category doesn't exist, create it
			if err := db.Create(&category).Error; err != nil {
				log.Printf("❌ Failed to seed category %s for brand ID %d: %v", category.Name, category.BrandID, err)
				return err
			}
			log.Printf("✅ Category '%s' for brand ID %d created successfully", category.Name, category.BrandID)
		} else {
			log.Printf("✅ CategorySeeder: Category '%s' for brand ID %d already exists, skipping...", category.Name, category.BrandID)
		}
	}

	return nil
}
