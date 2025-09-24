package model

import (
	"time"

	"gorm.io/gorm"
)

type ProductUnit struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Key to Product
	ProductID   uint     `gorm:"not null" json:"product_id"`
	LocationID  uint     `gorm:"not null" json:"location_id"`
	Name        *string  `json:"name"`        // Nullable name
	Quantity    *float64 `json:"quantity"`    // Nullable quantity
	UnitPrice   *float64 `json:"unit_price"`  // Nullable unit price
	Barcode     *string  `json:"barcode"`     // Nullable barcode
	Description *string  `json:"description"` // Nullable description
	// Audit Trail Fields
	UserIns  *uint `json:"user_ins,omitempty"`  // Pointer untuk allow null
	UserUpdt *uint `json:"user_updt,omitempty"` // Pointer untuk allow null
	// Relationships
	Product    Product  `gorm:"foreignKey:ProductID;constraint:OnDelete:RESTRICT" json:"product"`
	Location   Location `gorm:"foreignKey:LocationID;constraint:OnDelete:RESTRICT" json:"location"`
	InsertedBy *User    `gorm:"foreignKey:UserIns;constraint:OnDelete:RESTRICT" json:"inserted_by,omitempty"`
	UpdatedBy  *User    `gorm:"foreignKey:UserUpdt;constraint:OnDelete:SET NULL" json:"updated_by,omitempty"`
}
