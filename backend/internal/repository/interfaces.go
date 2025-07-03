package repository

import (
	"context"

	"github.com/masterbrent/electrical-bidding-app/internal/models"
)

// ItemRepository defines the interface for item data access
type ItemRepository interface {
	Create(ctx context.Context, item *models.Item) error
	GetByID(ctx context.Context, id string) (*models.Item, error)
	List(ctx context.Context, filter map[string]interface{}) ([]*models.Item, error)
	Update(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
}

// Errors
var (
	ErrNotFound = NewRepositoryError("not found", nil)
	ErrDuplicate = NewRepositoryError("duplicate entry", nil)
	ErrInvalidInput = NewRepositoryError("invalid input", nil)
)

// RepositoryError represents a repository-level error
type RepositoryError struct {
	Message string
	Err     error
}

func NewRepositoryError(message string, err error) *RepositoryError {
	return &RepositoryError{
		Message: message,
		Err:     err,
	}
}

func (e *RepositoryError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}