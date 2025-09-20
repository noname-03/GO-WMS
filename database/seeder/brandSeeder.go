package seeder

import (
	"log"
	"myapp/internal/model"

	"gorm.io/gorm"
)

type BrandSeeder struct{}

func NewBrandSeeder() SeederInterface {
	return &BrandSeeder{}
}

func (s *BrandSeeder) GetName() string {
	return "BrandSeeder"
}

func (s *BrandSeeder) Seed(db *gorm.DB) error {
	log.Println("🌱 Running BrandSeeder...")

	// Description as pointer to string for nullable field
	desc1 := "Premium automotive brand"
	desc2 := "Consumer electronics and technology"
	desc3 := "Sports and athletic wear"

	brands := []model.Brand{
		{Name: "Toyota", Description: &desc1},
		{Name: "Samsung", Description: &desc2},
		{Name: "Nike", Description: &desc3},
		{Name: "Apple", Description: nil}, // Null description example
	}

	for _, brand := range brands {
		var existing model.Brand
		result := db.Where("name = ?", brand.Name).First(&existing)
		if result.Error != nil {
			// Brand doesn't exist, create it
			if err := db.Create(&brand).Error; err != nil {
				log.Printf("❌ Failed to seed brand %s: %v", brand.Name, err)
				return err
			}
			log.Printf("✅ Brand '%s' created successfully", brand.Name)
		} else {
			log.Printf("✅ BrandSeeder: Brand '%s' already exists, skipping...", brand.Name)
		}
	}

	return nil
}
