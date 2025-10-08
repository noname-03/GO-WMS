package model

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// File details
	ModelType string  `gorm:"not null;size:100" json:"model_type"` // e.g., 'product', 'user', 'category'
	ModelID   uint    `gorm:"not null" json:"model_id"`            // ID of the related model
	Ext       string  `gorm:"not null;size:10" json:"ext"`         // File extension (jpg, png, pdf, etc.)
	FileURL   *string `gorm:"size:500" json:"file_url"`            // S3 URL or file path (nullable)

	// Audit Trail Fields
	UserIns  *uint `json:"user_ins,omitempty"`  // Pointer untuk allow null
	UserUpdt *uint `json:"user_updt,omitempty"` // Pointer untuk allow null

	// Relationships
	InsertedBy *User `gorm:"foreignKey:UserIns;constraint:OnDelete:RESTRICT" json:"inserted_by,omitempty"`
	UpdatedBy  *User `gorm:"foreignKey:UserUpdt;constraint:OnDelete:SET NULL" json:"updated_by,omitempty"`
}

// TableName specifies the table name for the File model
func (File) TableName() string {
	return "files"
}
