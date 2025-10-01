# Script untuk Menambahkan Deleted & Restore Endpoints

## Models yang perlu diupdate:

1. ✅ Brand - DONE
2. ⚙️ Category - In Progress  
3. ⏳ Product
4. ⏳ ProductBatch
5. ⏳ ProductUnit
6. ⏳ ProductStock
7. ⏳ ProductItem
8. ⏳ Location
9. ⏳ User
10. ⏳ Role

## Pattern untuk setiap model:

### Repository (tambahkan 2 method):
```go
// GetDeleted{ModelName}s returns all soft deleted records
func (r *{ModelName}Repository) GetDeleted{ModelName}s() ([]responseStruct, error) {
    var records []responseStruct
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

### Service (tambahkan 2 method):
```go
// GetDeleted{ModelName}s returns all soft deleted records
func (s *{ModelName}Service) GetDeleted{ModelName}s() ([]model.{ModelName}, error) {
    return s.{modelName}Repo.GetDeleted{ModelName}s()
}

// Restore{ModelName} restores a soft deleted record
func (s *{ModelName}Service) Restore{ModelName}(id uint, userID uint) (*model.{ModelName}, error) {
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
    return &restored{ModelName}, nil
}
```

### Handler (tambahkan 2 function):
```go
func GetDeleted{ModelName}s(c *fiber.Ctx) error {
    log.Printf("[{MODELNAME}] Get deleted {modelName}s request from IP: %s", c.IP())

    {modelName}s, err := {modelName}Service.GetDeleted{ModelName}s()
    if err != nil {
        log.Printf("[{MODELNAME}] Get deleted {modelName}s failed - error: %v", err)
        return helper.Fail(c, 500, "Failed to fetch deleted {modelName}s", err.Error())
    }

    log.Printf("[{MODELNAME}] Get deleted {modelName}s successful - Found %d deleted {modelName}s", len({modelName}s))
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

    log.Printf("[{MODELNAME}] Restore {modelName} successful - {ModelName} ID: %d, Restored by User ID: %d", {modelName}.ID, userID)
    return helper.Success(c, 200, "{ModelName} restored successfully", {modelName})
}
```

### Routes (tambahkan 2 route):
```go
// Specific routes (no parameters)
{modelName}s.Get("/", handler.Get{ModelName}s)
{modelName}s.Get("/deleted", handler.GetDeleted{ModelName}s) // ✨ NEW

// Parameterized routes (MUST be at the end)
{modelName}s.Put("/:id/restore", handler.Restore{ModelName}) // ✨ NEW
```

## Endpoints yang akan dibuat:

1. `GET /api/v1/brands/deleted` - Get deleted brands
2. `PUT /api/v1/brands/:id/restore` - Restore brand
3. `GET /api/v1/categories/deleted` - Get deleted categories  
4. `PUT /api/v1/categories/:id/restore` - Restore category
5. ... (dan seterusnya untuk semua model)