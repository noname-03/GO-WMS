package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
	"strings"
)

type LocationService struct {
	locationRepo *repository.LocationRepository
}

func NewLocationService() *LocationService {
	return &LocationService{
		locationRepo: repository.NewLocationRepository(),
	}
}

func (s *LocationService) GetAllLocations() (interface{}, error) {
	locations, err := s.locationRepo.GetAllLocations()
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (s *LocationService) GetLocationsByUser(userID uint) (interface{}, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	locations, err := s.locationRepo.GetLocationsByUser(userID)
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (s *LocationService) GetLocationsByType(locationType string) (interface{}, error) {
	if locationType == "" {
		return nil, errors.New("location type is required")
	}

	// Validate location type
	validTypes := []string{"gudang", "reseller"}
	isValid := false
	for _, validType := range validTypes {
		if strings.ToLower(locationType) == validType {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, errors.New("invalid location type. Must be 'gudang' or 'reseller'")
	}

	locations, err := s.locationRepo.GetLocationsByType(strings.ToLower(locationType))
	if err != nil {
		return nil, err
	}
	return locations, nil
}

func (s *LocationService) GetLocationByID(id uint) (interface{}, error) {
	location, err := s.locationRepo.GetLocationByID(id)
	if err != nil {
		return nil, err
	}
	return location, nil
}

func (s *LocationService) CreateLocation(userID uint, name string, address *string, phoneNumber *string, locationType string, createdByUserID uint) (interface{}, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}

	if name == "" {
		return nil, errors.New("location name is required")
	}

	if locationType == "" {
		return nil, errors.New("location type is required")
	}

	if createdByUserID == 0 {
		return nil, errors.New("created by user ID is required for audit trail")
	}

	// Validate location type
	validTypes := []string{"gudang", "reseller"}
	isValid := false
	for _, validType := range validTypes {
		if strings.ToLower(locationType) == validType {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, errors.New("invalid location type. Must be 'gudang' or 'reseller'")
	}

	// Check if user exists
	userExists, err := s.locationRepo.CheckUserExists(userID)
	if err != nil {
		return nil, err
	}
	if !userExists {
		return nil, errors.New("user not found")
	}

	// Check if location name already exists for this user
	nameExists, err := s.locationRepo.CheckLocationNameExists(userID, name, 0)
	if err != nil {
		return nil, err
	}
	if nameExists {
		return nil, errors.New("location name already exists for this user")
	}

	location := &model.Location{
		UserID:      userID,
		Name:        name,
		Address:     address,
		PhoneNumber: phoneNumber,
		Type:        strings.ToLower(locationType),
		UserIns:     &createdByUserID,
	}

	err = s.locationRepo.CreateLocation(location)
	if err != nil {
		return nil, err
	}

	// Fetch the created location with user details
	createdLocation, err := s.locationRepo.GetLocationByID(location.ID)
	if err != nil {
		return nil, err
	}

	return createdLocation, nil
}

func (s *LocationService) UpdateLocation(id uint, userID uint, name *string, address *string, phoneNumber *string, locationType *string, updatedByUserID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid location ID")
	}

	if updatedByUserID == 0 {
		return nil, errors.New("updated by user ID is required for audit trail")
	}

	// Check if location exists
	oldLocation, err := s.locationRepo.GetLocationModelByID(id)
	if err != nil {
		return nil, errors.New("location not found")
	}

	// If user ID is being changed, check if new user exists
	if userID != 0 && userID != oldLocation.UserID {
		userExists, err := s.locationRepo.CheckUserExists(userID)
		if err != nil {
			return nil, err
		}
		if !userExists {
			return nil, errors.New("user not found")
		}
	}

	// Use existing user if not provided
	if userID == 0 {
		userID = oldLocation.UserID
	}

	// Check if new name conflicts with existing locations for this user
	if name != nil && *name != "" && *name != oldLocation.Name {
		nameExists, err := s.locationRepo.CheckLocationNameExists(userID, *name, id)
		if err != nil {
			return nil, err
		}
		if nameExists {
			return nil, errors.New("location name already exists for this user")
		}
	}

	// Validate location type if provided
	if locationType != nil && *locationType != "" {
		validTypes := []string{"gudang", "reseller"}
		isValid := false
		for _, validType := range validTypes {
			if strings.ToLower(*locationType) == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			return nil, errors.New("invalid location type. Must be 'gudang' or 'reseller'")
		}
	}

	// Prepare update data with audit trail
	updateData := make(map[string]interface{})
	if userID != oldLocation.UserID {
		updateData["user_id"] = userID
	}
	if name != nil && *name != "" {
		updateData["name"] = *name
	}
	if address != nil {
		updateData["address"] = address
	}
	if phoneNumber != nil {
		updateData["phone_number"] = phoneNumber
	}
	if locationType != nil && *locationType != "" {
		updateData["type"] = strings.ToLower(*locationType)
	}
	// Always set the user who updated
	updateData["user_updt"] = updatedByUserID

	err = s.locationRepo.UpdateLocation(id, updateData)
	if err != nil {
		return nil, err
	}

	updatedLocation, err := s.locationRepo.GetLocationByID(id)
	if err != nil {
		return nil, err
	}

	return updatedLocation, nil
}

func (s *LocationService) DeleteLocation(id uint, userID uint) error {
	if id == 0 {
		return errors.New("invalid location ID")
	}

	if userID == 0 {
		return errors.New("user ID is required for audit trail")
	}

	// Check if location exists
	_, err := s.locationRepo.GetLocationModelByID(id)
	if err != nil {
		return errors.New("location not found")
	}

	return s.locationRepo.DeleteLocationWithAudit(id, userID)
}
