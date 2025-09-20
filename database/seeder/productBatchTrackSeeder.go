package seeder

import (
	"log"
	"myapp/internal/model"
	"myapp/internal/utils"

	"gorm.io/gorm"
)

type ProductBatchTrackSeeder struct{}

func NewProductBatchTrackSeeder() SeederInterface {
	return &ProductBatchTrackSeeder{}
}

func (s *ProductBatchTrackSeeder) GetName() string {
	return "ProductBatchTrackSeeder"
}

func (s *ProductBatchTrackSeeder) Seed(db *gorm.DB) error {
	log.Println("🌱 Running ProductBatchTrackSeeder...")

	// Check if tracking records already exist
	var existingTracks []model.ProductBatchTrack
	result := db.Find(&existingTracks)
	if result.Error != nil {
		return result.Error
	}

	if len(existingTracks) > 0 {
		log.Printf("ProductBatchTrack seeder: Found %d existing records, skipping seed", len(existingTracks))
		return nil
	}

	log.Println("ProductBatchTrack seeder: Starting seed process...")

	// Get some existing product batches to create tracking records
	var productBatches []model.ProductBatch
	db.Limit(5).Find(&productBatches)

	if len(productBatches) == 0 {
		log.Println("ProductBatchTrack seeder: No product batches found, creating basic tracking records only")
		return nil
	}

	// Create sample tracking records
	trackingRecords := []model.ProductBatchTrack{
		{
			ProductBatchID: productBatches[0].ID,
			Description:    utils.GenerateCreateDescription(productBatches[0]),
			UserInst:       1, // Assuming admin user has ID 1
		},
	}

	// Add more tracking records if we have more product batches
	if len(productBatches) > 1 {
		trackingRecords = append(trackingRecords,
			model.ProductBatchTrack{
				ProductBatchID: productBatches[1].ID,
				Description:    utils.GenerateCreateDescription(productBatches[1]),
				UserInst:       1,
			})
	}

	if len(productBatches) > 2 {
		// Add an update tracking record
		trackingRecords = append(trackingRecords,
			model.ProductBatchTrack{
				ProductBatchID: productBatches[2].ID,
				Description:    "Product batch updated: changed unit price from 10.00 to 12.50",
				UserInst:       2, // Assuming user with ID 2 made this change
			})
	}

	if len(productBatches) > 3 {
		trackingRecords = append(trackingRecords,
			model.ProductBatchTrack{
				ProductBatchID: productBatches[3].ID,
				Description:    utils.GenerateCreateDescription(productBatches[3]),
				UserInst:       1,
			})
	}

	if len(productBatches) > 4 {
		// Add another update tracking record
		trackingRecords = append(trackingRecords,
			model.ProductBatchTrack{
				ProductBatchID: productBatches[4].ID,
				Description:    "Product batch updated: changed expiry date from 2024-12-31 to 2025-06-30",
				UserInst:       2,
			})
	}

	// Insert tracking records
	for _, track := range trackingRecords {
		result := db.Create(&track)
		if result.Error != nil {
			log.Printf("ProductBatchTrack seeder: Failed to create tracking record for ProductBatch ID %d: %v", track.ProductBatchID, result.Error)
			return result.Error
		}
	}

	log.Printf("ProductBatchTrack seeder: Successfully created %d tracking records", len(trackingRecords))
	return nil
}
