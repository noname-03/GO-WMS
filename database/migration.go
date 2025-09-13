package database

import (
    "myapp/internal/model"
    "log"
)

func Migrate() error {
    log.Println("Starting database migration...")
    
    err := DB.AutoMigrate(&model.User{}, &model.Role{})
    if err != nil {
        log.Println("Migration failed:", err)
        return err
    }
    
    log.Println("Migration completed successfully!")
    return nil
}