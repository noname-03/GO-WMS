# 🛠️ Development Guide
Complete development setup and guidelines for GO-WMS (Warehouse Management System).

## 🚀 Quick Start

### Prerequisites
- **Go**: Version 1.23.12 or higher
- **PostgreSQL**: Version 16 or higher
- **Git**: Latest version
- **Code Editor**: VS Code (recommended) or any Go-compatible IDE

### Installation Steps

1. **Clone Repository**
```bash
git clone https://github.com/noname-03/GO-WMS.git
cd GO-WMS
```

2. **Install Dependencies**
```bash
go mod tidy
```

3. **Setup Environment**
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. **Setup Database**
```bash
# Create PostgreSQL database
createdb go_wms_development

# Run the application (auto-migration will handle schema)
go run main.go
```

5. **Verify Installation**
```bash
# Application should be running on http://localhost:8080
curl http://localhost:8080/health
```

## 🏗️ Development Environment

### Recommended Tools

#### Code Editor Setup (VS Code)
**Required Extensions:**
- Go (golang.org)
- REST Client or Thunder Client
- GitLens
- Error Lens

**VS Code Settings (`settings.json`):**
```json
{
    "go.useLanguageServer": true,
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.vetOnSave": "package",
    "go.buildOnSave": "package",
    "go.testTimeout": "30s"
}
```

#### Database Tools
- **pgAdmin**: PostgreSQL administration
- **DBeaver**: Universal database tool
- **TablePlus**: Modern database client

#### API Testing
- **Postman**: API testing and documentation
- **Insomnia**: REST API client
- **Thunder Client**: VS Code extension

### Development Dependencies

#### Go Tools Installation
```bash
# Install useful Go tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/air-verse/air@latest

# Install debugging tools
go install github.com/go-delve/delve/cmd/dlv@latest
```

#### Air (Hot Reload)
```bash
# Install Air for hot reloading
go install github.com/air-verse/air@latest

# Create .air.toml configuration
air init

# Start development with hot reload
air
```

**Air Configuration (`.air.toml`):**
```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
```

## 📁 Project Structure Guidelines

### Directory Organization
```
internal/
├── handler/        # HTTP handlers (controllers)
│   ├── auth_handler.go
│   ├── user_handler.go
│   └── product_handler.go
├── service/        # Business logic layer
│   ├── auth_service.go
│   ├── user_service.go
│   └── product_service.go
├── repository/     # Data access layer
│   ├── user_repository.go
│   └── product_repository.go
├── model/          # Data models
│   ├── user.go
│   └── product.go
├── middleware/     # HTTP middleware
│   └── auth.go
├── routes/         # Route definitions
│   ├── routes.go
│   └── v1/
└── utils/          # Utility functions
    └── jwt.go
```

### Naming Conventions

#### Files and Directories
- **Snake case** for file names: `user_handler.go`
- **Lowercase** for directory names: `internal/handler/`
- **Descriptive names**: `product_batch_service.go`

#### Go Code Conventions
```go
// Package names: lowercase, single word
package handler

// Interface names: noun + "er" suffix
type UserService interface {}

// Struct names: PascalCase
type UserHandler struct {}

// Function names: PascalCase for public, camelCase for private
func CreateUser() {} // public
func validateUser() {} // private

// Constants: PascalCase or UPPER_CASE
const MaxRetries = 3
const DEFAULT_PAGE_SIZE = 10

// Variables: camelCase
var userRepository UserRepository
```

## 🔧 Development Workflow

### Git Workflow

#### Branch Strategy
```bash
# Main branches
main        # Production-ready code
develop     # Integration branch for features

# Feature branches
feature/user-authentication
feature/product-management
feature/batch-tracking

# Hotfix branches
hotfix/security-patch
hotfix/critical-bug-fix
```

#### Commit Guidelines
```bash
# Commit message format
<type>(<scope>): <subject>

# Types
feat:     # New feature
fix:      # Bug fix
docs:     # Documentation changes
style:    # Code style changes (formatting, etc.)
refactor: # Code refactoring
test:     # Adding or updating tests
chore:    # Maintenance tasks

# Examples
feat(auth): implement JWT authentication
fix(user): resolve email validation issue
docs(api): update endpoint documentation
refactor(service): extract user validation logic
```

#### Git Hooks Setup
```bash
# Install pre-commit hooks
go install github.com/dnephin/pre-commit-golang@latest

# Create .pre-commit-config.yaml
```

**Pre-commit Configuration:**
```yaml
repos:
  - repo: https://github.com/dnephin/pre-commit-golang
    rev: master
    hooks:
      - id: go-fmt
      - id: go-vet-mod
      - id: go-mod-tidy
      - id: golangci-lint
```

### Code Quality Standards

#### Linting Configuration
**Create `.golangci.yml`:**
```yaml
run:
  deadline: 5m
  tests: true

linters-settings:
  govet:
    check-shadowing: true
  golint:
    min-confidence: 0.8
  gocyclo:
    min-complexity: 15
  maligned:
    suggest-new: true
  dupl:
    threshold: 100
  goconst:
    min-len: 2
    min-occurrences: 2

linters:
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - dupl
    - errcheck
    - funlen
    - gochecknoinits
    - goconst
    - gocritic
    - gocyclo
    - gofmt
    - goimports
    - golint
    - gomnd
    - goprintffuncname
    - gosec
    - gosimple
    - govet
    - ineffassign
    - interfacer
    - lll
    - misspell
    - nakedret
    - rowserrcheck
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck
    - whitespace
```

#### Code Formatting
```bash
# Format code
go fmt ./...

# Organize imports
goimports -w .

# Run linter
golangci-lint run

# Fix auto-fixable issues
golangci-lint run --fix
```

### Testing Strategy

#### Test File Organization
```
internal/
├── handler/
│   ├── user_handler.go
│   └── user_handler_test.go
├── service/
│   ├── user_service.go
│   └── user_service_test.go
└── repository/
    ├── user_repository.go
    └── user_repository_test.go
```

#### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -run TestCreateUser ./internal/handler

# Run tests in verbose mode
go test -v ./...

# Run tests with race detection
go test -race ./...
```

#### Test Writing Guidelines
```go
// Example test structure
func TestUserService_CreateUser(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    defer teardownTestDB(t, db)
    
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    
    // Test cases
    tests := []struct {
        name    string
        user    *model.User
        wantErr bool
    }{
        {
            name: "valid user creation",
            user: &model.User{
                Email:    "test@example.com",
                Password: "password123",
            },
            wantErr: false,
        },
        {
            name: "duplicate email",
            user: &model.User{
                Email:    "existing@example.com",
                Password: "password123",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := userService.CreateUser(tt.user)
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## 🗃️ Database Development

### Database Setup for Development
```bash
# Create development database
createdb go_wms_development

# Create test database
createdb go_wms_test

# Setup environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=yourpassword
export DB_NAME=go_wms_development
```

### Migration Development
```go
// Add new migration in database/migration.go
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &model.User{},
        &model.Role{},
        &model.Category{},
        &model.Brand{},
        &model.Product{},
        &model.ProductBatch{},
        &model.ProductBatchTrack{},
    )
}
```

### Seeder Development
```go
// Create new seeder in database/seeder/
type NewEntitySeeder struct {
    db *gorm.DB
}

func (s *NewEntitySeeder) Seed() error {
    // Seeding logic
    return nil
}

// Register in database/seeder/registry.go
func GetSeeders(db *gorm.DB) []seeder.Seeder {
    return []seeder.Seeder{
        &NewEntitySeeder{db: db},
        // ... other seeders
    }
}
```

## 🌐 API Development

### Handler Development Pattern
```go
// Handler structure
type EntityHandler struct {
    entityService service.EntityService
}

// Constructor
func NewEntityHandler(entityService service.EntityService) *EntityHandler {
    return &EntityHandler{
        entityService: entityService,
    }
}

// CRUD methods
func (h *EntityHandler) GetAll(c *fiber.Ctx) error {
    entities, err := h.entityService.GetAll()
    if err != nil {
        return c.Status(500).JSON(helper.ErrorResponse(err.Error()))
    }
    return c.JSON(helper.SuccessResponse("Entities retrieved", entities))
}

func (h *EntityHandler) Create(c *fiber.Ctx) error {
    var entity model.Entity
    if err := c.BodyParser(&entity); err != nil {
        return c.Status(400).JSON(helper.ErrorResponse("Invalid request body"))
    }
    
    if err := h.entityService.Create(&entity); err != nil {
        return c.Status(500).JSON(helper.ErrorResponse(err.Error()))
    }
    
    return c.Status(201).JSON(helper.SuccessResponse("Entity created", entity))
}
```

### Service Development Pattern
```go
// Service interface
type EntityService interface {
    GetAll() ([]*model.Entity, error)
    GetByID(id uint) (*model.Entity, error)
    Create(entity *model.Entity) error
    Update(entity *model.Entity) error
    Delete(id uint) error
}

// Service implementation
type entityService struct {
    entityRepo repository.EntityRepository
}

func NewEntityService(entityRepo repository.EntityRepository) EntityService {
    return &entityService{
        entityRepo: entityRepo,
    }
}

func (s *entityService) Create(entity *model.Entity) error {
    // Validation logic
    if entity.Name == "" {
        return errors.New("name is required")
    }
    
    // Business logic
    entity.CreatedAt = time.Now()
    
    // Repository call
    return s.entityRepo.Create(entity)
}
```

### Repository Development Pattern
```go
// Repository interface
type EntityRepository interface {
    GetAll() ([]*model.Entity, error)
    GetByID(id uint) (*model.Entity, error)
    Create(entity *model.Entity) error
    Update(entity *model.Entity) error
    Delete(id uint) error
}

// Repository implementation
type entityRepository struct {
    db *gorm.DB
}

func NewEntityRepository(db *gorm.DB) EntityRepository {
    return &entityRepository{db: db}
}

func (r *entityRepository) Create(entity *model.Entity) error {
    return r.db.Create(entity).Error
}

func (r *entityRepository) GetByID(id uint) (*model.Entity, error) {
    var entity model.Entity
    err := r.db.First(&entity, id).Error
    if err != nil {
        return nil, err
    }
    return &entity, nil
}
```

## 🔧 Debugging

### Debug Configuration (VS Code)
**Create `.vscode/launch.json`:**
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/main.go",
            "env": {
                "DB_HOST": "localhost",
                "DB_PORT": "5432",
                "DB_USER": "postgres",
                "DB_PASSWORD": "yourpassword",
                "DB_NAME": "go_wms_development"
            },
            "args": []
        }
    ]
}
```

### Debugging Techniques
```go
// Log debugging
import "log"

func (s *userService) CreateUser(user *model.User) error {
    log.Printf("Creating user: %+v", user)
    
    // ... business logic
    
    log.Printf("User created successfully: ID=%d", user.ID)
    return nil
}

// Delve debugger usage
// Set breakpoint and run:
dlv debug main.go

// Remote debugging
dlv debug --headless --listen=:2345 --api-version=2 main.go
```

### Performance Profiling
```go
// Add profiling endpoints
import _ "net/http/pprof"

// In main.go
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

// Profile memory
go tool pprof http://localhost:6060/debug/pprof/heap

// Profile CPU
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

## 📝 Documentation Guidelines

### Code Documentation
```go
// Package documentation
// Package handler provides HTTP request handlers for the GO-WMS API.
// It implements the controller layer of the MVC architecture.
package handler

// Function documentation
// CreateUser handles HTTP POST requests to create a new user.
// It validates the request body, calls the user service, and returns
// the appropriate HTTP response.
//
// Parameters:
//   - c: Fiber context containing the HTTP request and response
//
// Returns:
//   - error: nil on success, error on failure
//
// Example:
//   POST /api/v1/users
//   {
//     "email": "user@example.com",
//     "password": "secure123"
//   }
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    // Implementation
}
```

### API Documentation
```bash
# Generate swagger documentation
go install github.com/swaggo/swag/cmd/swag@latest

# Add swagger comments to handlers
// @Summary Create a new user
// @Description Create a new user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.User true "User data"
// @Success 201 {object} helper.Response
// @Failure 400 {object} helper.Response
// @Router /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {}

# Generate docs
swag init
```

## 🚀 Performance Optimization

### Database Query Optimization
```go
// Use preloading to avoid N+1 queries
func (r *productRepository) GetAllWithBatches() ([]*model.Product, error) {
    var products []*model.Product
    err := r.db.Preload("ProductBatches").Find(&products).Error
    return products, err
}

// Use specific fields selection
func (r *userRepository) GetUserList() ([]*model.User, error) {
    var users []*model.User
    err := r.db.Select("id", "email", "created_at").Find(&users).Error
    return users, err
}

// Use pagination
func (r *productRepository) GetPaginated(page, limit int) ([]*model.Product, error) {
    var products []*model.Product
    offset := (page - 1) * limit
    err := r.db.Offset(offset).Limit(limit).Find(&products).Error
    return products, err
}
```

### Memory Optimization
```go
// Pool expensive objects
var userPool = sync.Pool{
    New: func() interface{} {
        return &model.User{}
    },
}

func (s *userService) ProcessUser() {
    user := userPool.Get().(*model.User)
    defer userPool.Put(user)
    
    // Use user object
}
```

## 📋 Development Checklist

### Before Committing
- [ ] Code compiles without errors
- [ ] All tests pass (`go test ./...`)
- [ ] Code is formatted (`go fmt ./...`)
- [ ] Imports are organized (`goimports -w .`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] Documentation is updated
- [ ] Commit message follows convention

### Before Pull Request
- [ ] Feature branch is up to date with develop
- [ ] All acceptance criteria are met
- [ ] Integration tests pass
- [ ] API documentation is updated
- [ ] Database migrations are tested
- [ ] Performance impact is evaluated
- [ ] Security implications are reviewed

### Before Release
- [ ] All tests pass in CI/CD
- [ ] Performance benchmarks meet requirements
- [ ] Security scan passes
- [ ] Documentation is complete
- [ ] Deployment scripts are tested
- [ ] Rollback plan is prepared