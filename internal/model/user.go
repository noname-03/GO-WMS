package model

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"` // "-" means don't include in JSON response
}

// SafeLogString returns a safe string representation of User without password
func (u *User) SafeLogString() string {
	return fmt.Sprintf("User{ID: %d, Name: %s, Email: %s, CreatedAt: %v, UpdatedAt: %v}",
		u.ID, u.Name, u.Email, u.CreatedAt, u.UpdatedAt)
}
