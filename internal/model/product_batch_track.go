package model

import (
	"time"

	"gorm.io/gorm"
)

type ProductBatchTrack struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	ProductBatchID uint   `json:"product_batch_id" gorm:"not null;index"`
	Description    string `json:"description" gorm:"type:text;not null"` // Description of what changed
	UserInst       uint   `json:"user_inst" gorm:"not null"`             // User who made the change
	UserUpdt       *uint  `json:"user_updt"`                             // For future updates (nullable)

	// Relationships
	ProductBatch ProductBatch `json:"product_batch" gorm:"foreignKey:ProductBatchID;references:ID"`
	Creator      User         `json:"creator" gorm:"foreignKey:UserInst;references:ID"`
	Updater      *User        `json:"updater,omitempty" gorm:"foreignKey:UserUpdt;references:ID"`
}

// TrackingAction constants for tracking different types of changes
const (
	TrackActionCreate = "CREATED"
	TrackActionUpdate = "UPDATED"
	TrackActionDelete = "DELETED"
)
