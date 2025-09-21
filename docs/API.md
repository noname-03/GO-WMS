# 📚 API Documentation

Complete API endpoints documentation for GO-WMS (Warehouse Management System).

## Base URL
```
http://localhost:8080
```

## Authentication
All protected endpoints require JWT token in Authorization header:
```
Authorization: Bearer <jwt_token>
```

## 🔐 Authentication Endpoints

### Login
```http
POST /api/v1/auth/login
```

**Request Body:**
```json
{
  "email": "alice@mail.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "Alice Johnson",
      "email": "alice@mail.com"
    }
  }
}
```

### Register
```http
POST /api/v1/auth/register
```

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@mail.com",
  "password": "password123"
}
```

### Get Profile
```http
GET /api/v1/auth/profile
```
*Protected endpoint*

## 👥 User Management

### Get All Users
```http
GET /api/v1/users
```
*Protected endpoint*

### Get Users with Minimal Data
```http
GET /api/v1/users/minimal
```
*Protected endpoint*

### Search Users
```http
GET /api/v1/users/search?q=keyword
```
*Protected endpoint*

### Get User Statistics
```http
GET /api/v1/users/stats
```
*Protected endpoint*

### Get User by ID
```http
GET /api/v1/users/:id
```
*Protected endpoint*

### Create User
```http
POST /api/v1/users
```
*Protected endpoint*

**Request Body:**
```json
{
  "name": "New User",
  "email": "newuser@mail.com",
  "password": "password123"
}
```

### Update User
```http
PUT /api/v1/users/:id
```
*Protected endpoint*

### Delete User
```http
DELETE /api/v1/users/:id
```
*Protected endpoint*

## 🎭 Role Management

### Get All Roles
```http
GET /api/v1/roles
```
*Protected endpoint*

### Get Role by ID
```http
GET /api/v1/roles/:id
```
*Protected endpoint*

### Create Role
```http
POST /api/v1/roles
```
*Protected endpoint*

**Request Body:**
```json
{
  "name": "Manager",
  "description": "Manager role with extended permissions"
}
```

### Update Role
```http
PUT /api/v1/roles/:id
```
*Protected endpoint*

### Delete Role
```http
DELETE /api/v1/roles/:id
```
*Protected endpoint*

## 🏷️ Brand Management

### Get All Brands
```http
GET /api/v1/brands
```
*Protected endpoint*

### Get Brand by ID
```http
GET /api/v1/brands/:id
```
*Protected endpoint*

### Create Brand
```http
POST /api/v1/brands
```
*Protected endpoint*

**Request Body:**
```json
{
  "name": "Toyota",
  "description": "Japanese automotive manufacturer"
}
```

### Update Brand
```http
PUT /api/v1/brands/:id
```
*Protected endpoint*

### Delete Brand
```http
DELETE /api/v1/brands/:id
```
*Protected endpoint*

## 📂 Category Management

### Get All Categories
```http
GET /api/v1/categories
```
*Protected endpoint*

### Get Category by ID
```http
GET /api/v1/categories/:id
```
*Protected endpoint*

### Create Category
```http
POST /api/v1/categories
```
*Protected endpoint*

**Request Body:**
```json
{
  "brandId": 1,
  "name": "Automotive",
  "description": "Automotive category"
}
```

### Update Category
```http
PUT /api/v1/categories/:id
```
*Protected endpoint*

### Delete Category
```http
DELETE /api/v1/categories/:id
```
*Protected endpoint*

## 📦 Product Management

### Get All Products
```http
GET /api/v1/products
```
*Protected endpoint*

### Get Product by ID
```http
GET /api/v1/products/:id
```
*Protected endpoint*

### Create Product
```http
POST /api/v1/products
```
*Protected endpoint*

**Request Body:**
```json
{
  "categoryId": 1,
  "name": "Toyota Camry",
  "description": "Mid-size sedan"
}
```

### Update Product
```http
PUT /api/v1/products/:id
```
*Protected endpoint*

### Delete Product
```http
DELETE /api/v1/products/:id
```
*Protected endpoint*

### Get Product Batches by Product
```http
GET /api/v1/products/:productId/batches
```
*Protected endpoint*

## 📋 Product Batch Management

### Get All Product Batches
```http
GET /api/v1/product-batches
```
*Protected endpoint*

### Get Product Batch by ID
```http
GET /api/v1/product-batches/:id
```
*Protected endpoint*

### Create Product Batch
```http
POST /api/v1/product-batches
```
*Protected endpoint*

**Request Body:**
```json
{
  "productId": 1,
  "codeBatch": "BATCH-CAM-2024-001",
  "unitPrice": 35000.00,
  "expDate": "2024-12-31",
  "description": "Toyota Camry batch Q1 2024"
}
```

### Update Product Batch
```http
PUT /api/v1/product-batches/:id
```
*Protected endpoint*

### Delete Product Batch
```http
DELETE /api/v1/product-batches/:id
```
*Protected endpoint*

## 🏥 Health Check

### Global Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "OK",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### API v1 Health Check
```http
GET /api/v1/health
```

## 📊 Response Format

### Success Response
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    // Response data here
  }
}
```

### Error Response
```json
{
  "code": 400,
  "message": "Error message",
  "error": "Detailed error information"
}
```

## 🔍 HTTP Status Codes

| Code | Description |
|------|-------------|
| 200  | OK - Request successful |
| 201  | Created - Resource created successfully |
| 400  | Bad Request - Invalid request data |
| 401  | Unauthorized - Authentication required |
| 403  | Forbidden - Insufficient permissions |
| 404  | Not Found - Resource not found |
| 409  | Conflict - Resource already exists |
| 500  | Internal Server Error - Server error |

## 🔐 Authentication Flow

1. **Login** with email and password
2. Receive **JWT token** in response
3. Include token in **Authorization header** for protected endpoints
4. Token expires after configured time (default: 24 hours)

## 📝 Request Examples

### Using cURL

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@mail.com",
    "password": "password123"
  }'
```

**Get Products (with token):**
```bash
curl -X GET http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Create Category:**
```bash
curl -X POST http://localhost:8080/api/v1/categories \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "brandId": 1,
    "name": "Electronics",
    "description": "Electronic devices and components"
  }'
```

## 🔄 Pagination

Most list endpoints support pagination using query parameters:

```http
GET /api/v1/products?page=1&limit=10
```

## 🔍 Search and Filtering

Search functionality is available on specific endpoints:

```http
GET /api/v1/users/search?q=alice
GET /api/v1/products/search?q=toyota
```

## ⚠️ Rate Limiting

API endpoints may be rate-limited. Respect the following headers:
- `X-RateLimit-Limit`: Maximum requests per window
- `X-RateLimit-Remaining`: Remaining requests in current window
- `X-RateLimit-Reset`: Time when the rate limit window resets