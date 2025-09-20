package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService() *ProductService {
	return &ProductService{
		productRepo: repository.NewProductRepository(),
	}
}

// Business logic methods
func (s *ProductService) GetAllProducts() ([]model.Product, error) {
	return s.productRepo.GetAllProducts()
}

func (s *ProductService) GetProductsByCategory(categoryID uint) ([]model.Product, error) {
	if categoryID == 0 {
		return nil, errors.New("invalid category ID")
	}

	// Check if category exists
	categoryExists, err := s.productRepo.CheckCategoryExists(categoryID)
	if err != nil {
		return nil, err
	}
	if !categoryExists {
		return nil, errors.New("category not found")
	}

	return s.productRepo.GetProductsByCategory(categoryID)
}

func (s *ProductService) GetProductByID(id uint) (*model.Product, error) {
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *ProductService) CreateProduct(categoryID uint, name string, description *string, userID uint) (*model.Product, error) {
	if categoryID == 0 {
		return nil, errors.New("category ID is required")
	}

	if name == "" {
		return nil, errors.New("product name is required")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if category exists
	categoryExists, err := s.productRepo.CheckCategoryExists(categoryID)
	if err != nil {
		return nil, err
	}
	if !categoryExists {
		return nil, errors.New("category not found")
	}

	// Check if product exists for this category
	exists, err := s.productRepo.CheckProductExists(name, categoryID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("product already exists for this category")
	}

	product := &model.Product{
		CategoryID:  categoryID,
		Name:        name,
		Description: description,
		UserIns:     &userID, // Set pointer to userID
	}

	err = s.productRepo.CreateProduct(product)
	if err != nil {
		return nil, err
	}

	// Fetch the created product with relationships
	createdProduct, err := s.productRepo.GetProductByID(product.ID)
	if err != nil {
		return nil, err
	}

	return &createdProduct, nil
}

func (s *ProductService) UpdateProduct(id uint, categoryID uint, name string, description *string, userID uint) (*model.Product, error) {
	if id == 0 {
		return nil, errors.New("invalid product ID")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if product exists
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// If category ID is being changed, check if new category exists
	if categoryID != 0 && categoryID != product.CategoryID {
		categoryExists, err := s.productRepo.CheckCategoryExists(categoryID)
		if err != nil {
			return nil, err
		}
		if !categoryExists {
			return nil, errors.New("category not found")
		}
	}

	// Use existing category if not provided
	if categoryID == 0 {
		categoryID = product.CategoryID
	}

	// Note: No duplicate check for update since we don't use excludeID
	// This means update will fail if trying to use existing name in same category

	// Prepare update data with audit trail
	updateData := make(map[string]interface{})
	if categoryID != product.CategoryID {
		updateData["category_id"] = categoryID
	}
	if name != "" {
		updateData["name"] = name
	}
	if description != nil {
		updateData["description"] = description
	}
	// Always set the user who updated
	updateData["user_updt"] = userID

	err = s.productRepo.UpdateProduct(id, updateData)
	if err != nil {
		return nil, err
	}

	updatedProduct, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	return &updatedProduct, nil
}

func (s *ProductService) DeleteProduct(id uint, userID uint) error {
	if id == 0 {
		return errors.New("invalid product ID")
	}

	if userID == 0 {
		return errors.New("user ID is required for audit trail")
	}

	// Check if product exists
	_, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	return s.productRepo.DeleteProductWithAudit(id, userID)
}
