package seeder

import (
	"log"
	"myapp/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserSeeder struct{}

func NewUserSeeder() SeederInterface {
	return &UserSeeder{}
}

func (s *UserSeeder) GetName() string {
	return "UserSeeder"
}

func (s *UserSeeder) Seed(db *gorm.DB) error {
	log.Printf("🌱 Running %s...", s.GetName())
	
	// Check if users already exist
	var count int64
	db.Model(&model.User{}).Count(&count)
	
	if count > 0 {
		log.Printf("✅ %s: Users already exist, skipping...", s.GetName())
		return nil
	}

	// Hash default password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create sample users
	users := []model.User{
		{
			Name:     "Admin User",
			Email:    "admin@wms.com",
			Password: string(hashedPassword),
			Acusername: "system",
		},
		{
			Name:     "User",
			Email:    "user@wms.com",
			Password: string(hashedPassword),
			Acusername: "system",
		},
	}

	result := db.Create(&users)
	if result.Error != nil {
		log.Printf("❌ %s: Error creating users: %v", s.GetName(), result.Error)
		return result.Error
	}

	log.Printf("✅ %s: Successfully created %d users", s.GetName(), len(users))
	return nil
}