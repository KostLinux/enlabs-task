package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"enlabs-task/pkg/model"
)

type UserInterface interface {
	GetByID(id uint64) (*model.User, error)
	Exists(id uint64) (bool, error)
}

// UserRepository handles operations on user data
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID retrieves a user by their ID
func (repository *UserRepository) GetByID(id uint64) (*model.User, error) {
	var user model.User

	result := repository.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}

		return nil, fmt.Errorf("error fetching user: %w", result.Error)
	}

	return &user, nil
}

// Exists checks if a user with the given ID exists
func (repository *UserRepository) Exists(id uint64) (bool, error) {
	var count int64

	result := repository.db.Model(&model.User{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("error checking user existence: %w", result.Error)
	}

	return count > 0, nil
}
