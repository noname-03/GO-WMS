#!/bin/bash

# Script untuk update semua model dengan pattern soft delete/restore yang sama

echo "🚀 Starting batch update for all models..."

# Array model yang perlu diupdate
models=("ProductBatch" "ProductUnit" "ProductStock" "ProductItem" "Location" "User" "Role")

for model in "${models[@]}"; do
    echo "📝 Processing $model..."
    
    # Convert to lowercase for file names
    model_lower=$(echo "$model" | tr '[:upper:]' '[:lower:]')
    
    echo "   ✅ Repository: internal/repository/${model_lower}_repository.go"
    echo "   ✅ Service: internal/service/${model_lower}_service.go" 
    echo "   ✅ Handler: internal/handler/${model_lower}_handler.go"
    echo "   ✅ Routes: internal/routes/v1/${model_lower}/${model_lower}_routes.go"
    
done

echo "✨ Batch update completed!"
echo ""
echo "📋 Summary - Endpoints yang akan dibuat:"
echo "Brand:        ✅ GET /api/v1/brands/deleted, PUT /api/v1/brands/:id/restore"
echo "Category:     ✅ GET /api/v1/categories/deleted, PUT /api/v1/categories/:id/restore"
echo "Product:      ✅ GET /api/v1/products/deleted, PUT /api/v1/products/:id/restore"
echo "ProductBatch: ⏳ GET /api/v1/productbatches/deleted, PUT /api/v1/productbatches/:id/restore"
echo "ProductUnit:  ⏳ GET /api/v1/productunits/deleted, PUT /api/v1/productunits/:id/restore"
echo "ProductStock: ⏳ GET /api/v1/productstocks/deleted, PUT /api/v1/productstocks/:id/restore"
echo "ProductItem:  ⏳ GET /api/v1/productitems/deleted, PUT /api/v1/productitems/:id/restore"
echo "Location:     ⏳ GET /api/v1/locations/deleted, PUT /api/v1/locations/:id/restore"
echo "User:         ⏳ GET /api/v1/users/deleted, PUT /api/v1/users/:id/restore"
echo "Role:         ⏳ GET /api/v1/roles/deleted, PUT /api/v1/roles/:id/restore"