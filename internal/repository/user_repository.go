package repository

import (
	"errors"
	"log"
	"myapp/database"
	"myapp/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Basic GORM queries
func (r *UserRepository) GetAllUsers() ([]model.User, error) {
	var users []model.User
	result := database.DB.Select("id, name, email").Find(&users)
	return users, result.Error
}

func (r *UserRepository) GetUsersMinimal() ([]UserMinimal, error) {
	var users []UserMinimal
	result := database.DB.Model(&model.User{}).Select("id, name").Find(&users)
	return users, result.Error
}

func (r *UserRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	query := "SELECT * FROM users WHERE email = ?"
	result := database.DB.Raw(query, email).Scan(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("User with email %s not found", email)
			return nil, nil // Tidak ditemukan
		}
		log.Printf("Database error: %v", result.Error)
		return nil, result.Error
	}

	// log.Printf("User found with ID: %d, Email: %s", user.ID, user.Email)
	return &user, nil
}

func (r *UserRepository) CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

// Raw SQL Queries
func (r *UserRepository) GetUsersWithRawSQL() ([]model.User, error) {
	var users []model.User

	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC"
	result := database.DB.Raw(query).Scan(&users)

	return users, result.Error
}

func (r *UserRepository) GetUsersWithStats() ([]UserResult, error) {
	var users []UserResult

	query := "SELECT id, name, email, COUNT(*) OVER() as total FROM users WHERE deleted_at IS NULL"
	result := database.DB.Raw(query).Scan(&users)

	return users, result.Error
}

func (r *UserRepository) SearchUsersRaw(keyword string, limit int, offset int) ([]model.User, error) {
	var users []model.User

	query := `
		SELECT id, name, email, created_at, updated_at 
		FROM users 
		WHERE (name ILIKE ? OR email ILIKE ?) 
		AND deleted_at IS NULL 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`

	searchTerm := "%" + keyword + "%"
	result := database.DB.Raw(query, searchTerm, searchTerm, limit, offset).Scan(&users)

	return users, result.Error
}

func (r *UserRepository) GetUsersStats() (*UserStats, error) {
	var stats UserStats

	query := `
		SELECT 
			COUNT(*) as total_users,
			COUNT(CASE WHEN deleted_at IS NULL THEN 1 END) as active_users,
			COUNT(CASE WHEN deleted_at IS NOT NULL THEN 1 END) as deleted_users
		FROM users
	`

	result := database.DB.Raw(query).Scan(&stats)
	if result.Error != nil {
		return nil, result.Error
	}

	return &stats, nil
}

// CheckEmailExists checks if email already exists (excluding a specific user ID)
func (r *UserRepository) CheckEmailExists(email string, excludeID uint) (bool, error) {
	var count int64
	query := database.DB.Model(&model.User{}).Where("email = ?", email)

	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	result := query.Count(&count)
	return count > 0, result.Error
}

// Struct definitions
type UserMinimal struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UserResult struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Total int    `json:"total"`
}

type UserStats struct {
	TotalUsers   int `json:"total_users"`
	ActiveUsers  int `json:"active_users"`
	DeletedUsers int `json:"deleted_users"`
}
