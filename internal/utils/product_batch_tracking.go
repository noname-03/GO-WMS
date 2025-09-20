package utils

import (
	"fmt"
	"myapp/internal/model"
	"time"
)

// ProductBatchTrackingUtils provides utility functions for tracking ProductBatch changes
type ProductBatchTrackingUtils struct{}

// GenerateCreateDescription creates description for product batch creation
func (u *ProductBatchTrackingUtils) GenerateCreateDescription(productBatch model.ProductBatch) string {
	description := "Product batch created"
	if productBatch.CodeBatch != nil {
		description += " with code: " + *productBatch.CodeBatch
	}
	if productBatch.UnitPrice != nil {
		description += ", unit price: " + u.formatFloat(*productBatch.UnitPrice)
	}
	description += ", expiry date: " + productBatch.ExpDate.Format("2006-01-02")

	return description
}

// GenerateUpdateDescriptionFromChanges creates description based on actual changes made
func (u *ProductBatchTrackingUtils) GenerateUpdateDescriptionFromChanges(updateData map[string]interface{}, oldBatch model.ProductBatch) string {
	var changes []string

	// Check each field that was updated
	if newProductID, exists := updateData["product_id"]; exists {
		if newProductID != oldBatch.ProductID {
			changes = append(changes, fmt.Sprintf("changed product from ID %d to ID %d", oldBatch.ProductID, newProductID))
		}
	}

	if newCodeBatch, exists := updateData["code_batch"]; exists {
		if newCodeBatch == nil {
			if oldBatch.CodeBatch != nil {
				changes = append(changes, "removed code batch")
			}
		} else {
			newCode := newCodeBatch.(*string)
			if oldBatch.CodeBatch == nil {
				changes = append(changes, "added code batch: "+*newCode)
			} else if *oldBatch.CodeBatch != *newCode {
				changes = append(changes, "changed code batch from "+*oldBatch.CodeBatch+" to "+*newCode)
			}
		}
	}

	if newUnitPrice, exists := updateData["unit_price"]; exists {
		if newUnitPrice == nil {
			if oldBatch.UnitPrice != nil {
				changes = append(changes, "removed unit price")
			}
		} else {
			newPrice := newUnitPrice.(*float64)
			if oldBatch.UnitPrice == nil {
				changes = append(changes, "added unit price: "+u.formatFloat(*newPrice))
			} else if *oldBatch.UnitPrice != *newPrice {
				changes = append(changes, "changed unit price from "+u.formatFloat(*oldBatch.UnitPrice)+" to "+u.formatFloat(*newPrice))
			}
		}
	}

	if newExpDate, exists := updateData["exp_date"]; exists {
		if !newExpDate.(time.Time).Equal(oldBatch.ExpDate) {
			changes = append(changes, "changed expiry date from "+oldBatch.ExpDate.Format("2006-01-02")+" to "+newExpDate.(time.Time).Format("2006-01-02"))
		}
	}

	if newDescription, exists := updateData["description"]; exists {
		if newDescription == nil {
			if oldBatch.Description != nil {
				changes = append(changes, "removed description")
			}
		} else {
			newDesc := newDescription.(*string)
			if oldBatch.Description == nil {
				changes = append(changes, "added description")
			} else if *oldBatch.Description != *newDesc {
				changes = append(changes, "updated description")
			}
		}
	}

	if len(changes) == 0 {
		return "Product batch updated (no field changes detected)"
	}

	description := "Product batch updated: "
	for i, change := range changes {
		if i > 0 {
			description += ", "
		}
		description += change
	}

	return description
}

// GenerateDeleteDescription creates description for product batch deletion
func (u *ProductBatchTrackingUtils) GenerateDeleteDescription(productBatch model.ProductBatch) string {
	description := "Product batch deleted"
	if productBatch.CodeBatch != nil {
		description += " (code: " + *productBatch.CodeBatch + ")"
	}
	return description
}

// formatFloat formats float values for display
func (u *ProductBatchTrackingUtils) formatFloat(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

// NewProductBatchTrackingUtils creates a new instance of ProductBatchTrackingUtils
func NewProductBatchTrackingUtils() *ProductBatchTrackingUtils {
	return &ProductBatchTrackingUtils{}
}

// Standalone helper functions for backward compatibility and easier usage

// GenerateCreateDescription creates description for product batch creation (standalone function)
func GenerateCreateDescription(productBatch model.ProductBatch) string {
	utils := NewProductBatchTrackingUtils()
	return utils.GenerateCreateDescription(productBatch)
}

// GenerateUpdateDescriptionFromChanges creates description based on actual changes made (standalone function)
func GenerateUpdateDescriptionFromChanges(updateData map[string]interface{}, oldBatch model.ProductBatch) string {
	utils := NewProductBatchTrackingUtils()
	return utils.GenerateUpdateDescriptionFromChanges(updateData, oldBatch)
}

// GenerateDeleteDescription creates description for product batch deletion (standalone function)
func GenerateDeleteDescription(productBatch model.ProductBatch) string {
	utils := NewProductBatchTrackingUtils()
	return utils.GenerateDeleteDescription(productBatch)
}
