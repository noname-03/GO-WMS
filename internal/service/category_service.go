package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		categoryRepo: repository.NewCategoryRepository(),
	}
}

// Business logic methods
func (s *CategoryService) GetAllCategories() (interface{}, error) {
	return s.categoryRepo.GetAllCategories()
}

func (s *CategoryService) GetCategoriesByBrand(brandID uint) (interface{}, error) {
	if brandID == 0 {
		return nil, errors.New("invalid brand ID")
	}

	// Check if brand exists
	brandExists, err := s.categoryRepo.CheckBrandExists(brandID)
	if err != nil {
		return nil, err
	}
	if !brandExists {
		return nil, errors.New("brand not found")
	}

	return s.categoryRepo.GetCategoriesByBrand(brandID)
}

func (s *CategoryService) GetCategoryByID(id uint) (interface{}, error) {
	category, err := s.categoryRepo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) CreateCategory(brandID uint, name string, description *string, userID uint) (interface{}, error) {
	if brandID == 0 {
		return nil, errors.New("brand ID is required")
	}

	if name == "" {
		return nil, errors.New("category name is required")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if brand exists
	brandExists, err := s.categoryRepo.CheckBrandExists(brandID)
	if err != nil {
		return nil, err
	}
	if !brandExists {
		return nil, errors.New("brand not found")
	}

	// Check if category exists for this brand
	exists, err := s.categoryRepo.CheckCategoryExists(name, brandID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("category already exists for this brand")
	}

	category := &model.Category{
		BrandID:     brandID,
		Name:        name,
		Description: description,
		UserIns:     &userID, // Set pointer to userID
	}

	err = s.categoryRepo.CreateCategory(category)
	if err != nil {
		return nil, err
	}

	// Fetch the created category with brand name (consistent format)
	createdCategory, err := s.categoryRepo.GetCategoryByID(category.ID)
	if err != nil {
		return nil, err
	}

	return createdCategory, nil
}

func (s *CategoryService) UpdateCategory(id uint, brandID uint, name string, description *string, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid category ID")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if category exists
	category, err := s.categoryRepo.GetCategoryByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// If brand ID is being changed, check if new brand exists
	if brandID != 0 && brandID != category.BrandID {
		brandExists, err := s.categoryRepo.CheckBrandExists(brandID)
		if err != nil {
			return nil, err
		}
		if !brandExists {
			return nil, errors.New("brand not found")
		}
	}

	// Use existing brand if not provided
	if brandID == 0 {
		brandID = category.BrandID
	}

	// Check if new name conflicts with existing categories for this brand
	if name != "" && (name != category.Name || brandID != category.BrandID) {
		exists, err := s.categoryRepo.CheckCategoryExists(name, brandID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("category name already in use for this brand")
		}
	}

	// Prepare update data with audit trail
	updateData := make(map[string]interface{})
	if brandID != category.BrandID {
		updateData["brand_id"] = brandID
	}
	if name != "" {
		updateData["name"] = name
	}
	if description != nil {
		updateData["description"] = description
	}
	// Always set the user who updated
	updateData["user_updt"] = userID

	err = s.categoryRepo.UpdateCategory(id, updateData)
	if err != nil {
		return nil, err
	}

	updatedCategory, err := s.categoryRepo.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}
	return updatedCategory, nil
}

func (s *CategoryService) DeleteCategory(id uint, userID uint) error {
	if id == 0 {
		return errors.New("invalid category ID")
	}

	if userID == 0 {
		return errors.New("user ID is required for audit trail")
	}

	// Check if category exists
	_, err := s.categoryRepo.GetCategoryByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	return s.categoryRepo.DeleteCategoryWithAudit(id, userID)
}
