package model

import (
	"time"

	"gorm.io/gorm"
)

type DeliveryOrder struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Delivery Order Information
	DONumber        string    `gorm:"size:64;unique;not null" json:"do_number"`
	PurchaseOrderID uint      `gorm:"not null" json:"purchase_order_id"`
	DeliveryDate    time.Time `gorm:"not null" json:"delivery_date"`
	Status          string    `gorm:"size:16;not null;default:'draft'" json:"status"` // draft, shipped, received, closed
	Description     *string   `gorm:"type:text" json:"description"`

	// Audit Trail Fields
	UserIns  *uint `json:"user_ins,omitempty"`
	UserUpdt *uint `json:"user_updt,omitempty"`

	// Relationships
	PurchaseOrder      PurchaseOrder       `gorm:"foreignKey:PurchaseOrderID;constraint:OnDelete:RESTRICT" json:"purchase_order"`
	DeliveryOrderItems []DeliveryOrderItem `gorm:"foreignKey:DeliveryOrderID" json:"delivery_order_items,omitempty"`
	InsertedBy         *User               `gorm:"foreignKey:UserIns;constraint:OnDelete:RESTRICT" json:"inserted_by,omitempty"`
	UpdatedBy          *User               `gorm:"foreignKey:UserUpdt;constraint:OnDelete:SET NULL" json:"updated_by,omitempty"`
}
