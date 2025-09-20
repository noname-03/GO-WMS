package service

import (
	"fmt"
	"myapp/internal/model"
	"myapp/internal/repository"
	"myapp/internal/utils"
)

type ProductBatchTrackService struct {
	repository    *repository.ProductBatchTrackRepository
	trackingUtils *utils.ProductBatchTrackingUtils
}

func NewProductBatchTrackService() *ProductBatchTrackService {
	return &ProductBatchTrackService{
		repository:    repository.NewProductBatchTrackRepository(),
		trackingUtils: utils.NewProductBatchTrackingUtils(),
	}
}

// GetAllTracks retrieves all product batch tracking records
func (s *ProductBatchTrackService) GetAllTracks() ([]model.ProductBatchTrack, error) {
	return s.repository.GetAllTracks()
}

// GetTracksByProductBatchID retrieves all tracking records for a specific product batch
func (s *ProductBatchTrackService) GetTracksByProductBatchID(productBatchID uint) ([]model.ProductBatchTrack, error) {
	return s.repository.GetTracksByProductBatchID(productBatchID)
}

// GetTrackByID retrieves a specific tracking record by ID
func (s *ProductBatchTrackService) GetTrackByID(id uint) (model.ProductBatchTrack, error) {
	return s.repository.GetTrackByID(id)
}

// CreateTrackingRecord creates a new tracking record for product batch changes
func (s *ProductBatchTrackService) CreateTrackingRecord(productBatchID uint, description string, userID uint) (model.ProductBatchTrack, error) {
	track := model.ProductBatchTrack{
		ProductBatchID: productBatchID,
		Description:    description,
		UserInst:       userID,
	}

	err := s.repository.CreateTrack(&track)
	if err != nil {
		return track, fmt.Errorf("failed to create tracking record: %w", err)
	}

	// Fetch the created record with relationships
	return s.repository.GetTrackByID(track.ID)
}

// TrackCreate creates a tracking record for product batch creation
func (s *ProductBatchTrackService) TrackCreate(productBatch model.ProductBatch, userID uint) error {
	description := s.trackingUtils.GenerateCreateDescription(productBatch)
	_, err := s.CreateTrackingRecord(productBatch.ID, description, userID)
	return err
}

// TrackUpdateFromChanges creates a tracking record for product batch updates using update data
func (s *ProductBatchTrackService) TrackUpdateFromChanges(updateData map[string]interface{}, oldBatch model.ProductBatch, userID uint) error {
	description := s.trackingUtils.GenerateUpdateDescriptionFromChanges(updateData, oldBatch)
	_, err := s.CreateTrackingRecord(oldBatch.ID, description, userID)
	return err
}

// TrackUpdate creates a tracking record for product batch updates (legacy method, kept for compatibility)
func (s *ProductBatchTrackService) TrackUpdate(oldBatch, newBatch model.ProductBatch, userID uint) error {
	// This method is less precise, use TrackUpdateFromChanges instead
	description := "Product batch updated"
	_, err := s.CreateTrackingRecord(newBatch.ID, description, userID)
	return err
}

// TrackDelete creates a tracking record for product batch deletion
func (s *ProductBatchTrackService) TrackDelete(productBatch model.ProductBatch, userID uint) error {
	description := s.trackingUtils.GenerateDeleteDescription(productBatch)
	_, err := s.CreateTrackingRecord(productBatch.ID, description, userID)
	return err
}

// TrackCustomAction creates a tracking record for custom actions
func (s *ProductBatchTrackService) TrackCustomAction(productBatchID uint, customDescription string, userID uint) error {
	_, err := s.CreateTrackingRecord(productBatchID, customDescription, userID)
	return err
}

// GetTracksByUserID retrieves all tracking records made by a specific user
func (s *ProductBatchTrackService) GetTracksByUserID(userID uint) ([]model.ProductBatchTrack, error) {
	return s.repository.GetTracksByUserID(userID)
}

// GetLatestTrackForProductBatch retrieves the most recent tracking record for a product batch
func (s *ProductBatchTrackService) GetLatestTrackForProductBatch(productBatchID uint) (model.ProductBatchTrack, error) {
	return s.repository.GetLatestTrackForProductBatch(productBatchID)
}

// GetProductBatchHistory retrieves the complete history of changes for a product batch
func (s *ProductBatchTrackService) GetProductBatchHistory(productBatchID uint) ([]model.ProductBatchTrack, error) {
	tracks, err := s.repository.GetTracksByProductBatchID(productBatchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product batch history: %w", err)
	}
	return tracks, nil
}

// CountTracksByProductBatchID counts tracking records for a specific product batch
func (s *ProductBatchTrackService) CountTracksByProductBatchID(productBatchID uint) (int64, error) {
	return s.repository.CountTracksByProductBatchID(productBatchID)
}

// DeleteTrack removes a tracking record (rarely used, for admin purposes)
func (s *ProductBatchTrackService) DeleteTrack(id uint) error {
	return s.repository.DeleteTrack(id)
}
