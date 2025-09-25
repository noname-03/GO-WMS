# GO-WMS API Endpoints - Stock Management System

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
All endpoints require JWT token in Authorization header:
```
Authorization: Bearer <jwt_token>
```

---

## 🏗️ ProductStock Endpoints

### 1. Get All Product Stocks
- **GET** `/product-stocks`
- **Description**: Retrieve all product stocks with product name, batch code, and location name
- **Response**: Array of product stocks with joined data

### 2. Get Product Stock by ID
- **GET** `/product-stocks/:id`
- **Description**: Retrieve specific product stock by ID
- **Parameters**: 
  - `id` (path): Product stock ID

### 3. Get Product Stocks by Product
- **GET** `/product-stocks/product/:productId`
- **Description**: Retrieve all stocks for a specific product
- **Parameters**: 
  - `productId` (path): Product ID

### 4. Create Product Stock
- **POST** `/product-stocks`
- **Description**: Create new product stock entry
- **Request Body**:
```json
{
  "productBatchId": 1,
  "productId": 1,
  "locationId": 1,        // optional
  "quantity": 100.0       // optional, must be >= 0
}
```

### 5. Update Product Stock
- **PUT** `/product-stocks/:id`
- **Description**: Update existing product stock
- **Parameters**: 
  - `id` (path): Product stock ID
- **Request Body**:
```json
{
  "quantity": 150.0       // optional, must be >= 0
}
```

### 6. Delete Product Stock
- **DELETE** `/product-stocks/:id`
- **Description**: Soft delete product stock
- **Parameters**: 
  - `id` (path): Product stock ID

---

## 📊 ProductStockTrack Endpoints

### 1. Get All Product Stock Tracks
- **GET** `/product-stock-tracks`
- **Description**: Retrieve all stock tracking records

### 2. Get Product Stock Track by ID
- **GET** `/product-stock-tracks/:id`
- **Description**: Retrieve specific stock track by ID
- **Parameters**: 
  - `id` (path): Stock track ID

### 3. Get Stock Tracks by Stock
- **GET** `/product-stock-tracks/stock/:stockId`
- **Description**: Retrieve all tracking records for a specific stock
- **Parameters**: 
  - `stockId` (path): Product stock ID

### 4. Get Stock Tracks by Product
- **GET** `/product-stock-tracks/product/:productId`
- **Description**: Retrieve all tracking records for a specific product
- **Parameters**: 
  - `productId` (path): Product ID

### 5. Get Stock Tracks by Date Range
- **GET** `/product-stock-tracks/date-range?startDate=YYYY-MM-DD&endDate=YYYY-MM-DD`
- **Description**: Retrieve stock tracks within date range
- **Query Parameters**: 
  - `startDate`: Start date (YYYY-MM-DD format)
  - `endDate`: End date (YYYY-MM-DD format)

### 6. Create Product Stock Track
- **POST** `/product-stock-tracks`
- **Description**: Create new stock tracking record
- **Request Body**:
```json
{
  "productStockId": 1,
  "productId": 1,
  "productBatchId": 1,
  "dateTrack": "2025-09-25T10:00:00Z",
  "quantity": 10.0,       // optional, must be > 0
  "operation": "Plus",    // "Plus" or "Minus"
  "stock": 110.0          // optional, will be calculated if not provided
}
```

### 7. Update Product Stock Track
- **PUT** `/product-stock-tracks/:id`
- **Description**: Update existing stock track
- **Parameters**: 
  - `id` (path): Stock track ID
- **Request Body**:
```json
{
  "dateTrack": "2025-09-25T10:00:00Z",
  "quantity": 15.0,       // optional, must be > 0
  "operation": "Plus",    // optional, "Plus" or "Minus"
  "stock": 125.0          // optional, must be >= 0
}
```

### 8. Delete Product Stock Track
- **DELETE** `/product-stock-tracks/:id`
- **Description**: Soft delete stock track
- **Parameters**: 
  - `id` (path): Stock track ID

---

## 📦 ProductItem Endpoints

### 1. Get All Product Items
- **GET** `/product-items`
- **Description**: Retrieve all product items with related data

### 2. Get Product Item by ID
- **GET** `/product-items/:id`
- **Description**: Retrieve specific product item by ID
- **Parameters**: 
  - `id` (path): Product item ID

### 3. Get Product Items by Stock
- **GET** `/product-items/stock/:stockId`
- **Description**: Retrieve all items for a specific stock
- **Parameters**: 
  - `stockId` (path): Product stock ID

### 4. Get Product Items by Product
- **GET** `/product-items/product/:productId`
- **Description**: Retrieve all items for a specific product
- **Parameters**: 
  - `productId` (path): Product ID

### 5. Get Product Items by Location
- **GET** `/product-items/location/:locationId`
- **Description**: Retrieve all items for a specific location
- **Parameters**: 
  - `locationId` (path): Location ID

### 6. Get Items Summary by Product
- **GET** `/product-items/summary/by-product`
- **Description**: Get summary of items grouped by product (total stock in, out, quantity)

### 7. Create Product Item
- **POST** `/product-items`
- **Description**: Create new product item
- **Request Body**:
```json
{
  "productStockId": 1,
  "productId": 1,
  "productBatchId": 1,
  "locationId": 1,        // optional
  "stockIn": 50.0,        // optional, must be >= 0
  "stockOut": 10.0,       // optional, must be >= 0
  "quantity": 40.0        // optional, calculated from stockIn - stockOut
}
```

### 8. Update Product Item
- **PUT** `/product-items/:id`
- **Description**: Update existing product item
- **Parameters**: 
  - `id` (path): Product item ID
- **Request Body**:
```json
{
  "stockIn": 60.0,        // optional, must be >= 0
  "stockOut": 15.0,       // optional, must be >= 0
  "quantity": 45.0        // optional, will be recalculated if not provided
}
```

### 9. Delete Product Item
- **DELETE** `/product-items/:id`
- **Description**: Soft delete product item
- **Parameters**: 
  - `id` (path): Product item ID

---

## 🔍 ProductItemTrack Endpoints

### 1. Get All Product Item Tracks
- **GET** `/product-item-tracks`
- **Description**: Retrieve all item tracking records

### 2. Get Product Item Track by ID
- **GET** `/product-item-tracks/:id`
- **Description**: Retrieve specific item track by ID
- **Parameters**: 
  - `id` (path): Item track ID

### 3. Get Item Tracks by Item
- **GET** `/product-item-tracks/item/:itemId`
- **Description**: Retrieve all tracking records for a specific item
- **Parameters**: 
  - `itemId` (path): Product item ID

### 4. Get Item Tracks by Stock
- **GET** `/product-item-tracks/stock/:stockId`
- **Description**: Retrieve all item tracks for a specific stock
- **Parameters**: 
  - `stockId` (path): Product stock ID

### 5. Get Item Tracks by Product
- **GET** `/product-item-tracks/product/:productId`
- **Description**: Retrieve all item tracks for a specific product
- **Parameters**: 
  - `productId` (path): Product ID

### 6. Get Item Tracks by Date Range
- **GET** `/product-item-tracks/date-range?startDate=YYYY-MM-DD&endDate=YYYY-MM-DD`
- **Description**: Retrieve item tracks within date range
- **Query Parameters**: 
  - `startDate`: Start date (YYYY-MM-DD format)
  - `endDate`: End date (YYYY-MM-DD format)

### 7. Get Tracks by Operation
- **GET** `/product-item-tracks/operation/:operation`
- **Description**: Retrieve tracks filtered by operation type
- **Parameters**: 
  - `operation` (path): Operation type ("In", "Out", "Plus", "Minus")

### 8. Get Value Report by Product
- **GET** `/product-item-tracks/reports/value-by-product`
- **Description**: Get value report grouped by product (total transactions, value, avg unit price)

### 9. Create Product Item Track
- **POST** `/product-item-tracks`
- **Description**: Create new item tracking record
- **Request Body**:
```json
{
  "productItemId": 1,
  "productStockId": 1,
  "productId": 1,
  "productBatchId": 1,
  "dateTrack": "2025-09-25T10:00:00Z",
  "unitPrice": 15000.0,   // optional, must be > 0
  "quantity": 5.0,        // optional, must be > 0
  "operation": "In",      // "In", "Out", "Plus", or "Minus"
  "stock": 45.0           // optional, will be calculated if not provided
}
```

### 10. Update Product Item Track
- **PUT** `/product-item-tracks/:id`
- **Description**: Update existing item track
- **Parameters**: 
  - `id` (path): Item track ID
- **Request Body**:
```json
{
  "dateTrack": "2025-09-25T10:00:00Z",
  "unitPrice": 18000.0,   // optional, must be > 0
  "quantity": 8.0,        // optional, must be > 0
  "operation": "Out",     // optional, "In", "Out", "Plus", or "Minus"
  "stock": 37.0           // optional, must be >= 0
}
```

### 11. Delete Product Item Track
- **DELETE** `/product-item-tracks/:id`
- **Description**: Soft delete item track
- **Parameters**: 
  - `id` (path): Item track ID

---

## 📝 Notes

### Operation Types
- **Plus/Minus**: For stock adjustments (ProductStockTrack)
- **In/Out**: For item movements (ProductItemTrack)
- **Plus/Minus**: Also available for item movements (ProductItemTrack)

### Stock Calculations
- Stock values are automatically calculated based on previous stock and quantity
- For insufficient stock operations, the system will prevent negative stock
- All tracking operations maintain audit trail with user information

### Response Format
All endpoints return standardized response format:
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... }
}
```

### Error Responses
Error responses follow this format:
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error information"
}
```

### Total Endpoints: 32
- **ProductStock**: 6 endpoints
- **ProductStockTrack**: 8 endpoints  
- **ProductItem**: 9 endpoints
- **ProductItemTrack**: 11 endpoints (including 2 reporting endpoints)