package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
)

type BrandService struct {
	brandRepo *repository.BrandRepository
}

func NewBrandService() *BrandService {
	return &BrandService{
		brandRepo: repository.NewBrandRepository(),
	}
}

// Business logic methods
func (s *BrandService) GetAllBrands() ([]model.Brand, error) {
	return s.brandRepo.GetAllBrands()
}

func (s *BrandService) GetBrandByID(id uint) (*model.Brand, error) {
	brand, err := s.brandRepo.GetBrandByID(id)
	if err != nil {
		return nil, err
	}
	return &brand, nil
}

func (s *BrandService) CreateBrand(name string, description *string, userID uint) (*model.Brand, error) {
	if name == "" {
		return nil, errors.New("brand name is required")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if brand exists
	exists, err := s.brandRepo.CheckBrandExists(name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("brand already exists")
	}

	brand := &model.Brand{
		Name:        name,
		Description: description,
		UserIns:     &userID, // Set pointer to userID
	}

	err = s.brandRepo.CreateBrand(brand)
	return brand, err
}

func (s *BrandService) UpdateBrand(id uint, name string, description *string, userID uint) (*model.Brand, error) {
	if id == 0 {
		return nil, errors.New("invalid brand ID")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if brand exists
	brand, err := s.brandRepo.GetBrandByID(id)
	if err != nil {
		return nil, errors.New("brand not found")
	}

	// Check if new name conflicts with existing brands
	if name != "" && name != brand.Name {
		exists, err := s.brandRepo.CheckBrandExists(name)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("brand name already in use")
		}
	}

	// Prepare update data with audit trail
	updateData := make(map[string]interface{})
	if name != "" {
		updateData["name"] = name
	}
	if description != nil {
		updateData["description"] = description
	}
	// Always set the user who updated
	updateData["user_updt"] = userID

	err = s.brandRepo.UpdateBrand(id, updateData)
	if err != nil {
		return nil, err
	}

	updatedBrand, err := s.brandRepo.GetBrandByID(id)
	if err != nil {
		return nil, err
	}
	return &updatedBrand, nil
}

func (s *BrandService) DeleteBrand(id uint, userID uint) error {
	if id == 0 {
		return errors.New("invalid brand ID")
	}

	if userID == 0 {
		return errors.New("user ID is required for audit trail")
	}

	// Check if brand exists
	_, err := s.brandRepo.GetBrandByID(id)
	if err != nil {
		return errors.New("brand not found")
	}

	return s.brandRepo.DeleteBrandWithAudit(id, userID)
}

// GetDeletedBrands returns all soft deleted brands
func (s *BrandService) GetDeletedBrands() ([]model.Brand, error) {
	return s.brandRepo.GetDeletedBrands()
}

// RestoreBrand restores a soft deleted brand
func (s *BrandService) RestoreBrand(id uint, userID uint) (*model.Brand, error) {
	if id == 0 {
		return nil, errors.New("invalid brand ID")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if brand exists in deleted records
	deletedBrands, err := s.brandRepo.GetDeletedBrands()
	if err != nil {
		return nil, err
	}

	var foundBrand *model.Brand
	for _, brand := range deletedBrands {
		if brand.ID == id {
			foundBrand = &brand
			break
		}
	}

	if foundBrand == nil {
		return nil, errors.New("deleted brand not found")
	}

	err = s.brandRepo.RestoreBrand(id, userID)
	if err != nil {
		return nil, err
	}

	// Get restored brand
	restoredBrand, err := s.brandRepo.GetBrandByID(id)
	if err != nil {
		return nil, err
	}
	return &restoredBrand, nil
}
