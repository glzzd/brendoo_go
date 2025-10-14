package database

import (
	"product-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initialize creates a new database connection
func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		PrepareStmt: true, // Prepare statements for better performance
	})
	
	if err != nil {
		return nil, err
	}

	// Configure connection pool for high performance
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set maximum number of open connections
	sqlDB.SetMaxOpenConns(100)
	// Set maximum number of idle connections
	sqlDB.SetMaxIdleConns(10)

	return db, nil
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Product{})
	if err != nil {
		return err
	}

	// Create indexes for better performance
	if err := createIndexes(db); err != nil {
		return err
	}

	return nil
}

// createIndexes creates database indexes for performance optimization
func createIndexes(db *gorm.DB) error {
	// Composite index for store and is_active (most common query pattern)
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_products_store_active ON products(store, is_active)").Error; err != nil {
		return err
	}

	// Index for brand filtering
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_products_brand ON products(brand)").Error; err != nil {
		return err
	}

	// Index for category filtering
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_products_category ON products(category)").Error; err != nil {
		return err
	}

	// Index for stock status
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_products_stock_status ON products(stock_status)").Error; err != nil {
		return err
	}

	// Index for created_at for time-based queries
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at)").Error; err != nil {
		return err
	}

	// Partial index for active products only - optimized for main query
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_products_active_only ON products(store, created_at) WHERE is_active = true").Error; err != nil {
		return err
	}

	// GIN index for images field (JSONB type for array operations)
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_products_images_gin ON products USING gin(images)").Error; err != nil {
		return err
	}

	return nil
}