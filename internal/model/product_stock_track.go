package model

import (
	"time"

	"gorm.io/gorm"
)

type ProductStockTrack struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	ProductStockID uint `gorm:"not null" json:"product_stock_id"`
	ProductBatchID uint `gorm:"not null" json:"product_batch_id"`
	ProductID      uint `gorm:"not null" json:"product_id"`

	// Track Information
	Date        time.Time `gorm:"not null" json:"date"`
	Quantity    float64   `gorm:"not null" json:"quantity"`
	Operation   string    `gorm:"type:varchar(10);not null" json:"operation"` // Plus, Minus
	Stock       float64   `gorm:"not null" json:"stock"`
	Description *string   `gorm:"type:text" json:"description"`

	// Audit Trail Fields
	UserIns  *uint `json:"user_ins,omitempty"`
	UserUpdt *uint `json:"user_updt,omitempty"`

	// Relationships
	ProductStock ProductStock `gorm:"foreignKey:ProductStockID;constraint:OnDelete:RESTRICT" json:"product_stock"`
	ProductBatch ProductBatch `gorm:"foreignKey:ProductBatchID;constraint:OnDelete:RESTRICT" json:"product_batch"`
	Product      Product      `gorm:"foreignKey:ProductID;constraint:OnDelete:RESTRICT" json:"product"`
	InsertedBy   *User        `gorm:"foreignKey:UserIns;constraint:OnDelete:RESTRICT" json:"inserted_by,omitempty"`
	UpdatedBy    *User        `gorm:"foreignKey:UserUpdt;constraint:OnDelete:SET NULL" json:"updated_by,omitempty"`
}
