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

func (r *RoleRepository) DeleteRole(id uint) error {
    return database.DB.Delete(&model.Role{}, id).Error
}

func (r *RoleRepository) GetRoleByName(name string) (model.Role, error) {
    var role model.Role
    result := database.DB.Where("name = ?", name).First(&role)
    return role, result.Error
}

func (r *RoleRepository) CheckRoleExists(name string, excludeID uint) (bool, error) {
    var count int64
    query := database.DB.Model(&model.Role{}).Where("name = ?", name)
    
    if excludeID > 0 {
        query = query.Where("id != ?", excludeID)
    }
    
    result := query.Count(&count)
    return count > 0, result.Error
}