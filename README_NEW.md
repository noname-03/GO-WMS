# 🏭 GO-WMS (Warehouse Management System)

![Go Version](https://img.shields.io/badge/Go-1.23.12-blue.svg)
![Fiber Version](https://img.shields.io/badge/Fiber-v2.52.9-green.svg)
![GORM Version](https://img.shields.io/badge/GORM-Latest-orange.svg)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-blue.svg)
![License](https://img.shields.io/badge/License-MIT-yellow.svg)

A comprehensive Warehouse Management System built with Go, featuring modern web technologies and robust architecture for efficient inventory management with complete audit trails.

## ✨ Key Features

- 🔐 **JWT Authentication** - Secure user authentication system
- 👥 **Role-Based Access Control** - Multi-level permissions (Admin, Manager, Operator, Viewer)
- 📦 **Product Management** - Complete product lifecycle with categories and brands
- 📋 **Advanced Batch Tracking** - Real-time inventory tracking with automatic audit trails
- 🔍 **Comprehensive Logging** - Detailed tracking of all product batch operations
- 🚀 **RESTful API** - Clean and intuitive API design
- 🐳 **Docker Ready** - Containerized deployment support
- ⚡ **Hot Reload** - Development-friendly with Air

## 🏗️ Architecture

### Technology Stack
- **Backend**: Go 1.23.12 + Fiber v2.52.9 + GORM
- **Database**: PostgreSQL 16
- **Authentication**: JWT + bcrypt
- **Development**: Air (hot reload)
- **Deployment**: Docker + Docker Compose

### Project Structure
```
GO-WMS/
├── main.go                 # Application entry point
├── database/              # Database configuration & migrations
├── internal/              # Private application code
│   ├── handler/          # HTTP request handlers
│   ├── service/          # Business logic layer  
│   ├── repository/       # Data access layer
│   ├── model/           # Data models
│   ├── middleware/      # HTTP middleware
│   ├── routes/          # Route definitions
│   └── utils/           # Utility functions
├── pkg/helper/           # Helper functions
└── docs/                # Comprehensive documentation
```

## 🚀 Quick Start

### Prerequisites
- Go 1.23.12+
- PostgreSQL 16+
- Git

### Installation

1. **Clone & Setup**
   ```bash
   git clone https://github.com/noname-03/GO-WMS.git
   cd GO-WMS
   go mod tidy
   ```

2. **Configure Environment**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Create Database**
   ```bash
   createdb go_wms
   ```

4. **Run Application**
   ```bash
   # Development with hot reload
   air
   
   # Or run directly
   go run main.go
   ```

5. **Access Application**
   - API Base: `http://localhost:8080/api/v1`
   - Health Check: `http://localhost:8080/health`

### Docker Deployment
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

### Process Management Commands
```bash
# Check running processes
ps aux | grep main

# Kill specific process
kill -9 547218

# Start in background
nohup ./go-wms &

# Monitor logs
tail -f nohup.out
```

## 📚 Documentation

This project includes comprehensive documentation organized into specialized sections:

### 📖 **[API Documentation](docs/API.md)**
Complete API reference with all endpoints, request/response formats, and authentication details.

### 🏗️ **[Architecture Guide](docs/ARCHITECTURE.md)**  
Detailed system architecture, design patterns, and technical implementation details.

### 🗃️ **[Database Documentation](docs/DATABASE.md)**
Database schema, relationships, queries, and optimization strategies.

### 🛠️ **[Development Guide](docs/DEVELOPMENT.md)**
Development setup, coding standards, workflow, and best practices.

### 🚀 **[Deployment Guide](docs/DEPLOYMENT.md)**
Production deployment, hosting commands, monitoring, and maintenance procedures.

### 🧪 **[Testing Guide](docs/TESTING.md)**
Testing strategies, examples, coverage requirements, and CI/CD pipelines.

## 🔧 Core Features

### Authentication & Authorization
- JWT-based secure authentication
- Role-based access control (Admin, Manager, Operator, Viewer)
- Protected routes with middleware

### Product Management
- Complete CRUD operations for products
- Category and brand organization
- Advanced search and filtering

### Product Batch Tracking
- **Automatic Audit Trail**: Every operation is automatically tracked
- **Change Detection**: System detects and logs all modifications
- **User Attribution**: All actions are linked to authenticated users
- **Detailed Descriptions**: Comprehensive change descriptions

**Example Tracking:**
```json
{
  "action": "UPDATE",
  "description": "Quantity changed from 100 to 150; Status changed from 'pending' to 'active'",
  "user_id": 1,
  "created_at": "2024-03-15T11:00:00Z"
}
```

### Data Management
- Automatic database migrations
- Comprehensive data seeding
- Connection pooling and optimization

## 🔌 API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration

### Core Resources
- **Users**: `/api/v1/users` - User management
- **Products**: `/api/v1/products` - Product operations
- **Product Batches**: `/api/v1/product-batches` - Batch management with tracking
- **Categories**: `/api/v1/categories` - Category management
- **Brands**: `/api/v1/brands` - Brand management
- **Roles**: `/api/v1/roles` - Role management

### Tracking
- `GET /api/v1/product-batches/:id/tracking` - Get complete tracking history

For detailed API documentation with examples, see **[API Documentation](docs/API.md)**.

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific tests
go test -run TestUserHandler ./internal/handler
```

For comprehensive testing strategies and examples, see **[Testing Guide](docs/TESTING.md)**.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Follow the development guidelines in **[Development Guide](docs/DEVELOPMENT.md)**
4. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

- 📖 **Documentation**: Comprehensive guides in the `docs/` directory
- 🐛 **Issues**: Report bugs via GitHub Issues
- 💡 **Features**: Request features through GitHub Issues

## 🎯 What's Next?

- Advanced reporting and analytics
- Real-time notifications
- Multi-warehouse support
- Mobile API optimization
- Enhanced security features

---

**Built with ❤️ using Go, Fiber, GORM, and PostgreSQL**

> For detailed information about any aspect of the system, please refer to the appropriate documentation file in the `docs/` directory.
