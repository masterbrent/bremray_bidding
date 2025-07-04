package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrValidation = errors.New("validation error")
)

// Item represents an inventory item that can be used in jobs
type Item struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Nickname  string    `json:"nickname,omitempty" db:"nickname"` // Display name for job cards
	Unit      string    `json:"unit" db:"unit"`
	UnitPrice float64   `json:"unitPrice" db:"unit_price"`
	Category  string    `json:"category" db:"category"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// NewItem creates a new Item with validation
func NewItem(name, unit string, unitPrice float64, category string) (*Item, error) {
	if err := validateItemFields(name, unit, unitPrice); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Item{
		ID:        uuid.New().String(),
		Name:      name,
		Unit:      unit,
		UnitPrice: unitPrice,
		Category:  category,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Update updates item fields with validation
func (i *Item) Update(updates map[string]interface{}) error {
	// Extract values for validation
	name := i.Name
	unit := i.Unit
	unitPrice := i.UnitPrice

	// Apply updates to temp variables first
	if val, ok := updates["name"]; ok {
		name = val.(string)
	}
	if val, ok := updates["unit"]; ok {
		unit = val.(string)
	}
	if val, ok := updates["unit_price"]; ok {
		unitPrice = val.(float64)
	}

	// Validate all fields
	if err := validateItemFields(name, unit, unitPrice); err != nil {
		return err
	}

	// Apply updates if validation passes
	if val, ok := updates["name"]; ok {
		i.Name = val.(string)
	}
	if val, ok := updates["nickname"]; ok {
		i.Nickname = val.(string)
	}
	if val, ok := updates["unit"]; ok {
		i.Unit = val.(string)
	}
	if val, ok := updates["unit_price"]; ok {
		i.UnitPrice = val.(float64)
	}
	if val, ok := updates["category"]; ok {
		i.Category = val.(string)
	}

	i.UpdatedAt = time.Now()
	return nil
}

// validateItemFields validates item fields
func validateItemFields(name, unit string, unitPrice float64) error {
	if name == "" {
		return errors.New("item name is required")
	}
	if unit == "" {
		return errors.New("unit is required")
	}
	if unitPrice < 0 {
		return errors.New("unit price must be positive")
	}
	return nil
}