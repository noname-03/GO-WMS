# 🏗️ Architecture Documentation

Architecture overview and design patterns for GO-WMS (Warehouse Management System).

## 📁 Project Structure

```
GO-WMS/
├── .air.toml              # Air configuration for hot reload
├── docker-compose.yml     # Docker services configuration
├── Dockerfile            # Application container
├── go.mod                # Go modules
├── go.sum                # Go modules checksum
├── main.go               # Application entry point
├── README.md             # Project documentation
├── docs/                 # Documentation files
│   ├── API.md            # API documentation
│   ├── ARCHITECTURE.md   # This file
│   ├── DEPLOYMENT.md     # Deployment guide
│   ├── DEVELOPMENT.md    # Development guide
│   ├── DATABASE.md       # Database documentation
│   └── TESTING.md        # Testing guide
├── database/             # Database related files
│   ├── connection.go     # Database connection setup
│   ├── migration.go      # Database auto-migration
│   ├── seed.go          # Seeder execution
│   └── seeder/          # Database seeders
│       ├── registry.go   # Seeder registry
│       ├── seeder.go     # Seeder interface
│       ├── userSeeder.go # User seeder
│       ├── roleSeeder.go # Role seeder
│       ├── brandSeeder.go # Brand seeder
│       ├── categorySeeder.go # Category seeder
│       ├── productSeeder.go # Product seeder
│       ├── productBatchSeeder.go # Product batch seeder
│       └── productBatchTrackSeeder.go # Product batch tracking seeder
├── internal/             # Private application code
│   ├── handler/          # HTTP handlers
│   │   ├── auth_handler.go # Authentication handlers
│   │   ├── user_handler.go # User CRUD handlers
│   │   ├── role_handler.go # Role CRUD handlers
│   │   ├── brand_handler.go # Brand CRUD handlers
│   │   ├── category_handler.go # Category CRUD handlers
│   │   ├── product_handler.go # Product CRUD handlers
│   │   └── product_batch_handler.go # Product batch CRUD handlers
│   ├── middleware/       # HTTP middlewares
│   │   └── auth.go       # JWT authentication middleware
│   ├── model/            # Data models (GORM)
│   │   ├── user.go       # User model with relationships
│   │   ├── role.go       # Role model
│   │   ├── brand.go      # Brand model
│   │   ├── category.go   # Category model with brand relationship
│   │   ├── product.go    # Product model with category relationship
│   │   ├── product_batch.go # Product batch model with product relationship
│   │   └── product_batch_track.go # Product batch tracking model
│   ├── repository/       # Data access layer
│   │   ├── user_repository.go # User repository
│   │   ├── role_repository.go # Role repository
│   │   ├── brand_repository.go # Brand repository
│   │   ├── category_repository.go # Category repository
│   │   ├── product_repository.go # Product repository
│   │   ├── product_batch_repository.go # Product batch repository
│   │   └── product_batch_track_repository.go # Product batch tracking repository
│   ├── routes/           # Route definitions
│   │   ├── routes.go     # Main route setup
│   │   └── v1/           # API version 1 routes
│   │       ├── v1.go     # V1 route registry
│   │       ├── auth/     # Authentication routes
│   │       ├── user/     # User routes
│   │       ├── role/     # Role routes
│   │       ├── brand/    # Brand routes
│   │       ├── category/ # Category routes
│   │       ├── product/  # Product routes
│   │       └── productbatch/ # Product batch routes
│   ├── service/          # Business logic layer
│   │   ├── user_service.go # User business logic
│   │   ├── role_service.go # Role business logic
│   │   ├── brand_service.go # Brand business logic
│   │   ├── category_service.go # Category business logic
│   │   ├── product_service.go # Product business logic
│   │   ├── product_batch_service.go # Product batch business logic
│   │   └── product_batch_track_service.go # Product batch tracking business logic
│   └── utils/            # Utility functions
│       ├── jwt.go        # JWT utilities
│       └── product_batch_tracking.go # Product batch tracking utilities
└── pkg/                  # Public library code
    └── helper/           # Helper functions
        ├── response.go   # Standardized API responses
        └── database_error.go # Database error handling
```

## 🎯 Design Patterns

### Clean Architecture
The application follows Clean Architecture principles with clear separation of concerns:

1. **Handler Layer** - HTTP request/response handling
2. **Service Layer** - Business logic implementation
3. **Repository Layer** - Data access abstraction
4. **Model Layer** - Data structures and entities

### Dependency Injection
- Services are injected into handlers
- Repositories are injected into services
- Database connections are shared across repositories

### Repository Pattern
- Abstract data access logic
- Easy to test and mock
- Consistent CRUD operations across entities

## 🔄 Request Flow

```
HTTP Request → Router → Middleware → Handler → Service → Repository → Database
                                      ↓
HTTP Response ← Helper ← Handler ← Service ← Repository ← Database
```

### Detailed Flow:
1. **Router** receives HTTP request
2. **Middleware** validates JWT token (for protected routes)
3. **Handler** parses request and validates input
4. **Service** implements business logic and validation
5. **Repository** performs database operations
6. **Helper** formats standardized responses

## 🛡️ Security Architecture

### Authentication
- **JWT-based authentication** with configurable expiration
- **Password hashing** using bcrypt
- **Token validation middleware** for protected routes

### Authorization
- **Role-based access control** (RBAC)
- **User role assignment** with different permission levels
- **Audit trails** for all entity changes

### Data Protection
- **Soft deletes** for data recovery
- **Audit fields** tracking who created/updated/deleted records
- **Input validation** and sanitization

## 📊 Data Layer

### ORM (GORM)
- **Auto-migration** for database schema
- **Relationship mapping** with foreign keys
- **Soft delete** support
- **Connection pooling** for performance

### Database Design
- **PostgreSQL** as primary database
- **Foreign key constraints** for data integrity
- **Indexes** for query optimization
- **Transactions** for data consistency

## 🔧 Middleware Stack

### Authentication Middleware
```go
// JWT validation for protected routes
func AuthMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Extract and validate JWT token
        // Set user information in context
        return c.Next()
    }
}
```

### Error Handling
- **Centralized error handling** with standard HTTP codes
- **Database error conversion** to user-friendly messages
- **Validation error formatting** with detailed feedback

### Logging
- **Request/response logging** for debugging
- **Audit trail logging** for security
- **Performance monitoring** with execution times

## 🌐 API Design

### RESTful Principles
- **Resource-based URLs** (`/api/v1/products`)
- **HTTP methods** for different operations (GET, POST, PUT, DELETE)
- **Status codes** indicating operation results
- **JSON communication** for requests and responses

### Versioning
- **URL versioning** (`/api/v1/`)
- **Backward compatibility** considerations
- **Future version planning**

### Response Format
```json
{
  "code": 200,
  "message": "Success",
  "data": { /* response data */ },
  "error": null
}
```

## 🔄 Entity Relationships

### Core Entities
- **User** ← has role → **Role**
- **Brand** ← has many → **Category**
- **Category** ← has many → **Product**  
- **Product** ← has many → **ProductBatch**
- **ProductBatch** ← has many → **ProductBatchTrack**

### Audit Relationships
All entities have audit relationships:
- `created_by` → **User**
- `updated_by` → **User**
- `deleted_by` → **User**

## 🔌 External Dependencies

### Core Dependencies
```go
// Web framework
"github.com/gofiber/fiber/v2"

// ORM
"gorm.io/gorm"
"gorm.io/driver/postgres"

// JWT authentication
"github.com/golang-jwt/jwt/v5"

// Password hashing
"golang.org/x/crypto/bcrypt"

// Environment variables
"github.com/joho/godotenv"
```

### Development Dependencies
```go
// Hot reload
"github.com/cosmtrek/air"

// Testing
"testing"
"github.com/stretchr/testify"
```

## 🚀 Performance Considerations

### Database Optimization
- **Connection pooling** for concurrent requests
- **Query optimization** with proper indexes
- **Relationship preloading** to avoid N+1 queries
- **Pagination** for large datasets

### Memory Management
- **Proper struct definitions** with appropriate field types
- **Garbage collection** considerations
- **Connection pooling** to prevent resource leaks

### Caching Strategy
- **In-memory caching** for frequently accessed data
- **Redis integration** for distributed caching (planned)
- **Cache invalidation** strategies

## 🧪 Testing Architecture

### Unit Testing
- **Service layer testing** with mocked repositories
- **Repository testing** with test database
- **Handler testing** with mocked services

### Integration Testing
- **API endpoint testing** with test database
- **Database transaction testing**
- **Authentication flow testing**

## 🔄 Development Workflow

### Hot Reload
```bash
# Air configuration for development
air
```

### Code Organization
- **Feature-based organization** within layers
- **Clear naming conventions**
- **Comprehensive documentation**

### Git Workflow
- **Feature branches** for new development
- **Pull request reviews** for code quality
- **Conventional commits** for clear history

## 📈 Scalability Considerations

### Horizontal Scaling
- **Stateless application design**
- **Database connection pooling**
- **Load balancer ready**

### Vertical Scaling
- **Efficient memory usage**
- **Optimized database queries**
- **Connection pool tuning**

### Future Enhancements
- **Microservice architecture** migration path
- **Event-driven architecture** for complex workflows
- **Message queue integration** for async processing