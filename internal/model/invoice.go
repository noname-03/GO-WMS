package model

import (
	"time"

	"gorm.io/gorm"
)

type Invoice struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	InvoiceNumber   string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"invoiceNumber"`
	UserID          uint           `gorm:"not null" json:"userId"`
	PurchaseOrderID *uint          `gorm:"index" json:"purchaseOrderId"`
	DeliveryOrderID *uint          `gorm:"index" json:"deliveryOrderId"`
	InvoiceDate     time.Time      `gorm:"type:date;not null" json:"invoiceDate"`
	Status          string         `gorm:"type:varchar(16);not null" json:"status"` // draft, sent, paid, closed
	TotalAmount     float64        `gorm:"type:numeric(18,2);not null" json:"totalAmount"`
	Description     *string        `gorm:"type:text" json:"description"`
	UserIns         uint           `gorm:"not null" json:"userIns"`
	UserUpdt        uint           `gorm:"not null" json:"userUpdt"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	PurchaseOrder   *PurchaseOrder `gorm:"foreignKey:PurchaseOrderID;references:ID" json:"purchaseOrder,omitempty"`
	DeliveryOrder   *DeliveryOrder `gorm:"foreignKey:DeliveryOrderID;references:ID" json:"deliveryOrder,omitempty"`
	User            *User          `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
	InvoiceItems    []InvoiceItem  `gorm:"foreignKey:InvoiceID" json:"invoiceItems,omitempty"`
}

// TableName overrides the table name
func (Invoice) TableName() string {
	return "invoices"
}
