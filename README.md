# GO-WMS (Warehouse Management System)

A modern Warehouse Management System built with Go, Fiber, and PostgreSQL.

## 🚀 Features

- **JWT Authentication** - Secure user authentication
- **RESTful API** - Clean API design with versioning
- **Database Migration** - Automatic database setup
- **Modular Architecture** - Clean separation of concerns
- **Raw SQL Support** - Both GORM and raw SQL queries
- **Docker Support** - Containerized development

## 📋 Prerequisites

- Go 1.23+
- PostgreSQL 16+
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

### 5. Run Application
```bash
go run main.go
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

### Health Check
- `GET /health` - Global health check
- `GET /api/v1/health` - API v1 health check

## 🏗️ Project Structure

```
GO-WMS/
├── cmd/                    # Application entry points
├── config/                 # Configuration files
├── database/               # Database related files
│   ├── connection.go       # Database connection
│   ├── migration.go        # Database migrations
│   └── seeder/            # Database seeders
├── internal/              # Private application code
│   ├── handler/           # HTTP handlers
│   ├── middleware/        # HTTP middlewares
│   ├── model/             # Data models
│   ├── repository/        # Data access layer
│   ├── routes/            # Route definitions
│   ├── service/           # Business logic
│   └── utils/             # Utility functions
├── pkg/                   # Public library code
└── main.go               # Application entry point
```

## 🔧 Development

### Default Login Credentials
```
Email: alice@mail.com
Password: password123

Email: bob@mail.com
Password: password123
```

### Database Commands
```bash
# Start database
docker-compose up -d db

# View logs
docker-compose logs db

# Stop database
docker-compose down
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📦 Building

```bash
# Build for current platform
go build -o bin/go-wms main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o bin/go-wms-linux main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o bin/go-wms.exe main.go
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

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you have any questions or need help, please open an issue on GitHub.
