package model

import (
	"time"

	"gorm.io/gorm"
)

type ProductUnitTrack struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Key to Product Unit
	ProductUnitID uint `gorm:"not null" json:"product_unit_id"`

	Description string `gorm:"type:text;not null" json:"description"`

	// Audit Trail Fields
	UserIns  *uint `json:"user_ins,omitempty"`  // Pointer untuk allow null
	UserUpdt *uint `json:"user_updt,omitempty"` // Pointer untuk allow null

	// Relationships
	ProductUnit ProductUnit `gorm:"foreignKey:ProductUnitID;constraint:OnDelete:RESTRICT" json:"product_unit"`
	InsertedBy  *User       `gorm:"foreignKey:UserIns;constraint:OnDelete:RESTRICT" json:"inserted_by,omitempty"`
	UpdatedBy   *User       `gorm:"foreignKey:UserUpdt;constraint:OnDelete:SET NULL" json:"updated_by,omitempty"`
}
