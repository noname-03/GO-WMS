package service

import (
	"errors"
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

func (s *RoleService) CreateRole(name, description string) (*model.Role, error) {
	if name == "" {
		return nil, errors.New("role name is required")
	}
	
	// Check if role exists
	exists, err := s.roleRepo.CheckRoleExists(name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("role already exists")
	}

	role := &model.Role{
		Name:        name,
		Description: description,
	}

	err = s.roleRepo.CreateRole(role)
	return role, err
}

func (s *RoleService) UpdateRole(id uint, name, description string) (*model.Role, error) {
	if id == 0 {
		return nil, errors.New("invalid role ID")
	}

	// Check if role exists
	role, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return nil, errors.New("role not found")
	}

	// Check if new name conflicts with existing roles
	if name != "" && name != role.Name {
		exists, err := s.roleRepo.CheckRoleExists(name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("role name already in use")
		}
	}

	// Prepare update data
	updateData := make(map[string]interface{})
	if name != "" {
		updateData["name"] = name
	}
	if description != "" {
		updateData["description"] = description
	}

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

func (s *RoleService) DeleteRole(id uint) error {
	if id == 0 {
		return errors.New("invalid role ID")
	}
	
	// Check if role exists
	_, err := s.roleRepo.GetRoleByID(id)
	if err != nil {
		return errors.New("role not found")
	}

	return s.roleRepo.DeleteRole(id)
}