package model

import (
	"time"

	"gorm.io/gorm"
)

type Location struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Foreign Key to User
	UserID      uint    `gorm:"not null" json:"user_id"`
	Name        string  `gorm:"type:varchar(100);not null" json:"name"`
	Address     *string `gorm:"type:text" json:"address"`
	PhoneNumber *string `gorm:"type:varchar(20)" json:"phone_number"`
	Type        string  `gorm:"type:varchar(20);not null;check:type IN ('gudang', 'reseller')" json:"type"`

	// Audit Trail Fields
	UserIns  *uint `json:"user_ins,omitempty"`  // Pointer untuk allow null
	UserUpdt *uint `json:"user_updt,omitempty"` // Pointer untuk allow null

	// Relationships
	User       User  `gorm:"foreignKey:UserID;constraint:OnDelete:RESTRICT" json:"user"`
	InsertedBy *User `gorm:"foreignKey:UserIns;constraint:OnDelete:RESTRICT" json:"inserted_by,omitempty"`
	UpdatedBy  *User `gorm:"foreignKey:UserUpdt;constraint:OnDelete:SET NULL" json:"updated_by,omitempty"`
}
