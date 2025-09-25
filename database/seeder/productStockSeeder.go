package seeder

import (
	"log"
	"myapp/internal/model"

	"gorm.io/gorm"
)

type ProductStockSeeder struct{}

func NewProductStockSeeder() SeederInterface {
	return &ProductStockSeeder{}
}

func (s *ProductStockSeeder) GetName() string {
	return "ProductStockSeeder"
}

func (s *ProductStockSeeder) Seed(db *gorm.DB) error {
	log.Println("🌱 Running ProductStockSeeder...")

	// Get first user for audit trail
	var user model.User
	if err := db.First(&user).Error; err != nil {
		log.Printf("❌ Error getting user for ProductStockSeeder: %v", err)
		return err
	}

	// Get some specific products and batches by name for consistent seeding
	var camry, galaxy, airZoom, iphone15 model.Product
	db.Where("name LIKE ?", "%Camry%").First(&camry)
	db.Where("name LIKE ?", "%Galaxy S24%").First(&galaxy)
	db.Where("name LIKE ?", "%Air Zoom%").First(&airZoom)
	db.Where("name LIKE ?", "%iPhone 15 Pro%").First(&iphone15)

	// Get their corresponding batches
	var camryBatch, galaxyBatch, airZoomBatch, iphoneBatch model.ProductBatch
	db.Where("product_id = ?", camry.ID).First(&camryBatch)
	db.Where("product_id = ?", galaxy.ID).First(&galaxyBatch)
	db.Where("product_id = ?", airZoom.ID).First(&airZoomBatch)
	db.Where("product_id = ?", iphone15.ID).First(&iphoneBatch)

	// Get locations (use actual location names from locationSeeder)
	var warehouse, store, outlet model.Location
	db.Where("name LIKE ?", "%Gudang Pusat%").First(&warehouse)
	db.Where("name LIKE ?", "%Toko Sinar Jaya%").First(&store)
	db.Where("name LIKE ?", "%Minimarket%").First(&outlet)

	// Validate that we have valid locations
	if warehouse.ID == 0 || store.ID == 0 || outlet.ID == 0 {
		log.Printf("⚠️ Required locations not found - Warehouse ID: %d, Store ID: %d, Outlet ID: %d", warehouse.ID, store.ID, outlet.ID)
		log.Println("Make sure LocationSeeder has run first")
		return nil // Skip without error to avoid breaking the seeding process
	}

	// Define stock data with meaningful quantities
	qty1 := 150.0 // Toyota Camry at Warehouse
	qty2 := 25.0  // Toyota Camry at Store
	qty3 := 200.0 // Samsung Galaxy at Warehouse
	qty4 := 75.0  // Samsung Galaxy at Store
	qty5 := 100.0 // Nike shoes at Warehouse
	qty6 := 50.0  // Nike shoes at Outlet
	qty7 := 120.0 // iPhone at Warehouse
	qty8 := 30.0  // iPhone at Store

	stocks := []model.ProductStock{
		// Toyota Camry stock - multiple locations
		{ProductBatchID: camryBatch.ID, ProductID: camry.ID, LocationID: warehouse.ID, Quantity: &qty1, UserIns: &user.ID, UserUpdt: &user.ID},
		{ProductBatchID: camryBatch.ID, ProductID: camry.ID, LocationID: store.ID, Quantity: &qty2, UserIns: &user.ID, UserUpdt: &user.ID},

		// Samsung Galaxy stock - multiple locations
		{ProductBatchID: galaxyBatch.ID, ProductID: galaxy.ID, LocationID: warehouse.ID, Quantity: &qty3, UserIns: &user.ID, UserUpdt: &user.ID},
		{ProductBatchID: galaxyBatch.ID, ProductID: galaxy.ID, LocationID: store.ID, Quantity: &qty4, UserIns: &user.ID, UserUpdt: &user.ID},

		// Nike shoes stock - multiple locations
		{ProductBatchID: airZoomBatch.ID, ProductID: airZoom.ID, LocationID: warehouse.ID, Quantity: &qty5, UserIns: &user.ID, UserUpdt: &user.ID},
		{ProductBatchID: airZoomBatch.ID, ProductID: airZoom.ID, LocationID: outlet.ID, Quantity: &qty6, UserIns: &user.ID, UserUpdt: &user.ID},

		// iPhone stock - multiple locations
		{ProductBatchID: iphoneBatch.ID, ProductID: iphone15.ID, LocationID: warehouse.ID, Quantity: &qty7, UserIns: &user.ID, UserUpdt: &user.ID},
		{ProductBatchID: iphoneBatch.ID, ProductID: iphone15.ID, LocationID: store.ID, Quantity: &qty8, UserIns: &user.ID, UserUpdt: &user.ID},
	}

	for _, stock := range stocks {
		// Skip if required relationships don't exist
		if stock.ProductBatchID == 0 || stock.ProductID == 0 || stock.LocationID == 0 {
			log.Printf("⚠️ Skipping stock creation - missing product, batch, or location relationship (BatchID: %d, ProductID: %d, LocationID: %d)",
				stock.ProductBatchID, stock.ProductID, stock.LocationID)
			continue
		}

		var existing model.ProductStock
		result := db.Where("product_batch_id = ? AND product_id = ? AND location_id = ?",
			stock.ProductBatchID, stock.ProductID, stock.LocationID).First(&existing)

		if result.Error != nil {
			// ProductStock doesn't exist, create it
			if err := db.Create(&stock).Error; err != nil {
				log.Printf("❌ Failed to create product stock for product ID %d batch ID %d: %v",
					stock.ProductID, stock.ProductBatchID, err)
				return err
			}
			log.Printf("✅ Product stock created for product ID %d at location ID %v - Qty: %.2f",
				stock.ProductID, stock.LocationID, *stock.Quantity)
		} else {
			log.Printf("✅ ProductStockSeeder: Stock for product ID %d batch ID %d at location ID %v already exists, skipping...",
				stock.ProductID, stock.ProductBatchID, stock.LocationID)
		}
	}

	return nil
}
