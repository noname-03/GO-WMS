package seeder

import (
	"log"

	"gorm.io/gorm"
)

// GetAllSeeders returns all available seeders
func GetAllSeeders() *SeederRegistry {
	registry := NewSeederRegistry()

	// Register all seeders here
	// Add new seeders to this list to auto-run them
	registry.Register(NewUserSeeder())
	registry.Register(NewRoleSeeder())
	registry.Register(NewBrandSeeder())
	registry.Register(NewCategorySeeder())
	registry.Register(NewProductSeeder())
	registry.Register(NewProductBatchSeeder())
	registry.Register(NewProductBatchTrackSeeder())
	registry.Register(NewProductUnitSeeder())
	registry.Register(NewProductUnitTrackSeeder())
	registry.Register(NewLocationSeeder())
	// registry.Register(NewWarehouseSeeder())

	// Future seeders:
	// registry.Register(NewCategorySeeder())
	// registry.Register(NewOrderSeeder())
	// registry.Register(NewInventorySeeder())

	return registry
}

// RunAllSeeders executes all registered seeders
func RunAllSeeders(db *gorm.DB) error {
	log.Println("🌱 Starting database seeding...")

	registry := GetAllSeeders()

	log.Printf("📋 Found %d seeders to run", len(registry.GetSeeders()))

	err := registry.RunAll(db)
	if err != nil {
		log.Printf("❌ Seeding failed: %v", err)
		return err
	}

	log.Println("✅ Database seeding completed successfully!")
	return nil
}
