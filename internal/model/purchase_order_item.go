package model

import (
	"time"

	"gorm.io/gorm"
)

type PurchaseOrderItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	PurchaseOrderID uint `gorm:"not null" json:"purchase_order_id"`
	ProductID       uint `gorm:"not null" json:"product_id"`

	// Purchase Order Item Information
	QtyOrdered  *float64 `gorm:"type:numeric(18,2);not null" json:"qty_ordered"`
	UnitPrice   *float64 `gorm:"type:numeric(18,2);not null" json:"unit_price"`
	Discount    *float64 `gorm:"type:numeric(18,2);default:0" json:"discount"`
	TotalPrice  *float64 `gorm:"type:numeric(18,2);not null" json:"total_price"`
	Description *string  `gorm:"type:text" json:"description"`

	// Audit Trail Fields
	UserIns  *uint `json:"user_ins,omitempty"`
	UserUpdt *uint `json:"user_updt,omitempty"`

	// Relationships
	PurchaseOrder PurchaseOrder `gorm:"foreignKey:PurchaseOrderID;constraint:OnDelete:CASCADE" json:"purchase_order"`
	Product       Product       `gorm:"foreignKey:ProductID;constraint:OnDelete:RESTRICT" json:"product"`
	InsertedBy    *User         `gorm:"foreignKey:UserIns;constraint:OnDelete:RESTRICT" json:"inserted_by,omitempty"`
	UpdatedBy     *User         `gorm:"foreignKey:UserUpdt;constraint:OnDelete:SET NULL" json:"updated_by,omitempty"`
}
