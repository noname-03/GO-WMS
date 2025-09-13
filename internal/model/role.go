package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string `gorm:"unique;not null" json:"name"`
	Description string `json:"description"`
}