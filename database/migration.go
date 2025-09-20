package database

import (
	"log"
	"myapp/internal/model"
)

func Migrate() error {
	log.Println("Starting database migration...")

	err := DB.AutoMigrate(&model.User{}, &model.Role{}, &model.Brand{}, &model.Category{})
	if err != nil {
		log.Println("Migration failed:", err)
		return err
	}

	log.Println("Migration completed successfully!")
	return nil
}
