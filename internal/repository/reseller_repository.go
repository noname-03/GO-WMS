package repository

import (
	"myapp/database"
)

type ResellerRepository struct{}

// resellerResponse struct untuk response list reseller
type resellerResponse struct {
	UserID       uint    `json:"userId"`
	Name         string  `json:"name"`
	LocationID   uint    `json:"locationId"`
	LocationName string  `json:"locationName"`
	Address      *string `json:"address"`
	PhoneNumber  *string `json:"phoneNumber"`
}

func NewResellerRepository() *ResellerRepository {
	return &ResellerRepository{}
}

func (r *ResellerRepository) GetAllResellers() ([]resellerResponse, error) {
	var resellers []resellerResponse

	result := database.DB.Table("users u").
		Select("u.id AS user_id, u.name AS name, l.id AS location_id, l.name AS location_name, l.address, l.phone_number").
		Joins("JOIN locations l ON u.id = l.user_id").
		Where("u.deleted_at IS NULL").
		Where("l.deleted_at IS NULL").
		Where("l.type = ?", "reseller").
		Order("u.name ASC").
		Find(&resellers)

	return resellers, result.Error
}
