package service

import (
	"errors"
	"myapp/internal/model"
	"myapp/internal/repository"
	"myapp/pkg/helper/s3"
	"strings"
)

type FileService struct {
	fileRepo *repository.FileRepository
	s3Client *s3.S3Client
}

func NewFileService() *FileService {
	s3Client := s3.NewS3Client()
	if s3Client == nil {
		// Log warning but don't fail - allow service to work without S3
		// In production, you might want to fail here
		// log.Printf("Warning: S3 client initialization failed")
	}

	return &FileService{
		fileRepo: repository.NewFileRepository(),
		s3Client: s3Client,
	}
}

// Business logic methods
func (s *FileService) GetAllFiles() (interface{}, error) {
	return s.fileRepo.GetAllFiles()
}

func (s *FileService) GetFilesByModel(modelType string, modelID uint) (interface{}, error) {
	if modelType == "" {
		return nil, errors.New("model type is required")
	}
	if modelID == 0 {
		return nil, errors.New("model ID is required")
	}

	// Validate model type
	validModelTypes := []string{"product", "user", "category", "brand", "location"}
	isValid := false
	for _, validType := range validModelTypes {
		if strings.ToLower(modelType) == validType {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, errors.New("invalid model type. allowed: product, user, category, brand, location")
	}

	// Check if model exists
	exists, err := s.fileRepo.CheckModelExists(strings.ToLower(modelType), modelID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("model not found")
	}

	return s.fileRepo.GetFilesByModel(strings.ToLower(modelType), modelID)
}

func (s *FileService) GetFileByID(id uint) (interface{}, error) {
	file, err := s.fileRepo.GetFileByID(id)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (s *FileService) UploadFile(modelType string, modelID uint, ext string, fileData []byte, fileName string, userID uint) (interface{}, error) {
	// Check if S3Client is initialized
	if s.s3Client == nil {
		return nil, errors.New("S3 client not initialized")
	}

	if modelType == "" {
		return nil, errors.New("model type is required")
	}
	if modelID == 0 {
		return nil, errors.New("model ID is required")
	}
	if ext == "" {
		return nil, errors.New("file extension is required")
	}
	if len(fileData) == 0 {
		return nil, errors.New("file data is required")
	}
	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Validate model type
	validModelTypes := []string{"product", "user", "category", "brand", "location"}
	isValid := false
	for _, validType := range validModelTypes {
		if strings.ToLower(modelType) == validType {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, errors.New("invalid model type. allowed: product, user, category, brand, location")
	}

	// Check if model exists
	exists, err := s.fileRepo.CheckModelExists(strings.ToLower(modelType), modelID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("model not found")
	}

	// Validate file extension
	allowedExts := []string{"jpg", "jpeg", "png", "gif", "pdf", "doc", "docx", "xls", "xlsx"}
	isValidExt := false
	extLower := strings.ToLower(ext)
	for _, allowedExt := range allowedExts {
		if extLower == allowedExt {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		return nil, errors.New("invalid file extension. allowed: jpg, jpeg, png, gif, pdf, doc, docx, xls, xlsx")
	}

	// Upload to S3 with dynamic path
	fileURL, err := s.s3Client.UploadFile(fileData, strings.ToLower(modelType), modelID, extLower)
	if err != nil {
		return nil, errors.New("failed to upload file to S3: " + err.Error())
	}

	file := &model.File{
		ModelType: strings.ToLower(modelType),
		ModelID:   modelID,
		Ext:       extLower,
		FileURL:   &fileURL,
		UserIns:   &userID,
	}

	err = s.fileRepo.CreateFile(file)
	if err != nil {
		// If database save fails, try to delete the uploaded file from S3
		s.s3Client.DeleteFile(fileURL)
		return nil, err
	}

	// Fetch the created file with user details
	createdFile, err := s.fileRepo.GetFileByID(file.ID)
	if err != nil {
		return nil, err
	}

	return createdFile, nil
}

func (s *FileService) UpdateFile(id uint, modelType string, modelID uint, ext string, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid file ID")
	}
	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	// Check if file exists
	file, err := s.fileRepo.GetFileModelByID(id)
	if err != nil {
		return nil, errors.New("file not found")
	}

	// Prepare update data with audit trail
	updateData := make(map[string]interface{})

	if modelType != "" && strings.ToLower(modelType) != file.ModelType {
		// Validate model type
		validModelTypes := []string{"product", "user", "category", "brand", "location"}
		isValid := false
		for _, validType := range validModelTypes {
			if strings.ToLower(modelType) == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			return nil, errors.New("invalid model type. allowed: product, user, category, brand, location")
		}
		updateData["model_type"] = strings.ToLower(modelType)
	}

	if modelID != 0 && modelID != file.ModelID {
		updateData["model_id"] = modelID
	}

	if ext != "" && strings.ToLower(ext) != file.Ext {
		// Validate file extension
		allowedExts := []string{"jpg", "jpeg", "png", "gif", "pdf", "doc", "docx", "xls", "xlsx"}
		isValidExt := false
		extLower := strings.ToLower(ext)
		for _, allowedExt := range allowedExts {
			if extLower == allowedExt {
				isValidExt = true
				break
			}
		}
		if !isValidExt {
			return nil, errors.New("invalid file extension. allowed: jpg, jpeg, png, gif, pdf, doc, docx, xls, xlsx")
		}
		updateData["ext"] = extLower
	}

	// Always set the user who updated
	updateData["user_updt"] = userID

	err = s.fileRepo.UpdateFile(id, updateData)
	if err != nil {
		return nil, err
	}

	updatedFile, err := s.fileRepo.GetFileByID(id)
	if err != nil {
		return nil, err
	}
	return updatedFile, nil
}

func (s *FileService) DeleteFile(id uint, userID uint) error {
	if id == 0 {
		return errors.New("invalid file ID")
	}
	if userID == 0 {
		return errors.New("user ID is required for audit trail")
	}

	// Check if file exists
	file, err := s.fileRepo.GetFileModelByID(id)
	if err != nil {
		return errors.New("file not found")
	}

	// Delete from S3 if file URL exists
	if file.FileURL != nil && *file.FileURL != "" {
		err = s.s3Client.DeleteFile(*file.FileURL)
		if err != nil {
			// Log error but don't fail the deletion
			// Consider using a logger here
		}
	}

	return s.fileRepo.DeleteFileWithAudit(id, userID)
}

// GetDeletedFiles returns all soft deleted files
func (s *FileService) GetDeletedFiles() (interface{}, error) {
	return s.fileRepo.GetDeletedFiles()
}

// RestoreFile restores a soft deleted file
func (s *FileService) RestoreFile(id uint, userID uint) (interface{}, error) {
	if id == 0 {
		return nil, errors.New("invalid file ID")
	}
	if userID == 0 {
		return nil, errors.New("user ID is required for audit trail")
	}

	err := s.fileRepo.RestoreFile(id, userID)
	if err != nil {
		return nil, err
	}

	restoredFile, err := s.fileRepo.GetFileByID(id)
	if err != nil {
		return nil, err
	}
	return restoredFile, nil
}
