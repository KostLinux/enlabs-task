package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"enlabs-task/pkg/model"
)

// Database wraps the GORM DB connection
type Database struct {
	DB *gorm.DB
}

// NewPostgresDB creates a new connection to PostgreSQL using GORM
func NewPostgresConnection(cfg model.Database) (*Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)

	return &Database{
		DB: db,
	}, nil
}

// Close closes the database connection
func (postgres *Database) Close() error {
	sqlDB, err := postgres.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
