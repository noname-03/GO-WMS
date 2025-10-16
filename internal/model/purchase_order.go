package model

import (
	"time"

	"gorm.io/gorm"
)

type PurchaseOrder struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Purchase Order Information
	PONumber    string    `gorm:"size:64;unique;not null" json:"po_number"`
	UserID      uint      `gorm:"not null" json:"user_id"` // User who placed the order (reseller)
	OrderDate   time.Time `gorm:"not null" json:"order_date"`
	Status      string    `gorm:"size:16;not null;default:'draft'" json:"status"` // draft, submitted, approved, received, closed
	TotalAmount *float64  `gorm:"type:numeric(18,2);not null;default:0" json:"total_amount"`
	Description *string   `gorm:"type:text" json:"description"`

	// Audit Trail Fields
	UserIns  *uint `json:"user_ins,omitempty"`
	UserUpdt *uint `json:"user_updt,omitempty"`

	// Relationships
	User               User                `gorm:"foreignKey:UserID;constraint:OnDelete:RESTRICT" json:"user"`
	PurchaseOrderItems []PurchaseOrderItem `gorm:"foreignKey:PurchaseOrderID" json:"purchase_order_items,omitempty"`
	InsertedBy         *User               `gorm:"foreignKey:UserIns;constraint:OnDelete:RESTRICT" json:"inserted_by,omitempty"`
	UpdatedBy          *User               `gorm:"foreignKey:UserUpdt;constraint:OnDelete:SET NULL" json:"updated_by,omitempty"`
}
