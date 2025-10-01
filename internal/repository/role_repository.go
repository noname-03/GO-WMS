package repository

import (
	"myapp/database"
	"myapp/internal/model"
)

type RoleRepository struct{}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{}
}

func (r *RoleRepository) GetAllRoles() ([]model.Role, error) {
	var roles []model.Role
	result := database.DB.Find(&roles)
	return roles, result.Error
}

func (r *RoleRepository) GetRoleByID(id uint) (model.Role, error) {
	var role model.Role
	result := database.DB.First(&role, id)
	return role, result.Error
}

func (r *RoleRepository) CreateRole(role *model.Role) error {
	return database.DB.Create(role).Error
}

func (r *RoleRepository) UpdateRole(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.Role{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *RoleRepository) DeleteRoleWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the role
	updateData := map[string]interface{}{
		"user_updt": userID,
	}

	// Update the audit field first
	err := database.DB.Model(&model.Role{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.Role{}, id).Error
}

func (r *RoleRepository) CheckRoleExists(name string) (bool, error) {
	var count int64

	// Log the query parameters
	// log.Printf("CheckRoleExists - Searching for role name: %s", name)

	query := database.DB.Model(&model.Role{}).Unscoped().Where("name ILIKE ?", name)

	// Enable debug mode to see the actual SQL query
	// result := query.Debug().Count(&count)
	result := query.Count(&count)
	return count > 0, result.Error
}

// GetDeletedRoles returns all soft deleted roles
func (r *RoleRepository) GetDeletedRoles() ([]model.Role, error) {
	var roles []model.Role
	result := database.DB.Unscoped().Where("deleted_at IS NOT NULL").Order("deleted_at DESC").Find(&roles)
	return roles, result.Error
}

// RestoreRole restores a soft deleted role
func (r *RoleRepository) RestoreRole(id uint, userID uint) error {
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"deleted_at": nil,
	}
	return database.DB.Unscoped().Model(&model.Role{}).Where("id = ?", id).Updates(updateData).Error
}
