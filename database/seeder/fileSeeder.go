package seeder

import (
	"log"
	"myapp/internal/model"

	"gorm.io/gorm"
)

type FileSeeder struct{}

func NewFileSeeder() SeederInterface {
	return &FileSeeder{}
}

func (s *FileSeeder) GetName() string {
	return "FileSeeder"
}

func (s *FileSeeder) Seed(db *gorm.DB) error {
	log.Printf("🌱 Running %s...", s.GetName())

	// Check if files already exist
	var count int64
	db.Model(&model.File{}).Count(&count)

	if count > 0 {
		log.Printf("✅ %s: Files already exist, skipping...", s.GetName())
		return nil
	}

	// Create sample files
	files := []model.File{
		{
			ModelType: "product",
			ModelID:   1,
			Ext:       "jpg",
			FileURL:   stringPtr("https://mybucket.s3.amazonaws.com/2024/10/04/uuid1_product_image.jpg"),
			UserIns:   uintPtr(1),
		},
		{
			ModelType: "product",
			ModelID:   1,
			Ext:       "pdf",
			FileURL:   stringPtr("https://mybucket.s3.amazonaws.com/2024/10/04/uuid2_product_manual.pdf"),
			UserIns:   uintPtr(1),
		},
		{
			ModelType: "product",
			ModelID:   2,
			Ext:       "png",
			FileURL:   stringPtr("https://mybucket.s3.amazonaws.com/2024/10/04/uuid3_product_image.png"),
			UserIns:   uintPtr(1),
		},
		{
			ModelType: "user",
			ModelID:   1,
			Ext:       "jpg",
			FileURL:   stringPtr("https://mybucket.s3.amazonaws.com/2024/10/04/uuid4_profile_picture.jpg"),
			UserIns:   uintPtr(1),
		},
		{
			ModelType: "user",
			ModelID:   2,
			Ext:       "png",
			FileURL:   stringPtr("https://mybucket.s3.amazonaws.com/2024/10/04/uuid5_avatar.png"),
			UserIns:   uintPtr(2),
		},
		{
			ModelType: "category",
			ModelID:   1,
			Ext:       "jpg",
			FileURL:   stringPtr("https://mybucket.s3.amazonaws.com/2024/10/04/uuid6_category_banner.jpg"),
			UserIns:   uintPtr(1),
		},
		{
			ModelType: "brand",
			ModelID:   1,
			Ext:       "png",
			FileURL:   stringPtr("https://mybucket.s3.amazonaws.com/2024/10/04/uuid7_brand_logo.png"),
			UserIns:   uintPtr(1),
		},
		{
			ModelType: "location",
			ModelID:   1,
			Ext:       "jpg",
			FileURL:   stringPtr("https://mybucket.s3.amazonaws.com/2024/10/04/uuid8_location_photo.jpg"),
			UserIns:   uintPtr(1),
		},
		{
			ModelType: "product",
			ModelID:   3,
			Ext:       "gif",
			FileURL:   stringPtr("https://mybucket.s3.amazonaws.com/2024/10/04/uuid9_product_demo.gif"),
			UserIns:   uintPtr(2),
		},
		{
			ModelType: "product",
			ModelID:   1,
			Ext:       "docx",
			FileURL:   stringPtr("https://mybucket.s3.amazonaws.com/2024/10/04/uuid10_product_spec.docx"),
			UserIns:   uintPtr(1),
		},
	}

	result := db.Create(&files)
	if result.Error != nil {
		log.Printf("❌ %s: Error creating files: %v", s.GetName(), result.Error)
		return result.Error
	}

	log.Printf("✅ %s: Successfully created %d files", s.GetName(), len(files))
	return nil
}
