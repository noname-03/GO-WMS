package seeder

import (
	"log"
	"myapp/internal/model"
	"time"

	"gorm.io/gorm"
)

type ProductStockTrackSeeder struct{}

func NewProductStockTrackSeeder() SeederInterface {
	return &ProductStockTrackSeeder{}
}

func (s *ProductStockTrackSeeder) GetName() string {
	return "ProductStockTrackSeeder"
}

func (s *ProductStockTrackSeeder) Seed(db *gorm.DB) error {
	log.Println("🌱 Running ProductStockTrackSeeder...")

	// Get first user for audit trail
	var user model.User
	if err := db.First(&user).Error; err != nil {
		log.Printf("❌ Error getting user for ProductStockTrackSeeder: %v", err)
		return err
	}

	// Get specific product stocks that we know exist from ProductStockSeeder
	var stocks []model.ProductStock
	if err := db.Preload("Product").Preload("Location").Limit(8).Find(&stocks).Error; err != nil {
		log.Printf("❌ Error getting product stocks for ProductStockTrackSeeder: %v", err)
		return err
	}

	if len(stocks) == 0 {
		log.Println("⚠️ Skipping ProductStockTrackSeeder: No product stocks found")
		return nil
	}

	// Sample track data with realistic movements
	descriptions := []string{
		"Stock receipt from supplier",
		"Stock adjustment - count correction",
		"Sales transaction",
		"Return to supplier",
		"Transfer to other location",
		"Damaged goods write-off",
		"Sample allocation",
	}

	var tracks []model.ProductStockTrack
	now := time.Now()

	// Create meaningful tracks for each stock
	trackIndex := 0
	for _, stock := range stocks {
		// Start with current stock quantity
		currentStock := 0.0
		if stock.Quantity != nil {
			currentStock = *stock.Quantity
		}

		// Create 4-6 historical tracks leading to current stock
		numTracks := 4 + (trackIndex % 3) // 4-6 tracks per stock

		for j := 0; j < numTracks; j++ {
			// Create progressive quantities that build up to current stock
			var quantity, newStock float64
			var operation string

			if j < numTracks-1 {
				// Historical entries - mostly Plus operations to build stock
				if j%3 == 0 && currentStock > 30 {
					// Occasional Minus operation
					operation = "Minus"
					quantity = 10.0 + float64(j)*5.0
					newStock = currentStock - quantity
				} else {
					// Plus operation
					operation = "Plus"
					quantity = 20.0 + float64(j)*10.0
					newStock = currentStock + quantity
				}
			} else {
				// Final entry to match current stock
				if currentStock > 100 {
					operation = "Minus"
					quantity = 15.0
					newStock = currentStock
				} else {
					operation = "Plus"
					quantity = 25.0
					newStock = currentStock
				}
			}

			// Create track record with dates spread over last 2 months
			trackDate := now.AddDate(0, 0, -60+(j*10))
			description := descriptions[trackIndex%len(descriptions)]

			track := model.ProductStockTrack{
				ProductStockID: stock.ID,
				ProductID:      stock.ProductID,
				ProductBatchID: stock.ProductBatchID,
				Date:           trackDate,
				Quantity:       quantity,
				Operation:      operation,
				Stock:          newStock,
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
		var existing model.ProductStockTrack
		result := db.Where("product_stock_id = ? AND date = ? AND quantity = ?",
			track.ProductStockID, track.Date, track.Quantity).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			// Track doesn't exist, create it
			if err := db.Create(&track).Error; err != nil {
				log.Printf("⚠️  ProductStockTrackSeeder: Failed to create track for stock ID %d (might be duplicate), skipping: %v", track.ProductStockID, err)
				continue
			}
			log.Printf("✅ Product stock track created for stock ID %d - Operation: %s, Qty: %.2f",
				track.ProductStockID, track.Operation, track.Quantity)
		} else if result.Error != nil {
			log.Printf("⚠️  ProductStockTrackSeeder: Database error for stock ID %d, skipping: %v", track.ProductStockID, result.Error)
			continue
		} else {
			log.Printf("✅ ProductStockTrackSeeder: Track for stock ID %d on %s already exists, skipping...",
				track.ProductStockID, track.Date.Format("2006-01-02"))
		}
	}

	return nil
}
