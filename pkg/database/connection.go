package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"enlabs-task/pkg/model"
)

// PostgresDB represents the database connection
type PostgresDB struct {
	DB *sqlx.DB
}

// NewPostgresDB creates a new database connection
func NewPostgresDB(cfg model.Database) (*PostgresDB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)

	return &PostgresDB{DB: db}, nil
}

// Close closes the database connection
func (postgres *PostgresDB) Close() error {
	return postgres.DB.Close()
}
