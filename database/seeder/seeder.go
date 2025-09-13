package seeder

import "gorm.io/gorm"

// SeederInterface defines the contract for all seeders
type SeederInterface interface {
	Seed(db *gorm.DB) error
	GetName() string
}

// SeederRegistry holds all registered seeders
type SeederRegistry struct {
	seeders []SeederInterface
}

// NewSeederRegistry creates a new seeder registry
func NewSeederRegistry() *SeederRegistry {
	return &SeederRegistry{
		seeders: make([]SeederInterface, 0),
	}
}

// Register adds a seeder to the registry
func (r *SeederRegistry) Register(seeder SeederInterface) {
	r.seeders = append(r.seeders, seeder)
}

// RunAll executes all registered seeders
func (r *SeederRegistry) RunAll(db *gorm.DB) error {
	for _, seeder := range r.seeders {
		if err := seeder.Seed(db); err != nil {
			return err
		}
	}
	return nil
}

// GetSeeders returns all registered seeders
func (r *SeederRegistry) GetSeeders() []SeederInterface {
	return r.seeders
}