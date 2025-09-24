package seeder

import (
	"log"
	"myapp/internal/model"

	"gorm.io/gorm"
)

type LocationSeeder struct{}

func NewLocationSeeder() SeederInterface {
	return &LocationSeeder{}
}

func (s *LocationSeeder) GetName() string {
	return "LocationSeeder"
}

func (s *LocationSeeder) Seed(db *gorm.DB) error {
	log.Println("🌱 Running LocationSeeder...")

	// Check if data already exists
	var count int64
	db.Model(&model.Location{}).Count(&count)
	if count > 0 {
		log.Println("Location data already exists, skipping...")
		return nil
	}

	// Get some users for reference
	var users []model.User
	if err := db.Limit(5).Find(&users).Error; err != nil {
		return err
	}

	if len(users) == 0 {
		log.Println("No users found, skipping Location seeding...")
		return nil
	}

	// Get admin user for audit trail
	var adminUser model.User
	if err := db.Where("email = ?", "admin@example.com").First(&adminUser).Error; err != nil {
		log.Println("Admin user not found, using first user for audit trail")
		if len(users) > 0 {
			adminUser = users[0]
		}
	}

	var userID *uint
	if adminUser.ID != 0 {
		userID = &adminUser.ID
	}

	locations := []model.Location{
		{
			UserID:      users[0].ID,
			Name:        "Gudang Pusat Jakarta",
			Address:     stringPtrLocation("Jl. Raya Jakarta No. 123, Jakarta Pusat"),
			PhoneNumber: stringPtrLocation("021-1234567"),
			Type:        "gudang",
			UserIns:     userID,
		},
		{
			UserID:      users[0].ID,
			Name:        "Gudang Cabang Surabaya",
			Address:     stringPtrLocation("Jl. Ahmad Yani No. 456, Surabaya"),
			PhoneNumber: stringPtrLocation("031-9876543"),
			Type:        "gudang",
			UserIns:     userID,
		},
		{
			UserID:      users[1].ID,
			Name:        "Toko Sinar Jaya",
			Address:     stringPtrLocation("Jl. Pahlawan No. 789, Bandung"),
			PhoneNumber: stringPtrLocation("022-5555666"),
			Type:        "reseller",
			UserIns:     userID,
		},
		{
			UserID:      users[1].ID,
			Name:        "Minimarket Berkah",
			Address:     stringPtrLocation("Jl. Merdeka No. 321, Yogyakarta"),
			PhoneNumber: stringPtrLocation("0274-777888"),
			Type:        "reseller",
			UserIns:     userID,
		},
		{
			UserID:      users[2].ID,
			Name:        "Gudang Regional Medan",
			Address:     stringPtrLocation("Jl. Gatot Subroto No. 654, Medan"),
			PhoneNumber: stringPtrLocation("061-4444555"),
			Type:        "gudang",
			UserIns:     userID,
		},
	}

	// Additional locations if more users are available
	if len(users) > 3 {
		additionalLocations := []model.Location{
			{
				UserID:      users[3].ID,
				Name:        "Toko Makmur Sentosa",
				Address:     stringPtrLocation("Jl. Sudirman No. 147, Semarang"),
				PhoneNumber: stringPtrLocation("024-3333444"),
				Type:        "reseller",
				UserIns:     userID,
			},
			{
				UserID:      users[3].ID,
				Name:        "Gudang Distribusi Bali",
				Address:     stringPtrLocation("Jl. Bypass Ngurah Rai No. 258, Denpasar"),
				PhoneNumber: stringPtrLocation("0361-2222333"),
				Type:        "gudang",
				UserIns:     userID,
			},
		}
		locations = append(locations, additionalLocations...)
	}

	if len(users) > 4 {
		moreLocations := []model.Location{
			{
				UserID:      users[4].ID,
				Name:        "Supermarket Jaya Abadi",
				Address:     stringPtrLocation("Jl. Diponegoro No. 369, Malang"),
				PhoneNumber: stringPtrLocation("0341-1111222"),
				Type:        "reseller",
				UserIns:     userID,
			},
		}
		locations = append(locations, moreLocations...)
	}

	// Insert locations
	for _, location := range locations {
		if err := db.Create(&location).Error; err != nil {
			log.Printf("Failed to create location: %v", err)
			return err
		}
	}

	log.Printf("✅ Created %d locations", len(locations))
	return nil
}

// Helper function
func stringPtrLocation(s string) *string {
	return &s
}
