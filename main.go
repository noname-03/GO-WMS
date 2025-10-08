package main

import (
	"log"
	"myapp/database"
	"myapp/internal/routes"
	"myapp/pkg/helper/s3"
	"myapp/pkg/redis"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// 1. Koneksi DB
	if err := database.ConnectDB(); err != nil {
		log.Fatal("DB connection error: ", err)
	}

	// 2. Koneksi Redis
	if err := redis.InitRedis(); err != nil {
		log.Printf("Redis connection warning: %v", err)
		log.Println("Application will continue without Redis caching")
	}

	// 3. Initialize S3 Client
	s3Client := s3.NewS3Client()
	if s3Client != nil {
		log.Println("✅ S3 client initialized successfully")
		// Test connection (optional)
		if err := s3Client.TestConnection(); err != nil {
			log.Printf("⚠️  S3 connection test failed: %v", err)
		}
	} else {
		log.Println("⚠️  S3 client initialization failed - check your S3 configuration")
	}

	// 4. Migration
	if err := database.Migrate(); err != nil {
		log.Fatal("Migration error: ", err)
	}

	// 5. Seeder
	if err := database.Seed(); err != nil {
		log.Fatal("Seed error: ", err)
	}

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	// Setup all routes
	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
