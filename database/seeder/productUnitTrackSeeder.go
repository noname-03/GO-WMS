package seeder

import (
	"log"
	"myapp/internal/model"

	"gorm.io/gorm"
)

type ProductUnitTrackSeeder struct{}

func NewProductUnitTrackSeeder() SeederInterface {
	return &ProductUnitTrackSeeder{}
}

func (s *ProductUnitTrackSeeder) GetName() string {
	return "ProductUnitTrackSeeder"
}

func (s *ProductUnitTrackSeeder) Seed(db *gorm.DB) error {
	log.Println("🌱 Running ProductUnitTrackSeeder...")

	// Check if tracking records already exist
	var existingTracks []model.ProductUnitTrack
	result := db.Find(&existingTracks)
	if result.Error != nil {
		return result.Error
	}

	if len(existingTracks) > 0 {
		log.Printf("ProductUnitTrack seeder: Found %d existing records, skipping seed", len(existingTracks))
		return nil
	}

	log.Println("ProductUnitTrack seeder: Starting seed process...")

	// Get some existing product units to create tracking records
	var productUnits []model.ProductUnit
	db.Limit(10).Find(&productUnits)

	if len(productUnits) == 0 {
		log.Println("ProductUnitTrack seeder: No product units found, creating basic tracking records only")
		return nil
	}

	// Get admin user for audit trail
	var adminUser model.User
	var userID *uint
	if err := db.Where("username = ?", "admin").First(&adminUser).Error; err == nil {
		userID = &adminUser.ID
	}

	// Create sample tracking records
	trackingRecords := []model.ProductUnitTrack{
		{
			ProductUnitID: productUnits[0].ID,
			Description:   "Initial stock entry",
			UserIns:       userID,
		},
		{
			ProductUnitID: productUnits[0].ID,
			Description:   "Sales order fulfillment",
			UserIns:       userID,
		},
		{
			ProductUnitID: productUnits[0].ID,
			Description:   "Inventory adjustment - damaged items",
			UserIns:       userID,
		},
	}

	// Add more tracking records if we have more product units
	if len(productUnits) > 1 {
		additionalTracks := []model.ProductUnitTrack{
			{
				ProductUnitID: productUnits[1].ID,
				Description:   "New stock arrival",
				UserIns:       userID,
			},
			{
				ProductUnitID: productUnits[1].ID,
				Description:   "New stock arrival",
				UserIns:       userID,
			},
		}
		trackingRecords = append(trackingRecords, additionalTracks...)
	}

	if len(productUnits) > 2 {
		moreTracks := []model.ProductUnitTrack{
			{
				ProductUnitID: productUnits[2].ID,
				Description:   "Bulk purchase arrival",
				UserIns:       userID,
			},
		}
		trackingRecords = append(trackingRecords, moreTracks...)
	}

	// Insert tracking records
	for _, track := range trackingRecords {
		result := db.Create(&track)
		if result.Error != nil {
			log.Printf("ProductUnitTrack seeder: Failed to create tracking record for ProductUnit ID %d: %v", track.ProductUnitID, result.Error)
			return result.Error
		}
	}

	log.Printf("ProductUnitTrack seeder: Successfully created %d tracking records", len(trackingRecords))
	return nil
}
