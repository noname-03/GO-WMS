package seeder

import (
	"log"

	"gorm.io/gorm"
)

// GetAllSeeders returns all available seeders
func GetAllSeeders() *SeederRegistry {
	registry := NewSeederRegistry()

	// Register all seeders here in proper dependency order
	// Add new seeders to this list to auto-run them
	registry.Register(NewUserSeeder())         // Base users first
	registry.Register(NewRoleSeeder())         // Roles
	registry.Register(NewBrandSeeder())        // Product dependencies
	registry.Register(NewCategorySeeder())     // Product dependencies
	registry.Register(NewLocationSeeder())     // Location must be before ProductUnit
	registry.Register(NewProductSeeder())      // Products
	registry.Register(NewProductBatchSeeder()) // Product batches must be before ProductUnit
	registry.Register(NewProductBatchTrackSeeder())
	registry.Register(NewProductUnitSeeder()) // ProductUnit depends on Location & ProductBatch
	registry.Register(NewProductUnitTrackSeeder())
	registry.Register(NewProductStockSeeder())
	registry.Register(NewProductStockTrackSeeder())
	registry.Register(NewProductItemSeeder())
	registry.Register(NewProductItemTrackSeeder())
	registry.Register(NewFileSeeder()) // Files - depends on other models
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
