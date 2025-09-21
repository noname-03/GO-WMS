# 🗃️ Database Documentation
Comprehensive database schema, operations, and management guide for GO-WMS.

## 📊 Database Schema Overview

### Entity Relationship Diagram (ERD)
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│      Users      │    │      Roles      │    │   user_roles    │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ id (PK)         │    │ id (PK)         │    │ user_id (FK)    │
│ email (UNIQUE)  │◄──►│ name            │◄──►│ role_id (FK)    │
│ password        │    │ description     │    └─────────────────┘
│ created_at      │    │ created_at      │
│ updated_at      │    │ updated_at      │
└─────────────────┘    └─────────────────┘
         │
         │ 1:N
         ▼
┌─────────────────┐
│ProductBatchTrack│
├─────────────────┤
│ id (PK)         │
│ product_batch_id│◄──┐
│ user_id (FK)    │   │
│ action          │   │ N:1
│ description     │   │
│ created_at      │   │
└─────────────────┘   │
                      │
┌─────────────────┐   │
│   Categories    │   │
├─────────────────┤   │
│ id (PK)         │   │
│ name (UNIQUE)   │   │
│ description     │   │
│ created_at      │   │
│ updated_at      │   │
└─────────────────┘   │
         │             │
         │ 1:N         │
         ▼             │
┌─────────────────┐   │
│    Products     │   │
├─────────────────┤   │
│ id (PK)         │   │
│ name            │   │
│ description     │   │
│ category_id (FK)│   │
│ brand_id (FK)   │   │
│ created_at      │   │
│ updated_at      │   │
└─────────────────┘   │
         ▲             │
         │ 1:N         │
         │             │
┌─────────────────┐   │
│     Brands      │   │
├─────────────────┤   │
│ id (PK)         │   │
│ name (UNIQUE)   │   │
│ description     │   │
│ created_at      │   │
│ updated_at      │   │
└─────────────────┘   │
                      │
         ┌────────────┘
         │ 1:N
         ▼
┌─────────────────┐
│ ProductBatches  │
├─────────────────┤
│ id (PK)         │
│ product_id (FK) │
│ batch_number    │
│ quantity        │
│ manufacturing_dt│
│ expiry_date     │
│ status          │
│ created_at      │
│ updated_at      │
└─────────────────┘
```

## 🏗️ Table Structures

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
```

**Go Model:**
```go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
    Password  string    `json:"-" gorm:"not null" validate:"required,min=6"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    // Relationships
    Roles             []Role             `json:"roles" gorm:"many2many:user_roles;"`
    ProductBatchTrack []ProductBatchTrack `json:"product_batch_tracks"`
}
```

### Roles Table
```sql
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE UNIQUE INDEX idx_roles_name ON roles(name);
```

**Go Model:**
```go
type Role struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"uniqueIndex;not null" validate:"required"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    // Relationships
    Users []User `json:"users" gorm:"many2many:user_roles;"`
}
```

### User-Roles Junction Table
```sql
CREATE TABLE user_roles (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- Indexes
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
```

### Categories Table
```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE UNIQUE INDEX idx_categories_name ON categories(name);
```

**Go Model:**
```go
type Category struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"uniqueIndex;not null" validate:"required"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    // Relationships
    Products []Product `json:"products"`
}
```

### Brands Table
```sql
CREATE TABLE brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE UNIQUE INDEX idx_brands_name ON brands(name);
```

**Go Model:**
```go
type Brand struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"uniqueIndex;not null" validate:"required"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    // Relationships
    Products []Product `json:"products"`
}
```

### Products Table
```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    brand_id INTEGER REFERENCES brands(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_brand_id ON products(brand_id);
CREATE INDEX idx_products_name ON products(name);
```

**Go Model:**
```go
type Product struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null" validate:"required"`
    Description string    `json:"description"`
    CategoryID  *uint     `json:"category_id"`
    BrandID     *uint     `json:"brand_id"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    
    // Relationships
    Category      *Category      `json:"category" gorm:"foreignKey:CategoryID"`
    Brand         *Brand         `json:"brand" gorm:"foreignKey:BrandID"`
    ProductBatch  []ProductBatch `json:"product_batches"`
}
```

### Product Batches Table
```sql
CREATE TABLE product_batches (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    batch_number VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    manufacturing_date DATE,
    expiry_date DATE,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_product_batches_product_id ON product_batches(product_id);
CREATE INDEX idx_product_batches_batch_number ON product_batches(batch_number);
CREATE INDEX idx_product_batches_status ON product_batches(status);
CREATE INDEX idx_product_batches_expiry_date ON product_batches(expiry_date);

-- Composite indexes for common queries
CREATE INDEX idx_product_batches_product_status ON product_batches(product_id, status);
```

**Go Model:**
```go
type ProductBatch struct {
    ID                uint                 `json:"id" gorm:"primaryKey"`
    ProductID         uint                 `json:"product_id" gorm:"not null"`
    BatchNumber       string               `json:"batch_number" gorm:"not null" validate:"required"`
    Quantity          int                  `json:"quantity" gorm:"default:0"`
    ManufacturingDate *time.Time           `json:"manufacturing_date"`
    ExpiryDate        *time.Time           `json:"expiry_date"`
    Status            string               `json:"status" gorm:"default:active"`
    CreatedAt         time.Time            `json:"created_at"`
    UpdatedAt         time.Time            `json:"updated_at"`
    
    // Relationships
    Product           Product              `json:"product" gorm:"foreignKey:ProductID"`
    ProductBatchTrack []ProductBatchTrack  `json:"product_batch_tracks"`
}
```

### Product Batch Tracking Table
```sql
CREATE TABLE product_batch_tracks (
    id SERIAL PRIMARY KEY,
    product_batch_id INTEGER REFERENCES product_batches(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_product_batch_tracks_batch_id ON product_batch_tracks(product_batch_id);
CREATE INDEX idx_product_batch_tracks_user_id ON product_batch_tracks(user_id);
CREATE INDEX idx_product_batch_tracks_action ON product_batch_tracks(action);
CREATE INDEX idx_product_batch_tracks_created_at ON product_batch_tracks(created_at);

-- Composite indexes
CREATE INDEX idx_pbt_batch_action ON product_batch_tracks(product_batch_id, action);
CREATE INDEX idx_pbt_batch_created ON product_batch_tracks(product_batch_id, created_at DESC);
```

**Go Model:**
```go
type ProductBatchTrack struct {
    ID             uint         `json:"id" gorm:"primaryKey"`
    ProductBatchID uint         `json:"product_batch_id" gorm:"not null"`
    UserID         *uint        `json:"user_id"`
    Action         string       `json:"action" gorm:"not null" validate:"required"`
    Description    string       `json:"description"`
    CreatedAt      time.Time    `json:"created_at"`
    
    // Relationships
    ProductBatch   ProductBatch `json:"product_batch" gorm:"foreignKey:ProductBatchID"`
    User           *User        `json:"user" gorm:"foreignKey:UserID"`
}
```

## 🔧 Database Configuration

### Connection Setup
```go
// database/connection.go
package database

import (
    "fmt"
    "log"
    "os"
    "time"
    
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type Config struct {
    Host            string
    Port            string
    User            string
    Password        string
    DBName          string
    SSLMode         string
    MaxConnections  int
    IdleConnections int
    ConnMaxLifetime time.Duration
}

func NewConfig() *Config {
    return &Config{
        Host:            getEnv("DB_HOST", "localhost"),
        Port:            getEnv("DB_PORT", "5432"),
        User:            getEnv("DB_USER", "postgres"),
        Password:        getEnv("DB_PASSWORD", ""),
        DBName:          getEnv("DB_NAME", "go_wms"),
        SSLMode:         getEnv("DB_SSLMODE", "disable"),
        MaxConnections:  getEnvAsInt("DB_MAX_CONNECTIONS", 25),
        IdleConnections: getEnvAsInt("DB_IDLE_CONNECTIONS", 5),
        ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", time.Hour),
    }
}

func Connect(config *Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    
    // Configure connection pool
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
    }
    
    sqlDB.SetMaxOpenConns(config.MaxConnections)
    sqlDB.SetMaxIdleConns(config.IdleConnections)
    sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
    
    // Test connection
    if err := sqlDB.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    log.Println("Database connected successfully")
    return db, nil
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}
```

### Environment Variables
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=go_wms
DB_SSLMODE=disable

# Connection Pool Settings
DB_MAX_CONNECTIONS=25
DB_IDLE_CONNECTIONS=5
DB_CONN_MAX_LIFETIME=1h

# Development Settings
DB_LOG_LEVEL=info
DB_AUTO_MIGRATE=true
```

## 🔄 Database Migrations

### Auto-Migration Setup
```go
// database/migration.go
package database

import (
    "log"
    
    "gorm.io/gorm"
    "your-project/internal/model"
)

func AutoMigrate(db *gorm.DB) error {
    log.Println("Starting database migration...")
    
    err := db.AutoMigrate(
        &model.User{},
        &model.Role{},
        &model.Category{},
        &model.Brand{},
        &model.Product{},
        &model.ProductBatch{},
        &model.ProductBatchTrack{},
    )
    
    if err != nil {
        return fmt.Errorf("migration failed: %w", err)
    }
    
    log.Println("Database migration completed successfully")
    return nil
}

// Manual migration functions for specific changes
func MigrateAddIndexes(db *gorm.DB) error {
    // Add custom indexes
    queries := []string{
        "CREATE INDEX IF NOT EXISTS idx_product_batches_composite ON product_batches(product_id, status);",
        "CREATE INDEX IF NOT EXISTS idx_pbt_batch_created ON product_batch_tracks(product_batch_id, created_at DESC);",
    }
    
    for _, query := range queries {
        if err := db.Exec(query).Error; err != nil {
            return fmt.Errorf("failed to execute migration query: %w", err)
        }
    }
    
    return nil
}
```

### Version-Based Migrations (Advanced)
```go
// migrations/001_create_users_table.go
package migrations

import "gorm.io/gorm"

type Migration001 struct{}

func (m Migration001) Up(db *gorm.DB) error {
    return db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `).Error
}

func (m Migration001) Down(db *gorm.DB) error {
    return db.Exec("DROP TABLE IF EXISTS users;").Error
}

func (m Migration001) Version() string {
    return "001"
}
```

## 🌱 Database Seeding

### Seeder Registry
```go
// database/seeder/registry.go
package seeder

import (
    "gorm.io/gorm"
)

type Seeder interface {
    Seed() error
}

func GetSeeders(db *gorm.DB) []Seeder {
    return []Seeder{
        &RoleSeeder{db: db},
        &UserSeeder{db: db},
        &CategorySeeder{db: db},
        &BrandSeeder{db: db},
        &ProductSeeder{db: db},
        &ProductBatchSeeder{db: db},
        &ProductBatchTrackSeeder{db: db},
    }
}

func RunSeeders(db *gorm.DB) error {
    seeders := GetSeeders(db)
    
    for _, seeder := range seeders {
        if err := seeder.Seed(); err != nil {
            return fmt.Errorf("seeding failed: %w", err)
        }
    }
    
    log.Println("Database seeding completed successfully")
    return nil
}
```

### Example Seeders
```go
// database/seeder/role_seeder.go
package seeder

import (
    "gorm.io/gorm"
    "your-project/internal/model"
)

type RoleSeeder struct {
    db *gorm.DB
}

func (s *RoleSeeder) Seed() error {
    roles := []model.Role{
        {Name: "admin", Description: "System Administrator"},
        {Name: "manager", Description: "Warehouse Manager"},
        {Name: "operator", Description: "Warehouse Operator"},
        {Name: "viewer", Description: "Read-only Access"},
    }
    
    for _, role := range roles {
        var existingRole model.Role
        if err := s.db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                if err := s.db.Create(&role).Error; err != nil {
                    return err
                }
            } else {
                return err
            }
        }
    }
    
    return nil
}
```

## 📊 Database Queries

### Common Query Patterns

#### Basic CRUD Operations
```go
// Repository examples
type UserRepository struct {
    db *gorm.DB
}

// Create
func (r *UserRepository) Create(user *model.User) error {
    return r.db.Create(user).Error
}

// Read with relationships
func (r *UserRepository) GetByIDWithRoles(id uint) (*model.User, error) {
    var user model.User
    err := r.db.Preload("Roles").First(&user, id).Error
    return &user, err
}

// Update
func (r *UserRepository) Update(user *model.User) error {
    return r.db.Save(user).Error
}

// Delete (soft delete)
func (r *UserRepository) Delete(id uint) error {
    return r.db.Delete(&model.User{}, id).Error
}

// Hard delete
func (r *UserRepository) HardDelete(id uint) error {
    return r.db.Unscoped().Delete(&model.User{}, id).Error
}
```

#### Complex Queries
```go
// Product with batches and tracking
func (r *ProductRepository) GetProductWithFullTracking(id uint) (*model.Product, error) {
    var product model.Product
    err := r.db.
        Preload("Category").
        Preload("Brand").
        Preload("ProductBatch.ProductBatchTrack.User").
        First(&product, id).Error
    return &product, err
}

// Products by category with pagination
func (r *ProductRepository) GetByCategory(categoryID uint, page, limit int) ([]*model.Product, int64, error) {
    var products []*model.Product
    var total int64
    
    query := r.db.Where("category_id = ?", categoryID)
    
    // Count total
    if err := query.Model(&model.Product{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Get paginated results
    offset := (page - 1) * limit
    err := query.
        Preload("Category").
        Preload("Brand").
        Offset(offset).
        Limit(limit).
        Find(&products).Error
    
    return products, total, err
}

// Expired batches
func (r *ProductBatchRepository) GetExpiredBatches() ([]*model.ProductBatch, error) {
    var batches []*model.ProductBatch
    err := r.db.
        Where("expiry_date < ? AND status = ?", time.Now(), "active").
        Preload("Product").
        Find(&batches).Error
    return batches, err
}

// Low stock batches
func (r *ProductBatchRepository) GetLowStockBatches(threshold int) ([]*model.ProductBatch, error) {
    var batches []*model.ProductBatch
    err := r.db.
        Where("quantity < ? AND status = ?", threshold, "active").
        Preload("Product").
        Find(&batches).Error
    return batches, err
}
```

#### Advanced Queries with Raw SQL
```go
// Complex reporting query
func (r *ProductBatchRepository) GetInventoryReport() ([]InventoryReport, error) {
    type InventoryReport struct {
        ProductName    string `json:"product_name"`
        CategoryName   string `json:"category_name"`
        BrandName      string `json:"brand_name"`
        TotalQuantity  int    `json:"total_quantity"`
        ActiveBatches  int    `json:"active_batches"`
        ExpiredBatches int    `json:"expired_batches"`
    }
    
    var reports []InventoryReport
    
    query := `
        SELECT 
            p.name as product_name,
            c.name as category_name,
            b.name as brand_name,
            COALESCE(SUM(CASE WHEN pb.status = 'active' THEN pb.quantity END), 0) as total_quantity,
            COUNT(CASE WHEN pb.status = 'active' THEN 1 END) as active_batches,
            COUNT(CASE WHEN pb.expiry_date < NOW() THEN 1 END) as expired_batches
        FROM products p
        LEFT JOIN categories c ON p.category_id = c.id
        LEFT JOIN brands b ON p.brand_id = b.id
        LEFT JOIN product_batches pb ON p.id = pb.product_id
        GROUP BY p.id, p.name, c.name, b.name
        ORDER BY p.name
    `
    
    err := r.db.Raw(query).Scan(&reports).Error
    return reports, err
}
```

## 🔍 Database Indexing Strategy

### Performance Indexes
```sql
-- Primary indexes (automatically created)
-- id columns are automatically indexed as PRIMARY KEY

-- Unique indexes (automatically created)
-- email, name fields with unique constraints

-- Custom performance indexes
CREATE INDEX idx_products_name_gin ON products USING gin(to_tsvector('english', name));
CREATE INDEX idx_product_batches_expiry_status ON product_batches(expiry_date, status) WHERE status = 'active';
CREATE INDEX idx_pbt_created_desc ON product_batch_tracks(created_at DESC);

-- Partial indexes for common filtered queries
CREATE INDEX idx_active_batches ON product_batches(product_id) WHERE status = 'active';
CREATE INDEX idx_expired_batches ON product_batches(expiry_date) WHERE expiry_date < NOW();

-- Composite indexes for join queries
CREATE INDEX idx_products_category_brand ON products(category_id, brand_id);
CREATE INDEX idx_pbt_batch_user_action ON product_batch_tracks(product_batch_id, user_id, action);
```

### Index Monitoring
```go
// Check index usage
func (r *Repository) AnalyzeIndexUsage(db *gorm.DB) error {
    query := `
        SELECT 
            schemaname,
            tablename,
            indexname,
            idx_tup_read,
            idx_tup_fetch
        FROM pg_stat_user_indexes 
        ORDER BY idx_tup_read DESC;
    `
    
    var results []map[string]interface{}
    return db.Raw(query).Scan(&results).Error
}
```

## 🔐 Database Security

### User Permissions
```sql
-- Create application user with limited permissions
CREATE USER go_wms_app WITH PASSWORD 'secure_password';

-- Grant specific permissions
GRANT CONNECT ON DATABASE go_wms TO go_wms_app;
GRANT USAGE ON SCHEMA public TO go_wms_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO go_wms_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO go_wms_app;

-- Revoke dangerous permissions
REVOKE CREATE ON SCHEMA public FROM go_wms_app;
REVOKE ALL PRIVILEGES ON pg_catalog FROM go_wms_app;
```

### Connection Security
```go
// SSL configuration for production
func ConnectSecure(config *Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=require sslcert=%s sslkey=%s sslrootcert=%s",
        config.Host, config.Port, config.User, config.Password, config.DBName,
        config.SSLCert, config.SSLKey, config.SSLRootCert,
    )
    
    return gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent), // Don't log SQL in production
    })
}
```

## 💾 Backup and Recovery

### Automated Backup Script
```bash
#!/bin/bash
# backup-database.sh

DB_NAME="go_wms"
DB_USER="postgres"
BACKUP_DIR="/opt/backups/database"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="$BACKUP_DIR/${DB_NAME}_backup_$TIMESTAMP.sql"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Create database backup
pg_dump -h localhost -U $DB_USER -d $DB_NAME > $BACKUP_FILE

# Compress backup
gzip $BACKUP_FILE

# Remove backups older than 30 days
find $BACKUP_DIR -name "${DB_NAME}_backup_*.sql.gz" -mtime +30 -delete

echo "Backup completed: ${BACKUP_FILE}.gz"
```

### Recovery Procedures
```bash
# Restore from backup
gunzip go_wms_backup_20240315_120000.sql.gz
psql -h localhost -U postgres -d go_wms < go_wms_backup_20240315_120000.sql

# Point-in-time recovery (if WAL archiving is enabled)
pg_basebackup -h localhost -D /opt/recovery -U postgres -v -P -W -X stream
```

## 📈 Performance Optimization

### Connection Pool Tuning
```go
// Optimal connection pool settings
func ConfigureConnectionPool(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    
    // Maximum number of open connections
    sqlDB.SetMaxOpenConns(25)
    
    // Maximum number of idle connections
    sqlDB.SetMaxIdleConns(5)
    
    // Maximum lifetime of a connection
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return nil
}
```

### Query Optimization
```go
// Use proper preloading to avoid N+1 queries
func (r *ProductRepository) GetAllWithBatches() ([]*model.Product, error) {
    var products []*model.Product
    
    // Bad - causes N+1 queries
    // r.db.Find(&products)
    // for _, product := range products {
    //     r.db.Where("product_id = ?", product.ID).Find(&product.ProductBatch)
    // }
    
    // Good - single query with joins
    err := r.db.Preload("ProductBatch").Find(&products).Error
    return products, err
}

// Use select to limit returned columns
func (r *UserRepository) GetUsersList() ([]*model.User, error) {
    var users []*model.User
    err := r.db.Select("id", "email", "created_at").Find(&users).Error
    return users, err
}

// Use exists for checking relationships
func (r *ProductRepository) HasActiveBatches(productID uint) (bool, error) {
    var exists bool
    err := r.db.Raw(
        "SELECT EXISTS(SELECT 1 FROM product_batches WHERE product_id = ? AND status = 'active')",
        productID,
    ).Scan(&exists).Error
    return exists, err
}
```

## 📊 Database Monitoring

### Performance Monitoring Queries
```sql
-- Check slow queries
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    max_time
FROM pg_stat_statements 
ORDER BY mean_time DESC 
LIMIT 10;

-- Check table sizes
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size(tablename::regclass)) as size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(tablename::regclass) DESC;

-- Check index usage
SELECT 
    indexrelname,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes 
ORDER BY idx_tup_read DESC;

-- Check connection status
SELECT 
    state,
    count(*) 
FROM pg_stat_activity 
GROUP BY state;
```

### Application-Level Monitoring
```go
// Database health check
func (r *Repository) HealthCheck() error {
    sqlDB, err := r.db.DB()
    if err != nil {
        return err
    }
    
    // Check if database is reachable
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    return sqlDB.PingContext(ctx)
}

// Monitor connection pool stats
func (r *Repository) GetConnectionStats() map[string]interface{} {
    sqlDB, _ := r.db.DB()
    stats := sqlDB.Stats()
    
    return map[string]interface{}{
        "max_open_connections":     stats.MaxOpenConnections,
        "open_connections":         stats.OpenConnections,
        "in_use":                  stats.InUse,
        "idle":                    stats.Idle,
        "wait_count":              stats.WaitCount,
        "wait_duration":           stats.WaitDuration,
        "max_idle_closed":         stats.MaxIdleClosed,
        "max_lifetime_closed":     stats.MaxLifetimeClosed,
    }
}
```