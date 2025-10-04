package repository

import (
	"myapp/database"
	"myapp/internal/model"
)

type FileRepository struct{}

// fileWithUserResponse struct untuk response dengan user details
type fileWithUserResponse struct {
	ID           uint    `json:"id"`
	ModelType    string  `json:"modelType"`
	ModelID      uint    `json:"modelId"`
	Ext          string  `json:"ext"`
	FileURL      *string `json:"fileUrl"`
	UserInsName  *string `json:"userInsName"`
	UserUpdtName *string `json:"userUpdtName"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (r *FileRepository) GetAllFiles() ([]fileWithUserResponse, error) {
	var files []fileWithUserResponse

	result := database.DB.Table("files f").
		Select("f.id, f.model_type, f.model_id, f.ext, f.file_url, ui.name as user_ins_name, uu.name as user_updt_name, f.created_at, f.updated_at").
		Joins("LEFT JOIN users ui ON f.user_ins = ui.id AND ui.deleted_at IS NULL").
		Joins("LEFT JOIN users uu ON f.user_updt = uu.id AND uu.deleted_at IS NULL").
		Where("f.deleted_at IS NULL").
		Order("f.created_at DESC").
		Find(&files)

	return files, result.Error
}

func (r *FileRepository) GetFilesByModel(modelType string, modelID uint) ([]fileWithUserResponse, error) {
	var files []fileWithUserResponse

	result := database.DB.Table("files f").
		Select("f.id, f.model_type, f.model_id, f.ext, f.file_url, ui.name as user_ins_name, uu.name as user_updt_name, f.created_at, f.updated_at").
		Joins("LEFT JOIN users ui ON f.user_ins = ui.id AND ui.deleted_at IS NULL").
		Joins("LEFT JOIN users uu ON f.user_updt = uu.id AND uu.deleted_at IS NULL").
		Where("f.deleted_at IS NULL AND f.model_type = ? AND f.model_id = ?", modelType, modelID).
		Order("f.created_at DESC").
		Find(&files)

	return files, result.Error
}

func (r *FileRepository) GetFileByID(id uint) (fileWithUserResponse, error) {
	var file fileWithUserResponse

	result := database.DB.Table("files f").
		Select("f.id, f.model_type, f.model_id, f.ext, f.file_url, ui.name as user_ins_name, uu.name as user_updt_name, f.created_at, f.updated_at").
		Joins("LEFT JOIN users ui ON f.user_ins = ui.id AND ui.deleted_at IS NULL").
		Joins("LEFT JOIN users uu ON f.user_updt = uu.id AND uu.deleted_at IS NULL").
		Where("f.deleted_at IS NULL AND f.id = ?", id).
		First(&file)

	return file, result.Error
}

// GetFileModelByID returns model.File for service operations
func (r *FileRepository) GetFileModelByID(id uint) (model.File, error) {
	var file model.File
	result := database.DB.Where("id = ?", id).First(&file)
	return file, result.Error
}

func (r *FileRepository) CreateFile(file *model.File) error {
	return database.DB.Create(file).Error
}

func (r *FileRepository) UpdateFile(id uint, updateData map[string]interface{}) error {
	return database.DB.Model(&model.File{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *FileRepository) DeleteFileWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the file
	updateData := map[string]interface{}{
		"user_updt": userID,
	}

	// Update the audit field first
	err := database.DB.Model(&model.File{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	return database.DB.Delete(&model.File{}, id).Error
}

// GetDeletedFiles returns all soft deleted files
func (r *FileRepository) GetDeletedFiles() ([]fileWithUserResponse, error) {
	var files []fileWithUserResponse

	result := database.DB.Table("files f").
		Select("f.id, f.model_type, f.model_id, f.ext, f.file_url, ui.name as user_ins_name, uu.name as user_updt_name, f.created_at, f.updated_at").
		Joins("LEFT JOIN users ui ON f.user_ins = ui.id").
		Joins("LEFT JOIN users uu ON f.user_updt = uu.id").
		Where("f.deleted_at IS NOT NULL").
		Order("f.deleted_at DESC").
		Find(&files)

	return files, result.Error
}

// RestoreFile restores a soft deleted file
func (r *FileRepository) RestoreFile(id uint, userID uint) error {
	updateData := map[string]interface{}{
		"user_updt":  userID,
		"deleted_at": nil,
	}
	return database.DB.Unscoped().Model(&model.File{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *FileRepository) CheckModelExists(modelType string, modelID uint) (bool, error) {
	var count int64

	// Check based on model type
	switch modelType {
	case "product":
		result := database.DB.Model(&model.Product{}).Where("id = ?", modelID).Count(&count)
		return count > 0, result.Error
	case "user":
		result := database.DB.Model(&model.User{}).Where("id = ?", modelID).Count(&count)
		return count > 0, result.Error
	case "category":
		result := database.DB.Model(&model.Category{}).Where("id = ?", modelID).Count(&count)
		return count > 0, result.Error
	case "brand":
		result := database.DB.Model(&model.Brand{}).Where("id = ?", modelID).Count(&count)
		return count > 0, result.Error
	case "location":
		result := database.DB.Model(&model.Location{}).Where("id = ?", modelID).Count(&count)
		return count > 0, result.Error
	default:
		return false, nil
	}
}
