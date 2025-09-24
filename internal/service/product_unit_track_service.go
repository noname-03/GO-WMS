package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
	"myapp/internal/utils"
)

type ProductUnitTrackService struct {
	productUnitTrackRepo *repository.ProductUnitTrackRepository
	trackingUnitUtils    *utils.ProductUnitTrackingUtils
}

func NewProductUnitTrackService() *ProductUnitTrackService {
	return &ProductUnitTrackService{
		productUnitTrackRepo: repository.NewProductUnitTrackRepository(),
		trackingUnitUtils:    utils.NewProductUnitTrackingUtils(),
	}
}

// Business logic methods
func (s *ProductUnitTrackService) GetAllProductUnitTracks() (interface{}, error) {
	return s.productUnitTrackRepo.GetAllProductUnitTracks()
}

func (s *ProductUnitTrackService) GetProductUnitTracksByProductUnit(productUnitID uint) (interface{}, error) {
	if productUnitID == 0 {
		return nil, errors.New("invalid product unit ID")
	}

	// Check if product unit exists
	productUnitExists, err := s.productUnitTrackRepo.CheckProductUnitExists(productUnitID)
	if err != nil {
		return nil, err
	}
	if !productUnitExists {
		return nil, errors.New("product unit not found")
	}

	return s.productUnitTrackRepo.GetProductUnitTracksByProductUnit(productUnitID)
}

func (s *ProductUnitTrackService) GetProductUnitTracksByProduct(productID uint) (interface{}, error) {
	if productID == 0 {
		return nil, errors.New("invalid product ID")
	}

	return s.productUnitTrackRepo.GetProductUnitTracksByProduct(productID)
}

func (s *ProductUnitTrackService) GetProductUnitTrackByID(id uint) (interface{}, error) {
	track, err := s.productUnitTrackRepo.GetProductUnitTrackByID(id)
	if err != nil {
		return nil, err
	}
	return track, nil
}

func (s *ProductUnitTrackService) CreateProductUnitTrack(
	productUnitID uint,
	description string,
	userID uint) (interface{}, error) {

	if productUnitID == 0 {
		return nil, errors.New("product unit ID is required")
	}

	if description == "" {
		return nil, errors.New("transaction type is required")
	}

	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if product unit exists
	productUnitExists, err := s.productUnitTrackRepo.CheckProductUnitExists(productUnitID)
	if err != nil {
		return nil, err
	}
	if !productUnitExists {
		return nil, errors.New("product unit not found")
	}

	productUnitTrack := &model.ProductUnitTrack{
		ProductUnitID: productUnitID,
		Description:   description,
		UserIns:       &userID,
	}

	err = s.productUnitTrackRepo.CreateProductUnitTrack(productUnitTrack)
	if err != nil {
		return nil, err
	}

	// Fetch the created product unit track with details
	createdProductUnitTrack, err := s.productUnitTrackRepo.GetProductUnitTrackByID(productUnitTrack.ID)
	if err != nil {
		return nil, err
	}

	return createdProductUnitTrack, nil
}

func (s *ProductUnitTrackService) DeleteProductUnitTrack(id uint, userID uint) error {
	if id == 0 {
		return errors.New("invalid product unit track ID")
	}

	if userID == 0 {
		return errors.New("user ID is required for audit trail")
	}

	// Check if product unit track exists
	_, err := s.productUnitTrackRepo.GetProductUnitTrackModelByID(id)
	if err != nil {
		return errors.New("product unit track not found")
	}

	return s.productUnitTrackRepo.DeleteProductUnitTrackWithAudit(id, userID)
}

// TrackCreate creates a tracking record for product unit creation
func (s *ProductUnitTrackService) TrackCreate(productUnit model.ProductUnit, userID uint) error {
	description := s.trackingUnitUtils.GenerateCreateDescriptionProductUnit(productUnit)
	_, err := s.CreateProductUnitTrack(productUnit.ID, description, userID)
	return err
}

// TrackUpdateFromChanges creates a tracking record for product unit updates using update data
func (s *ProductUnitTrackService) TrackUpdateFromChangesProductUnit(updateData map[string]interface{}, oldBatch model.ProductUnit, userID uint) error {
	description := s.trackingUnitUtils.GenerateUpdateDescriptionFromChangesProductUnit(updateData, oldBatch)
	_, err := s.CreateProductUnitTrack(oldBatch.ID, description, userID)
	return err
}

// TrackDelete creates a tracking record for product unit deletion
func (s *ProductUnitTrackService) TrackDeleteProductUnit(productUnit model.ProductUnit, userID uint) error {
	description := s.trackingUnitUtils.GenerateDeleteDescriptionProductUnit(productUnit)
	_, err := s.CreateProductUnitTrack(productUnit.ID, description, userID)
	return err
}
