package repository

import (
	"fmt"

	"gorm.io/gorm"

	"enlabs-task/pkg/model"
)

type UserInterface interface {
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

// Exists checks if a user with the given ID exists
func (repository *UserRepository) Exists(id uint64) (bool, error) {
	var count int64

	result := repository.db.Model(&model.User{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("error checking user existence: %w", result.Error)
	}

	return count > 0, nil
}
