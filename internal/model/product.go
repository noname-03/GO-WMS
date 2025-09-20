package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Key to Category
	CategoryID  uint    `gorm:"not null" json:"category_id"`
	Name        string  `gorm:"not null" json:"name"`
	Description *string `json:"description"` // Nullable description

	// Audit Trail Fields
	UserIns  *uint `json:"user_ins,omitempty"`  // Pointer untuk allow null
	UserUpdt *uint `json:"user_updt,omitempty"` // Pointer untuk allow null

	// Relationships
	Category   Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT" json:"category"`
	InsertedBy *User    `gorm:"foreignKey:UserIns;constraint:OnDelete:RESTRICT" json:"inserted_by,omitempty"`
	UpdatedBy  *User    `gorm:"foreignKey:UserUpdt;constraint:OnDelete:SET NULL" json:"updated_by,omitempty"`

	// Has many ProductBatch
	ProductBatches []ProductBatch `gorm:"foreignKey:ProductID" json:"product_batches,omitempty"`
}
