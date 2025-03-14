package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"enlabs-task/pkg/model"
)

// UserRepository handles operations on user data
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(id uint64) (*model.User, error) {
	var user model.User

	query := `SELECT id, username, created_at FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	return &user, nil
}

// Exists checks if a user with the given ID exists
func (r *UserRepository) Exists(id uint64) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := r.db.Get(&exists, query, id)

	if err != nil {
		return false, fmt.Errorf("error checking user existence: %w", err)
	}

	return exists, nil
}
