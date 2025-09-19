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
    result := database.DB.Where("email = ?", email).First(&user)
    log.Printf("GetUserByEmail query - Email: %s, Result: %+v, Error: %v, RowsAffected: %d", 
               email, user, result.Error, result.RowsAffected)
    
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            log.Printf("User with email %s not found", email)
            return nil, nil // Tidak ditemukan
        }
        log.Printf("Database error: %v", result.Error)
        return nil, result.Error
    }
    
    log.Printf("User found with ID: %d", user.ID)
    return &user, nil
}

func (r *UserRepository) CreateUser(user *model.User) error {
	result := database.DB.Create(user)
	return result.Error
}

// Raw SQL Queries
func (r *UserRepository) GetUsersWithRawSQL() ([]model.User, error) {
	var users []model.User
	
	query := "SELECT id, name, email, created_at, updated_at FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC"
	result := database.DB.Raw(query).Scan(&users)
	
	return users, result.Error
}

func (r *UserRepository) GetUserByIDRaw(id string) (*model.User, error) {
	var user model.User
	
	query := "SELECT * FROM users WHERE id = ? AND deleted_at IS NULL"
	result := database.DB.Raw(query, id).Scan(&user)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	if result.RowsAffected == 0 {
		return nil, nil
	}
	
	return &user, nil
}

func (r *UserRepository) GetUsersWithStats() ([]UserResult, error) {
	var users []UserResult

	query := "SELECT id, name, email, COUNT(*) OVER() as total FROM users WHERE deleted_at IS NULL"
	result := database.DB.Raw(query).Scan(&users)
	
	return users, result.Error
}

func (r *UserRepository) GetUsersCountByStatus() (map[string]int, error) {
	type StatusCount struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}
	
	var results []StatusCount
	
	query := `
		SELECT 
			CASE 
				WHEN deleted_at IS NULL THEN 'active'
				ELSE 'deleted'
			END as status,
			COUNT(*) as count
		FROM users 
		GROUP BY status
	`
	
	result := database.DB.Raw(query).Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}
	
	// Convert to map
	statusMap := make(map[string]int)
	for _, r := range results {
		statusMap[r.Status] = r.Count
	}
	
	return statusMap, nil
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

// Legacy function for backward compatibility
func FetchUsersFromDB() ([]string, error) {
    return []string{"Alice", "Bob"}, nil
}