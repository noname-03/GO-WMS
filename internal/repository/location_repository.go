package repository

import (
	"myapp/database"
	"myapp/internal/model"
)

type LocationRepository struct{}

// locationWithUserResponse struct untuk response dengan user name
type locationWithUserResponse struct {
	ID          uint    `json:"id"`
	UserID      uint    `json:"userId"`
	UserName    string  `json:"userName"`
	Name        string  `json:"name"`
	Address     *string `json:"address"`
	PhoneNumber *string `json:"phoneNumber"`
	Type        string  `json:"type"`
}

func NewLocationRepository() *LocationRepository {
	return &LocationRepository{}
}

func (r *LocationRepository) GetAllLocations() ([]locationWithUserResponse, error) {
	var locations []locationWithUserResponse

	result := database.DB.Table("locations l").
		Select("l.id, l.user_id, u.name as user_name, l.name, l.address, l.phone_number, l.type").
		Joins("INNER JOIN users u ON l.user_id = u.id AND u.deleted_at IS NULL").
		Where("l.deleted_at IS NULL").
		Order("l.name ASC").
		Find(&locations)

	return locations, result.Error
}

func (r *LocationRepository) GetLocationsByUser(userID uint) ([]locationWithUserResponse, error) {
	var locations []locationWithUserResponse

	result := database.DB.Table("locations l").
		Select("l.id, l.user_id, u.name as user_name, l.name, l.address, l.phone_number, l.type").
		Joins("INNER JOIN users u ON l.user_id = u.id AND u.deleted_at IS NULL").
		Where("l.user_id = ? AND l.deleted_at IS NULL", userID).
		Order("l.name ASC").
		Find(&locations)

	return locations, result.Error
}

func (r *LocationRepository) GetLocationsByType(locationType string) ([]locationWithUserResponse, error) {
	var locations []locationWithUserResponse

	result := database.DB.Table("locations l").
		Select("l.id, l.user_id, u.name as user_name, l.name, l.address, l.phone_number, l.type").
		Joins("INNER JOIN users u ON l.user_id = u.id AND u.deleted_at IS NULL").
		Where("l.type = ? AND l.deleted_at IS NULL", locationType).
		Order("l.name ASC").
		Find(&locations)

	return locations, result.Error
}

func (r *LocationRepository) GetLocationByID(id uint) (locationWithUserResponse, error) {
	var location locationWithUserResponse

	result := database.DB.Table("locations l").
		Select("l.id, l.user_id, u.name as user_name, l.name, l.address, l.phone_number, l.type").
		Joins("INNER JOIN users u ON l.user_id = u.id AND u.deleted_at IS NULL").
		Where("l.id = ? AND l.deleted_at IS NULL", id).
		First(&location)

	return location, result.Error
}

// GetLocationModelByID returns model.Location for service operations
func (r *LocationRepository) GetLocationModelByID(id uint) (model.Location, error) {
	var location model.Location
	result := database.DB.Where("id = ?", id).First(&location)
	return location, result.Error
}

func (r *LocationRepository) CreateLocation(location *model.Location) error {
	return database.DB.Create(location).Error
}

func (r *LocationRepository) UpdateLocation(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.Location{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *LocationRepository) DeleteLocationWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the location
	updateData := map[string]interface{}{
		"user_updt": userID,
	}

	// Update the audit field first
	err := database.DB.Model(&model.Location{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.Location{}, id).Error
}

func (r *LocationRepository) CheckUserExists(userID uint) (bool, error) {
	var count int64
	result := database.DB.Model(&model.User{}).Where("id = ?", userID).Count(&count)
	return count > 0, result.Error
}

func (r *LocationRepository) CheckLocationNameExists(userID uint, name string, excludeID uint) (bool, error) {
	var count int64
	query := database.DB.Model(&model.Location{}).
		Where("user_id = ? AND name = ?", userID, name)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	result := query.Count(&count)
	return count > 0, result.Error
}

// GetDeletedLocations returns all soft deleted locations
func (r *LocationRepository) GetDeletedLocations() ([]locationWithUserResponse, error) {
	var locations []locationWithUserResponse

	result := database.DB.Table("locations l").
		Select("l.id, l.user_id, u.name as user_name, l.name, l.address, l.phone_number, l.type").
		Joins("INNER JOIN users u ON l.user_id = u.id").
		Where("l.deleted_at IS NOT NULL").
		Order("l.deleted_at DESC").
		Find(&locations)

	return locations, result.Error
}

// RestoreLocation restores a soft deleted location
func (r *LocationRepository) RestoreLocation(id uint, userID uint) error {
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"deleted_at": nil,
	}
	return database.DB.Unscoped().Model(&model.Location{}).Where("id = ?", id).Updates(updateData).Error
}
