package seeder

import (
	"log"
	"myapp/internal/model"
	"gorm.io/gorm"
)

type RoleSeeder struct{}

func NewRoleSeeder() SeederInterface {
	return &RoleSeeder{}
}

func (s *RoleSeeder) GetName() string {
	return "RoleSeeder"
}

func (s *RoleSeeder) Seed(db *gorm.DB) error {
	log.Println("🌱 Running RoleSeeder...")
	
	roles := []model.Role{
		{Name: "Admin", Description: "Administrator with full access"},
		{Name: "User", Description: "Regular user with limited access"},
		{Name: "Manager", Description: "Manager with elevated privileges"},
	}

	for _, role := range roles {
		var existing model.Role
		result := db.Where("name = ?", role.Name).First(&existing)
		if result.Error != nil {
			// Role doesn't exist, create it
			if err := db.Create(&role).Error; err != nil {
				log.Printf("❌ Failed to seed role %s: %v", role.Name, err)
				return err
			}
			log.Printf("✅ Role '%s' created successfully", role.Name)
		} else {
			log.Printf("✅ RoleSeeder: Role '%s' already exists, skipping...", role.Name)
		}
	}

	return nil
}