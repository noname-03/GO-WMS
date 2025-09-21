# 🧪 Testing Guide

Comprehensive testing strategies, examples, and best practices for GO-WMS.

## 🎯 Testing Strategy Overview

### Testing Pyramid
```
                    ▲
                   /|\
                  / | \
                 /  |  \
                /   |   \
               /    |    \
              /     |     \
             /  E2E |      \
            /  Tests|       \
           /_______|________\
          /        |         \
         /Integration|         \
        /   Tests    |          \
       /____________|___________\
      /             |            \
     /    Unit      |    Unit     \
    /     Tests     |    Tests     \
   /________________|_________________\
```

#### Test Distribution
- **Unit Tests (70%)**: Fast, isolated, testing individual components
- **Integration Tests (20%)**: Testing component interactions
- **End-to-End Tests (10%)**: Full application workflow testing

### Testing Principles
1. **Fast Feedback**: Unit tests should run in milliseconds
2. **Isolation**: Tests should not depend on each other
3. **Deterministic**: Same input should always produce same output
4. **Readable**: Tests should be self-documenting
5. **Maintainable**: Easy to update when code changes

## 🏗️ Test Environment Setup

### Test Dependencies
```bash
# Install testing tools
go install github.com/stretchr/testify@latest
go install github.com/golang/mock/gomock@latest
go install github.com/golang/mock/mockgen@latest
go install github.com/DATA-DOG/go-sqlmock@latest
```

### Test Configuration
```go
// test/config.go
package test

import (
    "os"
    "testing"
    
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type TestConfig struct {
    DB *gorm.DB
}

func SetupTestDB(t *testing.T) *gorm.DB {
    // Use in-memory SQLite for fast tests
    dsn := "file::memory:?cache=shared"
    db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    
    // Auto-migrate for tests
    err = db.AutoMigrate(
        &model.User{},
        &model.Role{},
        &model.Category{},
        &model.Brand{},
        &model.Product{},
        &model.ProductBatch{},
        &model.ProductBatchTrack{},
    )
    
    if err != nil {
        t.Fatalf("Failed to migrate test database: %v", err)
    }
    
    return db
}

func TeardownTestDB(t *testing.T, db *gorm.DB) {
    sqlDB, err := db.DB()
    if err != nil {
        t.Errorf("Failed to get underlying sql.DB: %v", err)
        return
    }
    
    if err := sqlDB.Close(); err != nil {
        t.Errorf("Failed to close test database: %v", err)
    }
}
```

### Test Utilities
```go
// test/helpers.go
package test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
)

// CreateTestRequest creates a test HTTP request
func CreateTestRequest(t *testing.T, method, url string, body interface{}) *http.Request {
    var buf bytes.Buffer
    if body != nil {
        err := json.NewEncoder(&buf).Encode(body)
        assert.NoError(t, err)
    }
    
    req := httptest.NewRequest(method, url, &buf)
    req.Header.Set("Content-Type", "application/json")
    return req
}

// ParseResponse parses a test response body
func ParseResponse(t *testing.T, resp *http.Response, target interface{}) {
    defer resp.Body.Close()
    err := json.NewDecoder(resp.Body).Decode(target)
    assert.NoError(t, err)
}

// CreateTestApp creates a test Fiber app
func CreateTestApp() *fiber.App {
    app := fiber.New(fiber.Config{
        Testing: true,
    })
    return app
}
```

## 🔧 Unit Testing

### Repository Testing
```go
// internal/repository/user_repository_test.go
package repository

import (
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "your-project/internal/model"
    "your-project/test"
)

type UserRepositoryTestSuite struct {
    suite.Suite
    db   *gorm.DB
    repo UserRepository
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
    suite.db = test.SetupTestDB(suite.T())
    suite.repo = NewUserRepository(suite.db)
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
    test.TeardownTestDB(suite.T(), suite.db)
}

func (suite *UserRepositoryTestSuite) SetupTest() {
    // Clean up before each test
    suite.db.Exec("DELETE FROM users")
}

func (suite *UserRepositoryTestSuite) TestCreate() {
    // Arrange
    user := &model.User{
        Email:    "test@example.com",
        Password: "hashedpassword",
    }
    
    // Act
    err := suite.repo.Create(user)
    
    // Assert
    assert.NoError(suite.T(), err)
    assert.NotZero(suite.T(), user.ID)
    assert.NotZero(suite.T(), user.CreatedAt)
}

func (suite *UserRepositoryTestSuite) TestGetByID() {
    // Arrange
    user := &model.User{
        Email:    "test@example.com",
        Password: "hashedpassword",
    }
    suite.repo.Create(user)
    
    // Act
    found, err := suite.repo.GetByID(user.ID)
    
    // Assert
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), user.Email, found.Email)
    assert.Equal(suite.T(), user.ID, found.ID)
}

func (suite *UserRepositoryTestSuite) TestGetByID_NotFound() {
    // Act
    found, err := suite.repo.GetByID(999)
    
    // Assert
    assert.Error(suite.T(), err)
    assert.Nil(suite.T(), found)
    assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)
}

func (suite *UserRepositoryTestSuite) TestUpdate() {
    // Arrange
    user := &model.User{
        Email:    "test@example.com",
        Password: "hashedpassword",
    }
    suite.repo.Create(user)
    
    // Act
    user.Email = "updated@example.com"
    err := suite.repo.Update(user)
    
    // Assert
    assert.NoError(suite.T(), err)
    
    // Verify update
    found, _ := suite.repo.GetByID(user.ID)
    assert.Equal(suite.T(), "updated@example.com", found.Email)
}

func (suite *UserRepositoryTestSuite) TestDelete() {
    // Arrange
    user := &model.User{
        Email:    "test@example.com",
        Password: "hashedpassword",
    }
    suite.repo.Create(user)
    
    // Act
    err := suite.repo.Delete(user.ID)
    
    // Assert
    assert.NoError(suite.T(), err)
    
    // Verify deletion
    found, err := suite.repo.GetByID(user.ID)
    assert.Error(suite.T(), err)
    assert.Nil(suite.T(), found)
}

// Run the test suite
func TestUserRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(UserRepositoryTestSuite))
}
```

### Service Testing with Mocks
```go
// internal/service/user_service_test.go
package service

import (
    "errors"
    "testing"
    
    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"
    "golang.org/x/crypto/bcrypt"
    "your-project/internal/model"
    "your-project/mocks"
)

func TestUserService_CreateUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockUserRepository(ctrl)
    service := NewUserService(mockRepo)
    
    tests := []struct {
        name    string
        user    *model.User
        setup   func()
        wantErr bool
        errMsg  string
    }{
        {
            name: "successful creation",
            user: &model.User{
                Email:    "test@example.com",
                Password: "password123",
            },
            setup: func() {
                mockRepo.EXPECT().
                    GetByEmail("test@example.com").
                    Return(nil, gorm.ErrRecordNotFound)
                mockRepo.EXPECT().
                    Create(gomock.Any()).
                    Return(nil)
            },
            wantErr: false,
        },
        {
            name: "duplicate email",
            user: &model.User{
                Email:    "existing@example.com",
                Password: "password123",
            },
            setup: func() {
                existingUser := &model.User{
                    ID:    1,
                    Email: "existing@example.com",
                }
                mockRepo.EXPECT().
                    GetByEmail("existing@example.com").
                    Return(existingUser, nil)
            },
            wantErr: true,
            errMsg:  "email already exists",
        },
        {
            name: "invalid email",
            user: &model.User{
                Email:    "invalid-email",
                Password: "password123",
            },
            setup:   func() {},
            wantErr: true,
            errMsg:  "invalid email format",
        },
        {
            name: "weak password",
            user: &model.User{
                Email:    "test@example.com",
                Password: "123",
            },
            setup:   func() {},
            wantErr: true,
            errMsg:  "password must be at least 6 characters",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setup()
            
            err := service.CreateUser(tt.user)
            
            if tt.wantErr {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                assert.NoError(t, err)
                // Verify password was hashed
                assert.NotEqual(t, "password123", tt.user.Password)
                err = bcrypt.CompareHashAndPassword([]byte(tt.user.Password), []byte("password123"))
                assert.NoError(t, err)
            }
        })
    }
}

func TestUserService_GetUserByID(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mocks.NewMockUserRepository(ctrl)
    service := NewUserService(mockRepo)
    
    expectedUser := &model.User{
        ID:    1,
        Email: "test@example.com",
    }
    
    mockRepo.EXPECT().
        GetByID(uint(1)).
        Return(expectedUser, nil)
    
    user, err := service.GetUserByID(1)
    
    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
}
```

### Handler Testing
```go
// internal/handler/user_handler_test.go
package handler

import (
    "encoding/json"
    "net/http"
    "testing"
    
    "github.com/golang/mock/gomock"
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "your-project/internal/model"
    "your-project/mocks"
    "your-project/pkg/helper"
    "your-project/test"
)

func TestUserHandler_CreateUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockService := mocks.NewMockUserService(ctrl)
    handler := NewUserHandler(mockService)
    
    app := test.CreateTestApp()
    app.Post("/users", handler.CreateUser)
    
    tests := []struct {
        name           string
        requestBody    interface{}
        setup          func()
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name: "successful creation",
            requestBody: map[string]interface{}{
                "email":    "test@example.com",
                "password": "password123",
            },
            setup: func() {
                mockService.EXPECT().
                    CreateUser(gomock.Any()).
                    Return(nil)
            },
            expectedStatus: 201,
            expectedBody: map[string]interface{}{
                "success": true,
                "message": "User created successfully",
            },
        },
        {
            name: "invalid request body",
            requestBody: map[string]interface{}{
                "invalid": "data",
            },
            setup:          func() {},
            expectedStatus: 400,
            expectedBody: map[string]interface{}{
                "success": false,
                "error":   "Invalid request body",
            },
        },
        {
            name: "service error",
            requestBody: map[string]interface{}{
                "email":    "test@example.com",
                "password": "password123",
            },
            setup: func() {
                mockService.EXPECT().
                    CreateUser(gomock.Any()).
                    Return(errors.New("email already exists"))
            },
            expectedStatus: 500,
            expectedBody: map[string]interface{}{
                "success": false,
                "error":   "email already exists",
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setup()
            
            req := test.CreateTestRequest(t, "POST", "/users", tt.requestBody)
            resp, err := app.Test(req)
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expectedStatus, resp.StatusCode)
            
            var responseBody map[string]interface{}
            test.ParseResponse(t, resp, &responseBody)
            
            for key, expectedValue := range tt.expectedBody {
                assert.Equal(t, expectedValue, responseBody[key])
            }
        })
    }
}

func TestUserHandler_GetUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockService := mocks.NewMockUserService(ctrl)
    handler := NewUserHandler(mockService)
    
    app := test.CreateTestApp()
    app.Get("/users/:id", handler.GetUser)
    
    expectedUser := &model.User{
        ID:    1,
        Email: "test@example.com",
    }
    
    mockService.EXPECT().
        GetUserByID(uint(1)).
        Return(expectedUser, nil)
    
    req := test.CreateTestRequest(t, "GET", "/users/1", nil)
    resp, err := app.Test(req)
    
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    var responseBody helper.Response
    test.ParseResponse(t, resp, &responseBody)
    
    assert.True(t, responseBody.Success)
    assert.Equal(t, "User retrieved successfully", responseBody.Message)
    
    // Parse user data
    userData, _ := json.Marshal(responseBody.Data)
    var user model.User
    json.Unmarshal(userData, &user)
    assert.Equal(t, expectedUser.ID, user.ID)
    assert.Equal(t, expectedUser.Email, user.Email)
}
```

## 🔗 Integration Testing

### Database Integration Tests
```go
// test/integration/database_test.go
package integration

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "your-project/database"
    "your-project/internal/model"
    "your-project/internal/repository"
    "your-project/internal/service"
)

type DatabaseIntegrationTestSuite struct {
    suite.Suite
    db          *gorm.DB
    userRepo    repository.UserRepository
    userService service.UserService
}

func (suite *DatabaseIntegrationTestSuite) SetupSuite() {
    // Use test database
    config := &database.Config{
        Host:     "localhost",
        Port:     "5432",
        User:     "postgres",
        Password: "password",
        DBName:   "go_wms_test",
        SSLMode:  "disable",
    }
    
    db, err := database.Connect(config)
    suite.Require().NoError(err)
    
    // Run migrations
    err = database.AutoMigrate(db)
    suite.Require().NoError(err)
    
    suite.db = db
    suite.userRepo = repository.NewUserRepository(db)
    suite.userService = service.NewUserService(suite.userRepo)
}

func (suite *DatabaseIntegrationTestSuite) TearDownSuite() {
    sqlDB, _ := suite.db.DB()
    sqlDB.Close()
}

func (suite *DatabaseIntegrationTestSuite) SetupTest() {
    // Clean up before each test
    suite.db.Exec("TRUNCATE users CASCADE")
}

func (suite *DatabaseIntegrationTestSuite) TestUserCreationFlow() {
    // Test the complete flow from service to database
    user := &model.User{
        Email:    "integration@test.com",
        Password: "password123",
    }
    
    // Create user through service
    err := suite.userService.CreateUser(user)
    assert.NoError(suite.T(), err)
    
    // Verify user exists in database
    found, err := suite.userRepo.GetByID(user.ID)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), user.Email, found.Email)
    assert.NotEqual(suite.T(), "password123", found.Password) // Should be hashed
}

func (suite *DatabaseIntegrationTestSuite) TestUserRoleAssignment() {
    // Create role
    role := &model.Role{
        Name:        "test_role",
        Description: "Test role",
    }
    suite.db.Create(role)
    
    // Create user
    user := &model.User{
        Email:    "user@test.com",
        Password: "password123",
    }
    err := suite.userService.CreateUser(user)
    assert.NoError(suite.T(), err)
    
    // Assign role to user
    err = suite.db.Model(user).Association("Roles").Append(role)
    assert.NoError(suite.T(), err)
    
    // Verify relationship
    var foundUser model.User
    err = suite.db.Preload("Roles").First(&foundUser, user.ID).Error
    assert.NoError(suite.T(), err)
    assert.Len(suite.T(), foundUser.Roles, 1)
    assert.Equal(suite.T(), role.Name, foundUser.Roles[0].Name)
}

func TestDatabaseIntegrationTestSuite(t *testing.T) {
    suite.Run(t, new(DatabaseIntegrationTestSuite))
}
```

### API Integration Tests
```go
// test/integration/api_test.go
package integration

import (
    "encoding/json"
    "testing"
    
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "your-project/internal/routes"
    "your-project/test"
)

type APIIntegrationTestSuite struct {
    suite.Suite
    app *fiber.App
    db  *gorm.DB
}

func (suite *APIIntegrationTestSuite) SetupSuite() {
    suite.db = test.SetupTestDB(suite.T())
    suite.app = fiber.New(fiber.Config{Testing: true})
    
    // Setup routes
    routes.SetupRoutes(suite.app, suite.db)
}

func (suite *APIIntegrationTestSuite) TearDownSuite() {
    test.TeardownTestDB(suite.T(), suite.db)
}

func (suite *APIIntegrationTestSuite) SetupTest() {
    // Clean database before each test
    suite.db.Exec("DELETE FROM users")
    suite.db.Exec("DELETE FROM roles")
}

func (suite *APIIntegrationTestSuite) TestUserRegistrationAndLogin() {
    // Register user
    registerData := map[string]interface{}{
        "email":    "test@example.com",
        "password": "password123",
    }
    
    req := test.CreateTestRequest(suite.T(), "POST", "/api/v1/auth/register", registerData)
    resp, err := suite.app.Test(req)
    
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), 201, resp.StatusCode)
    
    // Login user
    loginData := map[string]interface{}{
        "email":    "test@example.com",
        "password": "password123",
    }
    
    req = test.CreateTestRequest(suite.T(), "POST", "/api/v1/auth/login", loginData)
    resp, err = suite.app.Test(req)
    
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), 200, resp.StatusCode)
    
    var loginResponse map[string]interface{}
    test.ParseResponse(suite.T(), resp, &loginResponse)
    
    assert.True(suite.T(), loginResponse["success"].(bool))
    assert.Contains(suite.T(), loginResponse, "token")
}

func (suite *APIIntegrationTestSuite) TestProductBatchTracking() {
    // Create category, brand, product, and batch
    category := &model.Category{Name: "Test Category"}
    suite.db.Create(category)
    
    brand := &model.Brand{Name: "Test Brand"}
    suite.db.Create(brand)
    
    product := &model.Product{
        Name:       "Test Product",
        CategoryID: &category.ID,
        BrandID:    &brand.ID,
    }
    suite.db.Create(product)
    
    // Create product batch
    batchData := map[string]interface{}{
        "product_id":   product.ID,
        "batch_number": "BATCH001",
        "quantity":     100,
    }
    
    req := test.CreateTestRequest(suite.T(), "POST", "/api/v1/product-batches", batchData)
    resp, err := suite.app.Test(req)
    
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), 201, resp.StatusCode)
    
    // Verify tracking record was created
    var trackCount int64
    suite.db.Model(&model.ProductBatchTrack{}).Count(&trackCount)
    assert.Equal(suite.T(), int64(1), trackCount)
}

func TestAPIIntegrationTestSuite(t *testing.T) {
    suite.Run(t, new(APIIntegrationTestSuite))
}
```

## 🎭 End-to-End Testing

### E2E Test Setup
```go
// test/e2e/setup.go
package e2e

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "testing"
    "time"
    
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
)

type TestEnvironment struct {
    PostgresContainer testcontainers.Container
    AppURL           string
    DBConnectionString string
}

func SetupE2EEnvironment(t *testing.T) *TestEnvironment {
    ctx := context.Background()
    
    // Start PostgreSQL container
    postgresReq := testcontainers.ContainerRequest{
        Image:        "postgres:16",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_DB":       "go_wms_test",
            "POSTGRES_USER":     "postgres",
            "POSTGRES_PASSWORD": "testpass",
        },
        WaitingFor: wait.ForListeningPort("5432/tcp"),
    }
    
    postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: postgresReq,
        Started:          true,
    })
    
    if err != nil {
        t.Fatalf("Failed to start PostgreSQL container: %v", err)
    }
    
    // Get container connection details
    host, err := postgresContainer.Host(ctx)
    if err != nil {
        t.Fatalf("Failed to get container host: %v", err)
    }
    
    port, err := postgresContainer.MappedPort(ctx, "5432")
    if err != nil {
        t.Fatalf("Failed to get container port: %v", err)
    }
    
    dbConnectionString := fmt.Sprintf("host=%s port=%s user=postgres password=testpass dbname=go_wms_test sslmode=disable", host, port.Port())
    
    // Start application
    appURL := startApplication(t, dbConnectionString)
    
    return &TestEnvironment{
        PostgresContainer:  postgresContainer,
        AppURL:            appURL,
        DBConnectionString: dbConnectionString,
    }
}

func (env *TestEnvironment) Cleanup(t *testing.T) {
    ctx := context.Background()
    if err := env.PostgresContainer.Terminate(ctx); err != nil {
        t.Errorf("Failed to terminate container: %v", err)
    }
}

func startApplication(t *testing.T, dbConnectionString string) string {
    // Set environment variables
    os.Setenv("DB_CONNECTION_STRING", dbConnectionString)
    os.Setenv("APP_PORT", "8081")
    
    // Start application in goroutine
    go func() {
        // Start your application here
        // main.main()
    }()
    
    // Wait for application to start
    appURL := "http://localhost:8081"
    for i := 0; i < 30; i++ {
        resp, err := http.Get(appURL + "/health")
        if err == nil && resp.StatusCode == 200 {
            resp.Body.Close()
            return appURL
        }
        if resp != nil {
            resp.Body.Close()
        }
        time.Sleep(time.Second)
    }
    
    t.Fatal("Application failed to start within 30 seconds")
    return ""
}
```

### E2E Test Cases
```go
// test/e2e/warehouse_flow_test.go
package e2e

import (
    "encoding/json"
    "net/http"
    "strings"
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type WarehouseFlowTestSuite struct {
    suite.Suite
    env *TestEnvironment
}

func (suite *WarehouseFlowTestSuite) SetupSuite() {
    suite.env = SetupE2EEnvironment(suite.T())
}

func (suite *WarehouseFlowTestSuite) TearDownSuite() {
    suite.env.Cleanup(suite.T())
}

func (suite *WarehouseFlowTestSuite) TestCompleteWarehouseFlow() {
    baseURL := suite.env.AppURL + "/api/v1"
    
    // Step 1: Register user
    registerData := `{"email":"warehouse@test.com","password":"password123"}`
    resp, err := http.Post(baseURL+"/auth/register", "application/json", strings.NewReader(registerData))
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), 201, resp.StatusCode)
    resp.Body.Close()
    
    // Step 2: Login and get token
    loginData := `{"email":"warehouse@test.com","password":"password123"}`
    resp, err = http.Post(baseURL+"/auth/login", "application/json", strings.NewReader(loginData))
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), 200, resp.StatusCode)
    
    var loginResponse map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&loginResponse)
    resp.Body.Close()
    
    token := loginResponse["token"].(string)
    assert.NotEmpty(suite.T(), token)
    
    // Step 3: Create category
    categoryData := `{"name":"Electronics","description":"Electronic products"}`
    req, _ := http.NewRequest("POST", baseURL+"/categories", strings.NewReader(categoryData))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    resp, err = client.Do(req)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), 201, resp.StatusCode)
    
    var categoryResponse map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&categoryResponse)
    resp.Body.Close()
    
    categoryID := categoryResponse["data"].(map[string]interface{})["id"]
    
    // Step 4: Create brand
    brandData := `{"name":"TechBrand","description":"Technology brand"}`
    req, _ = http.NewRequest("POST", baseURL+"/brands", strings.NewReader(brandData))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err = client.Do(req)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), 201, resp.StatusCode)
    
    var brandResponse map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&brandResponse)
    resp.Body.Close()
    
    brandID := brandResponse["data"].(map[string]interface{})["id"]
    
    // Step 5: Create product
    productData := fmt.Sprintf(`{"name":"Laptop","category_id":%v,"brand_id":%v}`, categoryID, brandID)
    req, _ = http.NewRequest("POST", baseURL+"/products", strings.NewReader(productData))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err = client.Do(req)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), 201, resp.StatusCode)
    
    var productResponse map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&productResponse)
    resp.Body.Close()
    
    productID := productResponse["data"].(map[string]interface{})["id"]
    
    // Step 6: Create product batch
    batchData := fmt.Sprintf(`{"product_id":%v,"batch_number":"BATCH001","quantity":50}`, productID)
    req, _ = http.NewRequest("POST", baseURL+"/product-batches", strings.NewReader(batchData))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err = client.Do(req)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), 201, resp.StatusCode)
    
    var batchResponse map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&batchResponse)
    resp.Body.Close()
    
    batchID := batchResponse["data"].(map[string]interface{})["id"]
    
    // Step 7: Update product batch (should create tracking record)
    updateData := fmt.Sprintf(`{"quantity":75}`)
    req, _ = http.NewRequest("PUT", fmt.Sprintf("%s/product-batches/%v", baseURL, batchID), strings.NewReader(updateData))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err = client.Do(req)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), 200, resp.StatusCode)
    resp.Body.Close()
    
    // Step 8: Verify tracking records
    req, _ = http.NewRequest("GET", fmt.Sprintf("%s/product-batches/%v/tracking", baseURL, batchID), nil)
    req.Header.Set("Authorization", "Bearer "+token)
    
    resp, err = client.Do(req)
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), 200, resp.StatusCode)
    
    var trackingResponse map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&trackingResponse)
    resp.Body.Close()
    
    trackingData := trackingResponse["data"].([]interface{})
    assert.Len(suite.T(), trackingData, 2) // CREATE and UPDATE records
}

func TestWarehouseFlowTestSuite(t *testing.T) {
    suite.Run(t, new(WarehouseFlowTestSuite))
}
```

## 🚀 Performance Testing

### Load Testing
```go
// test/performance/load_test.go
package performance

import (
    "net/http"
    "sync"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
)

func TestUserCreationLoadTest(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping load test in short mode")
    }
    
    baseURL := "http://localhost:8080/api/v1"
    concurrency := 10
    requestsPerWorker := 100
    
    var wg sync.WaitGroup
    var successCount int64
    var failureCount int64
    var mu sync.Mutex
    
    start := time.Now()
    
    // Start concurrent workers
    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            
            for j := 0; j < requestsPerWorker; j++ {
                // Create unique user for each request
                email := fmt.Sprintf("load-test-%d-%d@example.com", workerID, j)
                userData := fmt.Sprintf(`{"email":"%s","password":"password123"}`, email)
                
                resp, err := http.Post(baseURL+"/users", "application/json", strings.NewReader(userData))
                
                mu.Lock()
                if err != nil || resp.StatusCode != 201 {
                    failureCount++
                } else {
                    successCount++
                }
                mu.Unlock()
                
                if resp != nil {
                    resp.Body.Close()
                }
            }
        }(i)
    }
    
    wg.Wait()
    duration := time.Since(start)
    
    totalRequests := int64(concurrency * requestsPerWorker)
    requestsPerSecond := float64(totalRequests) / duration.Seconds()
    
    t.Logf("Load test completed:")
    t.Logf("Total requests: %d", totalRequests)
    t.Logf("Successful requests: %d", successCount)
    t.Logf("Failed requests: %d", failureCount)
    t.Logf("Duration: %v", duration)
    t.Logf("Requests per second: %.2f", requestsPerSecond)
    
    // Assert performance criteria
    assert.True(t, requestsPerSecond > 50, "Should handle at least 50 requests per second")
    assert.True(t, float64(successCount)/float64(totalRequests) > 0.95, "Success rate should be above 95%")
}

func BenchmarkUserCreation(b *testing.B) {
    baseURL := "http://localhost:8080/api/v1"
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        i := 0
        for pb.Next() {
            email := fmt.Sprintf("bench-test-%d@example.com", i)
            userData := fmt.Sprintf(`{"email":"%s","password":"password123"}`, email)
            
            resp, err := http.Post(baseURL+"/users", "application/json", strings.NewReader(userData))
            if err != nil {
                b.Fatal(err)
            }
            resp.Body.Close()
            i++
        }
    })
}
```

## 📊 Test Coverage

### Coverage Configuration
```bash
# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# View coverage summary
go tool cover -func=coverage.out

# Set coverage threshold
go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//' | awk '{if($1 < 80) exit 1}'
```

### Coverage Targets
- **Overall**: Minimum 80% coverage
- **Critical paths**: Minimum 90% coverage
- **Business logic**: Minimum 95% coverage
- **Handlers**: Minimum 85% coverage

## 🔄 CI/CD Testing Pipeline

### GitHub Actions Configuration
```yaml
# .github/workflows/test.yml
name: Test Suite

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: go_wms_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.23.12
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run unit tests
      run: go test -v -race -coverprofile=coverage.out ./...
      env:
        DB_HOST: localhost
        DB_PORT: 5432
        DB_USER: postgres
        DB_PASSWORD: postgres
        DB_NAME: go_wms_test
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
        flags: unittests
        name: codecov-umbrella
    
    - name: Run integration tests
      run: go test -v -tags=integration ./test/integration/...
      env:
        DB_HOST: localhost
        DB_PORT: 5432
        DB_USER: postgres
        DB_PASSWORD: postgres
        DB_NAME: go_wms_test
    
    - name: Run E2E tests
      run: go test -v -tags=e2e ./test/e2e/...
      env:
        DB_HOST: localhost
        DB_PORT: 5432
        DB_USER: postgres
        DB_PASSWORD: postgres
        DB_NAME: go_wms_test
```

## 📋 Testing Best Practices

### Test Organization
1. **Separate test files**: `filename_test.go`
2. **Test suites**: Use testify suites for complex setups
3. **Test tags**: Use build tags for different test types
4. **Clear naming**: Test names should describe the scenario

### Test Data Management
```go
// test/fixtures/users.go
package fixtures

import "your-project/internal/model"

func CreateTestUser() *model.User {
    return &model.User{
        Email:    "test@example.com",
        Password: "hashedpassword",
    }
}

func CreateTestUserWithRoles() *model.User {
    user := CreateTestUser()
    user.Roles = []model.Role{
        {Name: "admin", Description: "Administrator"},
    }
    return user
}
```

### Mock Generation
```bash
# Generate mocks for interfaces
mockgen -source=internal/repository/user_repository.go -destination=mocks/user_repository_mock.go
mockgen -source=internal/service/user_service.go -destination=mocks/user_service_mock.go
```

### Test Commands
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestUserService_CreateUser ./internal/service

# Run tests with race detection
go test -race ./...

# Run tests with coverage
go test -cover ./...

# Run only unit tests
go test -short ./...

# Run only integration tests
go test -tags=integration ./test/integration/...

# Run only E2E tests
go test -tags=e2e ./test/e2e/...

# Benchmark tests
go test -bench=. ./...
```