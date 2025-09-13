package database

import (
	"myapp/database/seeder"
)

func Seed() error {
	// Auto-run all seeders from seeder folder
	return seeder.RunAllSeeders(DB)
}