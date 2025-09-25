package model

import (
	"time"

	"gorm.io/gorm"
)

type ProductStock struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	ProductBatchID uint `gorm:"not null" json:"product_batch_id"`
	ProductID      uint `gorm:"not null" json:"product_id"`
	LocationID     uint `gorm:"not null" json:"location_id"`

	// Stock Information
	Quantity *float64 `json:"quantity"`

	// Audit Trail Fields
	UserIns  *uint `json:"user_ins,omitempty"`
	UserUpdt *uint `json:"user_updt,omitempty"`

	// Relationships
	ProductBatch ProductBatch `gorm:"foreignKey:ProductBatchID;constraint:OnDelete:RESTRICT" json:"product_batch"`
	Product      Product      `gorm:"foreignKey:ProductID;constraint:OnDelete:RESTRICT" json:"product"`
	Location     *Location    `gorm:"foreignKey:LocationID;constraint:OnDelete:RESTRICT" json:"location,omitempty"`
	InsertedBy   *User        `gorm:"foreignKey:UserIns;constraint:OnDelete:RESTRICT" json:"inserted_by,omitempty"`
	UpdatedBy    *User        `gorm:"foreignKey:UserUpdt;constraint:OnDelete:SET NULL" json:"updated_by,omitempty"`
}
