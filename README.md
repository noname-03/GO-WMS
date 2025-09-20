# GO-WMS (Warehouse Management System)

A comprehensive Warehouse Management System built with Go, Fiber, GORM, and PostgreSQL. Features complete CRUD operations for warehouse entities with audit trails, JWT authentication, and automatic tracking.

## 🚀 Features

### **Core Features**
- **JWT Authentication** - Secure user authentication with role-based access
- **RESTful API** - Clean API design with versioning (v1)
- **Database Migration** - Automatic database setup with seeders
- **Modular Architecture** - Clean separation of concerns (handler/service/repository)
- **Audit Trails** - Complete tracking of who created/updated/deleted records
- **Soft Deletes** - Safe deletion with recovery capability
- **Hot Reload Development** - Air integration for faster development

### **WMS Entities**
- **User Management** - Complete user CRUD with role assignment
- **Role Management** - Role-based access control system
- **Brand Management** - Product brand categorization
- **Category Management** - Hierarchical product categorization with brand relationships
- **Product Management** - Complete product catalog with category relationships
- **Product Batch Management** - Batch tracking with expiry dates and pricing
- **Product Batch Tracking** - Automatic history tracking of all batch changes

### **Advanced Features**
- **Relationship Management** - Full entity relationships with preloading
- **Comprehensive Logging** - Detailed request/response logging
- **Error Handling** - Standardized error responses with appropriate HTTP codes
- **Data Validation** - Input validation with detailed error messages
- **Search & Filtering** - Advanced search capabilities across entities
- **Change Tracking** - Detailed history of what changed, when, and by whom

## 📋 Prerequisites

- Go 1.23.12+
- PostgreSQL 16+
- Air v1.63.0+ (for hot reload development)
- Docker & Docker Compose (optional)

## 🛠️ Setup

### 1. Clone the repository
```bash
git clone https://github.com/noname-03/GO-WMS.git
cd GO-WMS
```

### 2. Environment Configuration
```bash
cp .env.example .env
# Edit .env file with your configuration
```

### 3. Database Setup (Docker)
```bash
docker-compose up -d db
```

### 4. Install Dependencies
```bash
go mod tidy
```

### 5. Install Air for Hot Reload (Development)
```bash
go install github.com/cosmtrek/air@latest
```

### 6. Run Application
```bash
# Production mode
go run main.go

# Development mode with hot reload
air
```

The application will start on `http://localhost:8080`

## 📚 API Documentation

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration  
- `GET /api/v1/auth/profile` - Get user profile (Protected)

### Users
- `GET /api/v1/users` - Get all users (Protected)
- `GET /api/v1/users/minimal` - Get users with minimal data (Protected)
- `GET /api/v1/users/search?q=keyword` - Search users (Protected)
- `GET /api/v1/users/stats` - Get user statistics (Protected)
- `GET /api/v1/users/:id` - Get user by ID (Protected)
- `POST /api/v1/users` - Create new user (Protected)
- `PUT /api/v1/users/:id` - Update user (Protected)
- `DELETE /api/v1/users/:id` - Delete user (Protected)

### Roles
- `GET /api/v1/roles` - Get all roles (Protected)
- `GET /api/v1/roles/:id` - Get role by ID (Protected)
- `POST /api/v1/roles` - Create new role (Protected)
- `PUT /api/v1/roles/:id` - Update role (Protected)
- `DELETE /api/v1/roles/:id` - Delete role (Protected)

### Brands
- `GET /api/v1/brands` - Get all brands (Protected)
- `GET /api/v1/brands/:id` - Get brand by ID (Protected)
- `POST /api/v1/brands` - Create new brand (Protected)
- `PUT /api/v1/brands/:id` - Update brand (Protected)
- `DELETE /api/v1/brands/:id` - Delete brand (Protected)

### Categories
- `GET /api/v1/categories` - Get all categories (Protected)
- `GET /api/v1/categories/:id` - Get category by ID (Protected)
- `POST /api/v1/categories` - Create new category (Protected)
- `PUT /api/v1/categories/:id` - Update category (Protected)
- `DELETE /api/v1/categories/:id` - Delete category (Protected)

### Products
- `GET /api/v1/products` - Get all products (Protected)
- `GET /api/v1/products/:id` - Get product by ID (Protected)
- `POST /api/v1/products` - Create new product (Protected)
- `PUT /api/v1/products/:id` - Update product (Protected)
- `DELETE /api/v1/products/:id` - Delete product (Protected)
- `GET /api/v1/products/:productId/batches` - Get batches by product (Protected)

### Product Batches
- `GET /api/v1/product-batches` - Get all product batches (Protected)
- `GET /api/v1/product-batches/:id` - Get product batch by ID (Protected)
- `POST /api/v1/product-batches` - Create new product batch (Protected)
- `PUT /api/v1/product-batches/:id` - Update product batch (Protected)
- `DELETE /api/v1/product-batches/:id` - Delete product batch (Protected)

### Health Check
- `GET /health` - Global health check
- `GET /api/v1/health` - API v1 health check

## 🏗️ Project Structure

```
GO-WMS/
├── .air.toml              # Air configuration for hot reload
├── docker-compose.yml     # Docker services configuration
├── Dockerfile            # Application container
├── go.mod                # Go modules
├── go.sum                # Go modules checksum
├── main.go               # Application entry point
├── README.md             # Project documentation
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

## 🔧 Development

### Default Login Credentials
```
Admin User:
Email: alice@mail.com
Password: password123
Role: Admin

Regular User:
Email: bob@mail.com  
Password: password123
Role: User
```

### Development Commands
```bash
# Start with hot reload (recommended for development)
air

# Start without hot reload
go run main.go

# Build application
go build -o bin/go-wms main.go

# Run with specific environment
ENV=development go run main.go
```

### Database Commands
```bash
# Start database
docker-compose up -d db

# View database logs
docker-compose logs db

# Stop database
docker-compose down

# Reset database (drop and recreate)
docker-compose down -v
docker-compose up -d db
```

### Sample Data
The application includes comprehensive seeders that create:
- **2 Roles**: Admin, User
- **2 Users**: Alice (Admin), Bob (User)  
- **3 Brands**: Toyota, Samsung, Nike
- **6 Categories**: Automotive, Electronics, Smartphones, Sports, Footwear, Apparel
- **4 Products**: Toyota Camry, Samsung Galaxy S24, Nike Air Zoom, iPhone 15 Pro
- **8 Product Batches**: With various expiry dates and pricing
- **Sample Tracking Records**: Showing create/update/delete history

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Test specific package
go test ./internal/service/...

# Test with verbose output
go test -v ./...
```

## 📊 Entity Relationships

```
User ──┐
       ├── ProductBatch (created_by, updated_by, deleted_by)
       ├── Product (created_by, updated_by, deleted_by)  
       ├── Category (created_by, updated_by, deleted_by)
       ├── Brand (created_by, updated_by, deleted_by)
       ├── Role (created_by, updated_by, deleted_by)
       └── ProductBatchTrack (user_inst, user_updt)

Role ──── User (role_id)

Brand ──── Category (brand_id)

Category ──── Product (category_id)

Product ──── ProductBatch (product_id)

ProductBatch ──── ProductBatchTrack (product_batch_id)
```

## 🔍 Key Features Explained

### **Audit Trails**
Every entity tracks:
- `created_by` - User who created the record
- `updated_by` - User who last updated the record  
- `deleted_by` - User who deleted the record
- `created_at`, `updated_at`, `deleted_at` - Timestamps

### **Product Batch Tracking**
Automatic tracking of ProductBatch changes:
- **Create**: "Product batch created with code: BATCH001, unit price: 150.00, expiry date: 2024-12-31"
- **Update**: "Product batch updated: changed unit price from 150.00 to 175.00, updated description"
- **Delete**: "Product batch deleted (code: BATCH001)"

### **Soft Deletes**
All entities use soft deletes - records are marked as deleted but remain in database for audit purposes.

### **Relationship Preloading**
API responses include related data:
- Products include Category and Brand information
- Categories include Brand information
- Product Batches include Product, Category, and Brand information
- Tracking records include ProductBatch and User information

## 📦 Building

```bash
# Build for current platform
go build -o bin/go-wms main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/go-wms-linux main.go

# Build for Windows  
GOOS=windows GOARCH=amd64 go build -o bin/go-wms.exe main.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o bin/go-wms-macos main.go

# Build with version info
go build -ldflags "-X main.version=1.0.0" -o bin/go-wms main.go
```

## 🚀 Deployment

### **Production Setup**
```bash
# 1. Build application
go build -o go-wms main.go

# 2. Set environment variables
export DB_HOST=your-postgres-host
export DB_PORT=5432
export DB_USER=your-db-user
export DB_PASSWORD=your-db-password
export DB_NAME=go_wms
export JWT_SECRET=your-jwt-secret

# 3. Run application
./go-wms
```

### **Environment Variables**
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=go_wms

# Application Configuration  
APP_PORT=8080
APP_ENV=production
JWT_SECRET=your-super-secret-jwt-key

# Optional
LOG_LEVEL=info
```

## 🐳 Docker

```bash
# Build and run with Docker Compose
docker-compose up --build

# Run in background
docker-compose up -d

# Stop containers
docker-compose down
```

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### **Development Guidelines**
- Follow Go conventions and best practices
- Add tests for new features
- Update documentation for API changes
- Use conventional commit messages
- Ensure all tests pass before submitting PR

## 📈 Roadmap

### **Planned Features**
- [ ] Inventory Management
- [ ] Warehouse Location Management  
- [ ] Stock Movement Tracking
- [ ] Order Management
- [ ] Supplier Management
- [ ] Reports and Analytics
- [ ] Barcode/QR Code Integration
- [ ] Mobile App Support
- [ ] Advanced Search and Filtering
- [ ] Export/Import Functionality

### **Technical Improvements**
- [ ] Comprehensive Unit Tests
- [ ] Integration Tests
- [ ] Performance Optimization
- [ ] Caching Layer (Redis)
- [ ] Message Queue Integration
- [ ] Monitoring and Metrics
- [ ] API Documentation (Swagger)
- [ ] Rate Limiting
- [ ] Request Validation Middleware

## 📝 Changelog

### **v1.0.0** (Current)
- ✅ Complete authentication system with JWT
- ✅ User and Role management
- ✅ Brand and Category management
- ✅ Product and Product Batch management
- ✅ Automatic change tracking for Product Batches
- ✅ Comprehensive audit trails
- ✅ RESTful API with proper error handling
- ✅ Database migration and seeding
- ✅ Hot reload development setup
- ✅ Docker containerization

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you have any questions or need help, please open an issue on GitHub.
