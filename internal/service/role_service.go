package service

import (
	"errors"
	"log"
	"myapp/internal/model"
	"myapp/internal/repository"
)

type RoleService struct {
	roleRepo *repository.RoleRepository
}

func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo: repository.NewRoleRepository(),
	}
}

// Business logic methods
func (s *RoleService) GetAllRoles() ([]model.Role, error) {
	return s.roleRepo.GetAllRoles()
}

func (s *RoleService) GetRoleByID(id uint) (*model.Role, error) {
	role, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *RoleService) CreateRole(name, description string, userID uint) (*model.Role, error) {
	if name == "" {
		return nil, errors.New("role name is required")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if role exists
	exists, err := s.roleRepo.CheckRoleExists(name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("role already exists")
	}

	role := &model.Role{
		Name:        name,
		Description: description,
		UserIns:     &userID, // Set pointer to userID
	}

	err = s.roleRepo.CreateRole(role)
	return role, err
}

func (s *RoleService) UpdateRole(id uint, name, description string, userID uint) (*model.Role, error) {
	if id == 0 {
		return nil, errors.New("invalid role ID")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if role exists
	role, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return nil, errors.New("role not found")
	}

	// Check if new name conflicts with existing roles
	if name != "" && name != role.Name {
		exists, err := s.roleRepo.CheckRoleExists(name)
		if err != nil {
			return nil, err
		}
		log.Println("cek status", exists)
		if exists {
			return nil, errors.New("role name already in use")
		}
	}

	// Prepare update data with audit trail
	updateData := make(map[string]interface{})
	if name != "" {
		updateData["name"] = name
	}
	if description != "" {
		updateData["description"] = description
	}
	// Always set the user who updated
	updateData["user_updt"] = userID

	err = s.roleRepo.UpdateRole(id, updateData)
	if err != nil {
		return nil, err
	}

	updatedRole, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return nil, err
	}
	return &updatedRole, nil
}

func (s *RoleService) DeleteRole(id uint, userID uint) error {
	if id == 0 {
		return errors.New("invalid role ID")
	}

	if userID == 0 {
		return errors.New("user ID is required for audit trail")
	}

	// Check if role exists
	_, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return errors.New("role not found")
	}

	return s.roleRepo.DeleteRoleWithAudit(id, userID)
}

// GetDeletedRoles returns all soft deleted roles
func (s *RoleService) GetDeletedRoles() ([]model.Role, error) {
	return s.roleRepo.GetDeletedRoles()
}

// RestoreRole restores a soft deleted role
func (s *RoleService) RestoreRole(id uint, userID uint) (*model.Role, error) {
	if id == 0 {
		return nil, errors.New("invalid role ID")
	}
	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	err := s.roleRepo.RestoreRole(id, userID)
	if err != nil {
		return nil, err
	}

	restoredRole, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return nil, err
	}
	return &restoredRole, nil
}
