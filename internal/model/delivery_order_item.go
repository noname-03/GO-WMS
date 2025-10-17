package model

import (
	"time"

	"gorm.io/gorm"
)

type DeliveryOrderItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Keys
	DeliveryOrderID uint `gorm:"not null" json:"delivery_order_id"`
	ProductID       uint `gorm:"not null" json:"product_id"`

	// Delivery Order Item Information
	QtyDelivered *float64 `gorm:"type:numeric(18,2);not null" json:"qty_delivered"`
	Description  *string  `gorm:"type:text" json:"description"`

	// Audit Trail Fields
	UserIns  *uint `json:"user_ins,omitempty"`
	UserUpdt *uint `json:"user_updt,omitempty"`

	// Relationships
	DeliveryOrder DeliveryOrder `gorm:"foreignKey:DeliveryOrderID;constraint:OnDelete:CASCADE" json:"delivery_order"`
	Product       Product       `gorm:"foreignKey:ProductID;constraint:OnDelete:RESTRICT" json:"product"`
	InsertedBy    *User         `gorm:"foreignKey:UserIns;constraint:OnDelete:RESTRICT" json:"inserted_by,omitempty"`
	UpdatedBy     *User         `gorm:"foreignKey:UserUpdt;constraint:OnDelete:SET NULL" json:"updated_by,omitempty"`
}
