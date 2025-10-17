package model

import (
	"time"

	"gorm.io/gorm"
)

type InvoiceItem struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	InvoiceID   uint           `gorm:"not null;index" json:"invoiceId"`
	ProductID   uint           `gorm:"not null;index" json:"productId"`
	QtyInvoiced float64        `gorm:"type:numeric(18,2);not null" json:"qtyInvoiced"`
	UnitPrice   float64        `gorm:"type:numeric(18,2);not null" json:"unitPrice"`
	TotalPrice  float64        `gorm:"type:numeric(18,2);not null" json:"totalPrice"`
	Description *string        `gorm:"type:text" json:"description"`
	UserIns     uint           `gorm:"not null" json:"userIns"`
	UserUpdt    uint           `gorm:"not null" json:"userUpdt"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	Invoice     *Invoice       `gorm:"foreignKey:InvoiceID;references:ID" json:"invoice,omitempty"`
	Product     *Product       `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
}

// TableName overrides the table name
func (InvoiceItem) TableName() string {
	return "invoice_items"
}
