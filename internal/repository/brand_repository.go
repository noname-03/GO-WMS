package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"myapp/database"
	"myapp/internal/model"
	"myapp/pkg/redis"
)

type BrandRepository struct{}

func NewBrandRepository() *BrandRepository {
	return &BrandRepository{}
}

func (r *BrandRepository) GetAllBrands() ([]model.Brand, error) {
	cacheKey := "brands:all"

	// Try to get from cache first
	if cached, err := redis.Get(cacheKey); err == nil {
		var brands []model.Brand
		if err := json.Unmarshal([]byte(cached), &brands); err == nil {
			log.Printf("[REDIS] Cache hit for %s", cacheKey)
			log.Printf("[REDIS] Cached brands data: %+v", brands)
			return brands, nil
		}
	}

	// Cache miss, get from database
	var brands []model.Brand
	result := database.DB.Find(&brands)
	if result.Error != nil {
		return brands, result.Error
	}

	// Store in cache
	if data, err := json.Marshal(brands); err == nil {
		if err := redis.Set(cacheKey, string(data)); err != nil {
			log.Printf("[REDIS] Failed to cache %s: %v", cacheKey, err)
		} else {
			log.Printf("[REDIS] Cached %s", cacheKey)
		}
	}

	return brands, result.Error
}

func (r *BrandRepository) GetBrandByID(id uint) (model.Brand, error) {
	cacheKey := fmt.Sprintf("brand:id:%d", id)

	// Try to get from cache first
	if cached, err := redis.Get(cacheKey); err == nil {
		var brand model.Brand
		if err := json.Unmarshal([]byte(cached), &brand); err == nil {
			log.Printf("[REDIS] Cache hit for %s", cacheKey)
			return brand, nil
		}
	}

	// Cache miss, get from database
	var brand model.Brand
	result := database.DB.First(&brand, id)
	if result.Error != nil {
		return brand, result.Error
	}

	// Store in cache
	if data, err := json.Marshal(brand); err == nil {
		if err := redis.Set(cacheKey, string(data)); err != nil {
			log.Printf("[REDIS] Failed to cache %s: %v", cacheKey, err)
		} else {
			log.Printf("[REDIS] Cached %s", cacheKey)
		}
	}

	return brand, result.Error
}

func (r *BrandRepository) CreateBrand(brand *model.Brand) error {
	err := database.DB.Create(brand).Error
	if err != nil {
		return err
	}

	// Write-Through: Update cache instead of invalidate
	r.updateBrandCache()
	return nil
}

func (r *BrandRepository) UpdateBrand(id uint, updateData map[string]interface{}) error {
	err := database.DB.Model(&model.Brand{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Update both all brands cache and specific brand cache
	r.updateBrandCache()
	r.updateSpecificBrandCache(id)
	return nil
}

func (r *BrandRepository) DeleteBrandWithAudit(id uint, userID uint) error {
	// First update the user_updt field to track who deleted the brand
	updateData := map[string]interface{}{
		"user_updt": userID,
	}

	// Update the audit field first
	err := database.DB.Model(&model.Brand{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		return err
	}

	// Then perform the soft delete
	err = database.DB.Delete(&model.Brand{}, id).Error
	if err != nil {
		return err
	}

	// Update cache after delete
	r.updateBrandCache()
	r.invalidateSpecificBrandCache(id) // This one we can invalidate since it's deleted
	return nil
}

// updateBrandCache refreshes the all brands cache
func (r *BrandRepository) updateBrandCache() {
	var brands []model.Brand
	result := database.DB.Find(&brands)
	if result.Error != nil {
		log.Printf("[REDIS] Failed to fetch brands for cache update: %v", result.Error)
		return
	}

	if data, err := json.Marshal(brands); err == nil {
		if err := redis.Set("brands:all", string(data)); err != nil {
			log.Printf("[REDIS] Failed to update brands:all cache: %v", err)
		} else {
			log.Printf("[REDIS] Updated brands:all cache with %d brands", len(brands))
		}
	}
}

// updateSpecificBrandCache refreshes specific brand cache
func (r *BrandRepository) updateSpecificBrandCache(brandID uint) {
	var brand model.Brand
	result := database.DB.First(&brand, brandID)
	if result.Error != nil {
		log.Printf("[REDIS] Failed to fetch brand %d for cache update: %v", brandID, result.Error)
		return
	}

	cacheKey := fmt.Sprintf("brand:id:%d", brandID)
	if data, err := json.Marshal(brand); err == nil {
		if err := redis.Set(cacheKey, string(data)); err != nil {
			log.Printf("[REDIS] Failed to update %s cache: %v", cacheKey, err)
		} else {
			log.Printf("[REDIS] Updated %s cache", cacheKey)
		}
	}
}

// invalidateSpecificBrandCache removes specific brand cache
func (r *BrandRepository) invalidateSpecificBrandCache(brandID uint) {
	cacheKey := fmt.Sprintf("brand:id:%d", brandID)
	if err := redis.Delete(cacheKey); err != nil {
		log.Printf("[REDIS] Failed to invalidate %s cache: %v", cacheKey, err)
	} else {
		log.Printf("[REDIS] Invalidated %s cache", cacheKey)
	}
}

func (r *BrandRepository) CheckBrandExists(name string) (bool, error) {
	var count int64
	query := database.DB.Model(&model.Brand{}).Unscoped().Where("name ILIKE ?", name)

	result := query.Count(&count)
	return count > 0, result.Error
}
