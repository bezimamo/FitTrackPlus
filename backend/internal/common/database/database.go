package database

import (
	"fmt"
	"log"

	"fittrackplus/internal/common/config"
	"fittrackplus/internal/common/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is a global variable to hold our database connection
// In Go, we can use global variables for database connections
var DB *gorm.DB

// Connect establishes a connection to the PostgreSQL database
func Connect(config *config.Config) error {
	// Create the database connection string
	// This follows the PostgreSQL connection format
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	// Open the database connection using GORM
	// GORM is an ORM (Object-Relational Mapping) library for Go
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Configure GORM logger for development
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Store the database connection in our global variable
	DB = db

	log.Println("✅ Database connected successfully!")

	// Auto-migrate our models to create tables
	// This will create tables based on our struct definitions
	err = AutoMigrate()
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("✅ Database tables created successfully!")
	return nil
}

// AutoMigrate creates database tables based on our models
func AutoMigrate() error {
	// GORM will automatically create tables for all our models
	// The table names will be the plural form of the struct name
	return DB.AutoMigrate(
		&models.User{},
		&models.UserProfile{},
		&models.Plan{},
		&models.UserPlan{},
		&models.ProgressLog{},
		&models.Booking{},
		&models.Payment{},
	)
}

// GetDB returns the database connection
// This function allows other parts of our application to access the database
func GetDB() *gorm.DB {
	return DB
}

// Close closes the database connection
// This is useful for graceful shutdown
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
} 