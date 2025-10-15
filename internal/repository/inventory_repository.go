package repository

import (
	"myapp/database"
)

type InventoryRepository struct{}

type InventoryStockResponse struct {
	BrandID          uint     `json:"brandId"`
	CategoryID       uint     `json:"categoryId"`
	ProductID        uint     `json:"productId"`
	LocationID       uint     `json:"locationId"`
	BrandName        string   `json:"brandName"`
	CategoryName     string   `json:"categoryName"`
	ProductName      string   `json:"productName"`
	LocationName     string   `json:"locationName"`
	Stock            *float64 `json:"stock"`
	ProductUnitPrice *float64 `json:"productUnitPrice"`
	Barcode          *string  `json:"barcode"`
	CodeBatch        string   `json:"codeBatch"`
	PurchasePrice    *float64 `json:"purchasePrice"`
	ExpDate          *string  `json:"expDate"`
}

func NewInventoryRepository() *InventoryRepository {
	return &InventoryRepository{}
}

// GetInventoryStock returns inventory stock information with filters
func (r *InventoryRepository) GetInventoryStock(brandID, categoryID, productID, productBatchID, locationID *uint, barcode *string) ([]InventoryStockResponse, error) {
	var results []InventoryStockResponse

	query := database.DB.Table("brands AS b").
		Select(`b.id AS brand_id,
			c.id AS category_id,
			p.id AS product_id,
			l.id AS location_id,
			b.name AS brand_name,
			c.name AS category_name,
			p.name AS product_name,
			l.name AS location_name,
			ps.quantity AS stock,
			pu.unit_price AS product_unit_price,
			pu.barcode AS barcode,
			pb.code_batch AS code_batch,
			pb.unit_price AS purchase_price,
			pb.exp_date AS exp_date`).
		Joins("JOIN categories AS c ON c.brand_id = b.id AND c.deleted_at IS NULL").
		Joins("JOIN products AS p ON p.category_id = c.id AND p.deleted_at IS NULL").
		Joins("JOIN product_units AS pu ON pu.product_id = p.id AND pu.deleted_at IS NULL").
		Joins("JOIN product_batches AS pb ON pb.product_id = p.id AND pb.deleted_at IS NULL").
		Joins("JOIN product_stocks AS ps ON p.id = ps.product_id AND ps.deleted_at IS NULL AND ps.location_id = pu.location_id").
		Joins("JOIN locations AS l ON ps.location_id = l.id AND l.deleted_at IS NULL").
		Where("b.deleted_at IS NULL")

	// Apply filters
	if brandID != nil {
		query = query.Where("b.id = ?", *brandID)
	}
	if categoryID != nil {
		query = query.Where("c.id = ?", *categoryID)
	}
	if productID != nil {
		query = query.Where("p.id = ?", *productID)
	}
	if productBatchID != nil {
		query = query.Where("pb.id = ?", *productBatchID)
	}
	if locationID != nil {
		query = query.Where("l.id = ?", *locationID)
	}
	if barcode != nil {
		query = query.Where("pu.barcode = ?", *barcode)
	}

	result := query.Order("b.name, c.name, p.name, l.name").Find(&results)
	return results, result.Error
}
