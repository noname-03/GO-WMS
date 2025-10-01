package seeder

import (
	"log"
	"myapp/internal/model"
	"time"

	"gorm.io/gorm"
)

type ProductItemTrackSeeder struct{}

func NewProductItemTrackSeeder() SeederInterface {
	return &ProductItemTrackSeeder{}
}

func (s *ProductItemTrackSeeder) GetName() string {
	return "ProductItemTrackSeeder"
}

func (s *ProductItemTrackSeeder) Seed(db *gorm.DB) error {
	log.Println("🌱 Running ProductItemTrackSeeder...")

	// Get first user for audit trail
	var user model.User
	if err := db.First(&user).Error; err != nil {
		log.Printf("❌ Error getting user for ProductItemTrackSeeder: %v", err)
		return err
	}

	// Get product items that we know exist from ProductItemSeeder
	var items []model.ProductItem
	if err := db.Preload("ProductStock").Limit(10).Find(&items).Error; err != nil {
		log.Printf("❌ Error getting product items for ProductItemTrackSeeder: %v", err)
		return err
	}

	if len(items) == 0 {
		log.Println("⚠️ Skipping ProductItemTrackSeeder: No product items found")
		return nil
	}

	// Sample track data with realistic movements
	unitPrices := []string{"15000.00", "25000.00", "50000.00", "12000.00", "35000.00", "8000.00", "18000.00", "45000.00"}
	descriptions := []string{
		"Item received from supplier",
		"Item sold to customer",
		"Inventory adjustment",
		"Quality control check",
		"Transfer between locations",
		"Return processing",
		"Damage assessment",
	}

	var tracks []model.ProductItemTrack

	// Create meaningful tracks for each item
	trackIndex := 0
	now := time.Now()

	for _, item := range items {
		// Start with current item quantity
		currentStock := 0.0
		if item.Quantity != nil {
			currentStock = *item.Quantity
		}

		// Create 3-4 historical tracks leading to current stock
		numTracks := 3 + (trackIndex % 2) // 3-4 tracks per item

		for j := 0; j < numTracks; j++ {
			// Create progressive dates (spread over last 30 days)
			daysAgo := 30 - (j * 7) // 30, 23, 16, 9 days ago
			trackDate := now.AddDate(0, 0, -daysAgo)

			// Create progressive quantities that build to current stock
			var quantity, newStock float64
			var operation string

			if j < numTracks-1 {
				// Historical entries - mostly Plus operations to build stock
				if j%2 == 1 && currentStock > 20 {
					// Occasional Minus operation
					operation = "Minus"
					quantity = 5.0 + float64(j)*3.0
					newStock = currentStock - quantity
				} else {
					// Plus operation
					operation = "Plus"
					quantity = 10.0 + float64(j)*5.0
					newStock = currentStock + quantity
				}
			} else {
				// Final entry to match current stock
				if currentStock > 50 {
					operation = "Minus"
					quantity = 8.0
					newStock = currentStock
				} else {
					operation = "Plus"
					quantity = 12.0
					newStock = currentStock
				}
			}

			// Get unit price and description
			unitPrice := unitPrices[trackIndex%len(unitPrices)]
			description := descriptions[trackIndex%len(descriptions)]

			track := model.ProductItemTrack{
				ProductStockID: item.ProductStockID,
				ProductID:      item.ProductID,
				ProductBatchID: item.ProductBatchID,
				Date:           trackDate,
				Quantity:       quantity,
				Operation:      operation,
				Stock:          newStock,
				UnitPrice:      &unitPrice,
				Description:    &description,
				UserIns:        &user.ID,
				UserUpdt:       &user.ID,
			}

			tracks = append(tracks, track)
			trackIndex++
		}
	}

	// Create tracks individually to handle duplicates properly
	for _, track := range tracks {
		var existing model.ProductItemTrack
		result := db.Where("product_stock_id = ? AND date = ? AND quantity = ? AND operation = ?",
			track.ProductStockID, track.Date, track.Quantity, track.Operation).First(&existing)

		if result.Error != nil {
			// Track doesn't exist, create it
			if err := db.Create(&track).Error; err != nil {
				log.Printf("❌ Failed to create product item track for stock ID %d: %v", track.ProductStockID, err)
				return err
			}
			log.Printf("✅ Product item track created for stock ID %d - Date: %s, Operation: %s, Qty: %.2f, Stock: %.2f",
				track.ProductStockID, track.Date.Format("2006-01-02"), track.Operation, track.Quantity, track.Stock)
		} else {
			log.Printf("✅ ProductItemTrackSeeder: Track for stock ID %d on %s already exists, skipping...",
				track.ProductStockID, track.Date.Format("2006-01-02"))
		}
	}

	return nil
}
