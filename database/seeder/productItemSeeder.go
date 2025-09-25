package seeder

import (
	"log"
	"myapp/internal/model"

	"gorm.io/gorm"
)

type ProductItemSeeder struct{}

func NewProductItemSeeder() SeederInterface {
	return &ProductItemSeeder{}
}

func (s *ProductItemSeeder) GetName() string {
	return "ProductItemSeeder"
}

func (s *ProductItemSeeder) Seed(db *gorm.DB) error {
	log.Println("🌱 Running ProductItemSeeder...")

	// Get first user for audit trail
	var user model.User
	if err := db.First(&user).Error; err != nil {
		log.Printf("❌ Error getting user for ProductItemSeeder: %v", err)
		return err
	}

	// Get product stocks that we know exist from ProductStockSeeder
	var stocks []model.ProductStock
	if err := db.Preload("Product").Preload("Location").Limit(8).Find(&stocks).Error; err != nil {
		log.Printf("❌ Error getting product stocks for ProductItemSeeder: %v", err)
		return err
	}

	if len(stocks) == 0 {
		log.Println("⚠️ Skipping ProductItemSeeder: No product stocks found")
		return nil
	}

	// Sample item descriptions
	descriptions := []string{
		"Initial stock receipt",
		"Stock adjustment after inventory count",
		"Partial shipment from supplier",
		"Transfer from main warehouse",
		"Customer return processing",
		"Quality control sample",
	}

	var items []model.ProductItem

	// Create meaningful items for each stock
	itemIndex := 0
	for i, stock := range stocks {
		// Create 2-3 items per stock with different scenarios
		numItems := 2 + (i % 2) // 2-3 items per stock

		for j := 0; j < numItems; j++ {
			var stockIn, stockOut, quantity *float64
			description := descriptions[itemIndex%len(descriptions)]

			// Create different realistic scenarios
			scenario := (itemIndex + j) % 4
			switch scenario {
			case 0: // Pure stock in - receiving goods
				stockInVal := 50.0 + float64(itemIndex*10)
				stockIn = &stockInVal
				quantity = &stockInVal // Quantity equals stock in when no stock out

			case 1: // Stock in with partial out - receiving then selling some
				stockInVal := 80.0 + float64(itemIndex*8)
				stockOutVal := 20.0 + float64(j*5)
				if stockOutVal > stockInVal {
					stockOutVal = stockInVal - 10 // Ensure we don't go negative
				}
				stockIn = &stockInVal
				stockOut = &stockOutVal
				quantityVal := stockInVal - stockOutVal
				quantity = &quantityVal

			case 2: // Large stock in - bulk receiving
				stockInVal := 100.0 + float64(itemIndex*15)
				stockIn = &stockInVal
				quantity = &stockInVal

			case 3: // Stock in with significant out - high turnover
				stockInVal := 120.0 + float64(itemIndex*12)
				stockOutVal := 40.0 + float64(j*8)
				if stockOutVal > stockInVal {
					stockOutVal = stockInVal - 25
				}
				stockIn = &stockInVal
				stockOut = &stockOutVal
				quantityVal := stockInVal - stockOutVal
				quantity = &quantityVal
			}

			item := model.ProductItem{
				ProductStockID: stock.ID,
				ProductID:      stock.ProductID,
				ProductBatchID: stock.ProductBatchID,
				StockIn:        stockIn,
				StockOut:       stockOut,
				Quantity:       quantity,
				Description:    &description,
				UserIns:        &user.ID,
				UserUpdt:       &user.ID,
			}

			items = append(items, item)
			itemIndex++
		}
	}

	// Create items individually to handle duplicates properly
	for _, item := range items {
		var existing model.ProductItem
		result := db.Where("product_stock_id = ? AND stock_in = ? AND description = ?",
			item.ProductStockID, item.StockIn, item.Description).First(&existing)

		if result.Error != nil {
			// Item doesn't exist, create it
			if err := db.Create(&item).Error; err != nil {
				log.Printf("❌ Failed to create product item for stock ID %d: %v", item.ProductStockID, err)
				return err
			}

			stockInVal := 0.0
			stockOutVal := 0.0
			quantityVal := 0.0

			if item.StockIn != nil {
				stockInVal = *item.StockIn
			}
			if item.StockOut != nil {
				stockOutVal = *item.StockOut
			}
			if item.Quantity != nil {
				quantityVal = *item.Quantity
			}

			log.Printf("✅ Product item created for stock ID %d - In: %.2f, Out: %.2f, Qty: %.2f",
				item.ProductStockID, stockInVal, stockOutVal, quantityVal)
		} else {
			log.Printf("✅ ProductItemSeeder: Item for stock ID %d already exists, skipping...", item.ProductStockID)
		}
	}

	return nil
}
