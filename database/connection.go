package database

import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "os"
    "log"
)

var DB *gorm.DB

func ConnectDB() error {
    dsn := os.Getenv("DB_DSN")
    log.Println("Connecting to database with DSN:", dsn)
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Println("Database connection failed:", err)
        return err
    }
    
    log.Println("Database connected successfully!")
    DB = db
    return nil
}