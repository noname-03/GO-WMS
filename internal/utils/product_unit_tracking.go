package utils

import (
	"fmt"
	"myapp/internal/model"
)

// ProductUnitTrackingUtils provides utility functions for tracking ProductUnit changes
type ProductUnitTrackingUtils struct{}

// GenerateCreateDescription creates description for product unit creation
func (u *ProductUnitTrackingUtils) GenerateCreateDescriptionProductUnit(productUnit model.ProductUnit) string {
	description := "Product unit created"
	if productUnit.Name != nil {
		description += " with name: " + *productUnit.Name
	}
	if productUnit.Quantity != nil {
		description += " with Quantity: " + u.formatFloat(*productUnit.Quantity)
	}
	if productUnit.UnitPrice != nil {
		description += ", unit price: " + u.formatFloat(*productUnit.UnitPrice)
	}
	description += ", barcode: " + *productUnit.Barcode

	return description
}

// GenerateUpdateDescriptionFromChanges creates description based on actual changes made
func (u *ProductUnitTrackingUtils) GenerateUpdateDescriptionFromChangesProductUnit(updateData map[string]interface{}, oldBatch model.ProductUnit) string {
	var changes []string

	// Check each field that was updated
	if newProductID, exists := updateData["product_id"]; exists {
		if newProductID != oldBatch.ProductID {
			changes = append(changes, fmt.Sprintf("changed product from ID %d to ID %d", oldBatch.ProductID, newProductID))
		}
	}

	if newName, exists := updateData["name"]; exists {
		if newName == nil {
			if oldBatch.Name != nil {
				changes = append(changes, "removed name: "+*oldBatch.Name)
			}
		} else {
			newNameStr := newName.(*string)
			if oldBatch.Name == nil {
				changes = append(changes, "added name: "+*newNameStr)
			} else if *oldBatch.Name != *newNameStr {
				changes = append(changes, "changed name from "+*oldBatch.Name+" to "+*newNameStr)
			}
		}
	}

	if newUnitQuantity, exists := updateData["quantity"]; exists {
		if newUnitQuantity == nil {
			if oldBatch.Quantity != nil {
				changes = append(changes, "removed quantity")
			}
		} else {
			newQuantity := newUnitQuantity.(*float64)
			if oldBatch.Quantity == nil {
				changes = append(changes, "added quantity: "+u.formatFloat(*newQuantity))
			} else if *oldBatch.Quantity != *newQuantity {
				changes = append(changes, "changed quantity from "+u.formatFloat(*oldBatch.Quantity)+" to "+u.formatFloat(*newQuantity))
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

	if newUnitBarcode, exists := updateData["barcode"]; exists {
		if newUnitBarcode == nil {
			if oldBatch.Barcode != nil {
				changes = append(changes, "removed barcode: "+*oldBatch.Barcode)
			}
		} else {
			newBarcode := newUnitBarcode.(*string)
			if oldBatch.Barcode == nil {
				changes = append(changes, "added barcode: "+*newBarcode)
			} else if *oldBatch.Barcode != *newBarcode {
				changes = append(changes, "changed barcode from "+*oldBatch.Barcode+" to "+*newBarcode)
			}
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
func (u *ProductUnitTrackingUtils) GenerateDeleteDescriptionProductUnit(productUnit model.ProductUnit) string {
	description := "Product batch deleted"
	if productUnit.Barcode != nil {
		description += " (code: " + *productUnit.Barcode + ")"
	}
	return description
}

// formatFloat formats float values for display
func (u *ProductUnitTrackingUtils) formatFloat(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

// NewProductUnitTrackingUtils creates a new instance of ProductUnitTrackingUtils
func NewProductUnitTrackingUtils() *ProductUnitTrackingUtils {
	return &ProductUnitTrackingUtils{}
}

// Standalone helper functions for backward compatibility and easier usage

// GenerateCreateDescription creates description for product batch creation (standalone function)
func GenerateCreateDescriptionProductUnit(productUnit model.ProductUnit) string {
	utils := NewProductUnitTrackingUtils()
	return utils.GenerateCreateDescriptionProductUnit(productUnit)
}

// GenerateUpdateDescriptionFromChanges creates description based on actual changes made (standalone function)
func GenerateUpdateDescriptionFromChangesProductUnit(updateData map[string]interface{}, oldBatch model.ProductUnit) string {
	utils := NewProductUnitTrackingUtils()
	return utils.GenerateUpdateDescriptionFromChangesProductUnit(updateData, oldBatch)
}

// GenerateDeleteDescription creates description for product batch deletion (standalone function)
func GenerateDeleteDescriptionProductUnit(productUnit model.ProductUnit) string {
	utils := NewProductUnitTrackingUtils()
	return utils.GenerateDeleteDescriptionProductUnit(productUnit)
}
