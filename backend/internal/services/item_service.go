package services

import (
	"context"

	"github.com/masterbrent/electrical-bidding-app/internal/models"
	"github.com/masterbrent/electrical-bidding-app/internal/repository"
)

// ItemService implements business logic for items
type ItemService struct {
	repo repository.ItemRepository
}

// NewItemService creates a new item service
func NewItemService(repo repository.ItemRepository) *ItemService {
	return &ItemService{
		repo: repo,
	}
}

// Create creates a new item
func (s *ItemService) Create(ctx context.Context, name, unit string, unitPrice float64, category string) (*models.Item, error) {
	// Create the item with validation
	item, err := models.NewItem(name, unit, unitPrice, category)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

// GetByID retrieves an item by ID
func (s *ItemService) GetByID(ctx context.Context, id string) (*models.Item, error) {
	return s.repo.GetByID(ctx, id)
}

// List retrieves all items
func (s *ItemService) List(ctx context.Context) ([]*models.Item, error) {
	return s.repo.List(ctx, nil)
}

// Update updates an existing item
func (s *ItemService) Update(ctx context.Context, id string, updates map[string]interface{}) (*models.Item, error) {
	// Update through repository (which handles validation)
	if err := s.repo.Update(ctx, id, updates); err != nil {
		return nil, err
	}

	// Get the updated item
	return s.repo.GetByID(ctx, id)
}

// Delete removes an item
func (s *ItemService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}