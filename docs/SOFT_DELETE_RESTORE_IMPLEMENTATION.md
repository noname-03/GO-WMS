## ✅ SOFT DELETE & RESTORE ENDPOINTS - IMPLEMENTASI LENGKAP

### 🎯 Status Implementasi

| Model | Repository | Service | Handler | Routes | Status |
|-------|------------|---------|---------|--------|--------|
| ✅ Brand | ✅ GetDeletedBrands<br>✅ RestoreBrand | ✅ GetDeletedBrands<br>✅ RestoreBrand | ✅ GetDeletedBrands<br>✅ RestoreBrand | ✅ GET /deleted<br>✅ PUT /:id/restore | **COMPLETE** |
| ✅ Category | ✅ GetDeletedCategories<br>✅ RestoreCategory | ✅ GetDeletedCategories<br>✅ RestoreCategory | ✅ GetDeletedCategories<br>✅ RestoreCategory | ✅ GET /deleted<br>✅ PUT /:id/restore | **COMPLETE** |
| ✅ Product | ✅ GetDeletedProducts<br>✅ RestoreProduct | ✅ GetDeletedProducts<br>✅ RestoreProduct | ✅ GetDeletedProducts<br>✅ RestoreProduct | ✅ GET /deleted<br>✅ PUT /:id/restore | **COMPLETE** |
| ✅ ProductBatch | ✅ GetDeletedProductBatches<br>✅ RestoreProductBatch | ✅ GetDeletedProductBatches<br>✅ RestoreProductBatch | ✅ GetDeletedProductBatches<br>✅ RestoreProductBatch | ✅ GET /deleted<br>✅ PUT /:id/restore | **COMPLETE** |
| ✅ ProductUnit | ✅ GetDeletedProductUnits<br>✅ RestoreProductUnit | ✅ GetDeletedProductUnits<br>✅ RestoreProductUnit | ✅ GetDeletedProductUnits<br>✅ RestoreProductUnit | ✅ GET /deleted<br>✅ PUT /:id/restore | **COMPLETE** |
| ✅ Location | ✅ GetDeletedLocations<br>✅ RestoreLocation | ✅ GetDeletedLocations<br>✅ RestoreLocation | ✅ GetDeletedLocations<br>✅ RestoreLocation | ✅ GET /deleted<br>✅ PUT /:id/restore | **COMPLETE** |
| ✅ User | ✅ GetDeletedUsers<br>✅ RestoreUser | ✅ GetDeletedUsers<br>✅ RestoreUser | ✅ GetDeletedUsers<br>✅ RestoreUser | ✅ GET /deleted<br>✅ PUT /:id/restore | **COMPLETE** |
| ✅ Role | ✅ GetDeletedRoles<br>✅ RestoreRole | ✅ GetDeletedRoles<br>✅ RestoreRole | ✅ GetDeletedRoles<br>✅ RestoreRole | ✅ GET /deleted<br>✅ PUT /:id/restore | **COMPLETE** |

### 🌐 Endpoints Yang Tersedia

#### 1. Brand
- `GET /api/v1/brands/deleted` - Melihat semua brand yang sudah dihapus
- `PUT /api/v1/brands/:id/restore` - Restore brand yang sudah dihapus

#### 2. Category  
- `GET /api/v1/categories/deleted` - Melihat semua category yang sudah dihapus
- `PUT /api/v1/categories/:id/restore` - Restore category yang sudah dihapus

#### 3. Product
- `GET /api/v1/products/deleted` - Melihat semua product yang sudah dihapus
- `PUT /api/v1/products/:id/restore` - Restore product yang sudah dihapus

#### 4. Product Batch
- `GET /api/v1/product-batches/deleted` - Melihat semua product batch yang sudah dihapus
- `PUT /api/v1/product-batches/:id/restore` - Restore product batch yang sudah dihapus

#### 5. Product Unit
- `GET /api/v1/product-units/deleted` - Melihat semua product unit yang sudah dihapus
- `PUT /api/v1/product-units/:id/restore` - Restore product unit yang sudah dihapus

#### 6. Location
- `GET /api/v1/locations/deleted` - Melihat semua location yang sudah dihapus
- `PUT /api/v1/locations/:id/restore` - Restore location yang sudah dihapus

#### 7. User
- `GET /api/v1/users/deleted` - Melihat semua user yang sudah dihapus
- `PUT /api/v1/users/:id/restore` - Restore user yang sudah dihapus

#### 8. Role
- `GET /api/v1/roles/deleted` - Melihat semua role yang sudah dihapus
- `PUT /api/v1/roles/:id/restore` - Restore role yang sudah dihapus

### 🔐 Authentication & Authorization

Semua endpoint memerlukan:
- **JWT Authentication** - Bearer token wajib
- **User ID** untuk audit trail (diambil dari JWT token)

### 📋 Response Format

#### Get Deleted Items Response:
```json
{
  "status": "success",
  "status_code": 200,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "Deleted Item Name",
      // ... other fields sesuai dengan format GetAll masing-masing model
    }
  ]
}
```

#### Restore Item Response:
```json
{
  "status": "success", 
  "status_code": 200,
  "message": "{ModelName} restored successfully",
  "data": {
    "id": 1,
    "name": "Restored Item Name",
    // ... other fields sesuai dengan format GetByID masing-masing model
  }
}
```

### ⚙️ Technical Implementation

#### Repository Layer Pattern:
```go
// GetDeleted{ModelName}s returns all soft deleted records
func (r *{ModelName}Repository) GetDeleted{ModelName}s() ([]responseType, error) {
    var records []responseType
    result := database.DB.Unscoped().Where("deleted_at IS NOT NULL").Order("deleted_at DESC").Find(&records)
    return records, result.Error
}

// Restore{ModelName} restores a soft deleted record
func (r *{ModelName}Repository) Restore{ModelName}(id uint, userID uint) error {
    updateData := map[string]interface{}{
        "user_updt":  userID,
        "deleted_at": nil,
    }
    return database.DB.Unscoped().Model(&model.{ModelName}{}).Where("id = ?", id).Updates(updateData).Error
}
```

#### Service Layer Pattern:
```go
// GetDeleted{ModelName}s returns all soft deleted records
func (s *{ModelName}Service) GetDeleted{ModelName}s() (interface{}, error) {
    return s.{modelName}Repo.GetDeleted{ModelName}s()
}

// Restore{ModelName} restores a soft deleted record with validation
func (s *{ModelName}Service) Restore{ModelName}(id uint, userID uint) (interface{}, error) {
    if id == 0 {
        return nil, errors.New("invalid {modelName} ID")
    }
    if userID == 0 {
        return nil, errors.New("user ID is required for audit trail")
    }
    
    err := s.{modelName}Repo.Restore{ModelName}(id, userID)
    if err != nil {
        return nil, err
    }
    
    restored{ModelName}, err := s.{modelName}Repo.Get{ModelName}ByID(id)
    if err != nil {
        return nil, err
    }
    return restored{ModelName}, nil
}
```

#### Handler Layer Pattern:
```go
func GetDeleted{ModelName}s(c *fiber.Ctx) error {
    log.Printf("[{MODELNAME}] Get deleted {modelName}s request from IP: %s", c.IP())

    {modelName}s, err := {modelName}Service.GetDeleted{ModelName}s()
    if err != nil {
        log.Printf("[{MODELNAME}] Get deleted {modelName}s failed - error: %v", err)
        return helper.Fail(c, 500, "Failed to fetch deleted {modelName}s", err.Error())
    }

    log.Printf("[{MODELNAME}] Get deleted {modelName}s successful")
    return helper.Success(c, 200, "Success", {modelName}s)
}

func Restore{ModelName}(c *fiber.Ctx) error {
    id := c.Params("id")
    log.Printf("[{MODELNAME}] Restore {modelName} request - ID: %s from IP: %s", id, c.IP())

    idUint, err := strconv.ParseUint(id, 10, 32)
    if err != nil {
        log.Printf("[{MODELNAME}] Restore {modelName} failed - Invalid ID: %s, error: %v", id, err)
        return helper.Fail(c, 400, "Invalid {modelName} ID", err.Error())
    }

    userID, ok := c.Locals("user_id").(uint)
    if !ok {
        log.Printf("[{MODELNAME}] Restore {modelName} failed - User not authenticated for {ModelName} ID: %d", idUint)
        return helper.Fail(c, 401, "User not authenticated", "Failed to get user ID from token")
    }

    {modelName}, err := {modelName}Service.Restore{ModelName}(uint(idUint), userID)
    if err != nil {
        log.Printf("[{MODELNAME}] Restore {modelName} failed - {ModelName} ID: %d, error: %v", idUint, err)
        statusCode, message := handle{ModelName}Error(err)
        return helper.Fail(c, statusCode, message, err.Error())
    }

    log.Printf("[{MODELNAME}] Restore {modelName} successful - {ModelName} ID: %d, Restored by User ID: %d", idUint, userID)
    return helper.Success(c, 200, "{ModelName} restored successfully", {modelName})
}
```

### 🛡️ Security Features

1. **Audit Trail**: Setiap restore mencatat user yang melakukan restore
2. **JWT Protection**: Semua endpoint dilindungi JWT middleware
3. **Input Validation**: Validasi ID dan user authentication
4. **Error Handling**: Comprehensive error handling dengan logging
5. **SQL Injection Protection**: Menggunakan GORM parameterized queries

### 📊 Logging & Monitoring

Setiap operasi mencatat:
- User IP address
- User ID yang melakukan operasi  
- Timestamp operasi
- Status success/failure
- Error details (jika ada)

### 🎉 COMPLETION STATUS: 100% DONE!

Semua 7 model utama telah berhasil diimplementasikan dengan:
- ✅ Soft delete functionality dengan audit trail
- ✅ Comprehensive restore functionality  
- ✅ Consistent API patterns
- ✅ Proper error handling & logging
- ✅ JWT authentication protection
- ✅ Database relationship preservation
- ✅ Response format consistency

**Total Endpoints Added: 16 endpoints (8 GET deleted + 8 PUT restore)**