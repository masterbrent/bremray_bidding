package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// JobTemplate represents a template for creating jobs
type JobTemplate struct {
	ID          string           `json:"id" db:"id"`
	Name        string           `json:"name" db:"name"`
	Description string           `json:"description" db:"description"`
	Items       []TemplateItem   `json:"items"`
	Phases      []TemplatePhase  `json:"phases"`
	IsActive    bool             `json:"isActive" db:"is_active"`
	CreatedAt   time.Time        `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time        `json:"updatedAt" db:"updated_at"`
}

// TemplateItem represents an item in a job template
type TemplateItem struct {
	ID              string  `json:"id" db:"id"`
	TemplateID      string  `json:"templateId" db:"template_id"`
	ItemID          string  `json:"itemId" db:"item_id"`
	DefaultQuantity float64 `json:"defaultQuantity" db:"default_quantity"`
}

// NewJobTemplate creates a new JobTemplate with validation
func NewJobTemplate(name, description string, items []TemplateItem, phases []TemplatePhase) (*JobTemplate, error) {
	if name == "" {
		return nil, errors.New("template name is required")
	}
	if len(items) == 0 {
		return nil, errors.New("template must have at least one item")
	}

	now := time.Now()
	templateID := uuid.New().String()

	// Set template ID on all items
	for i := range items {
		items[i].ID = uuid.New().String()
		items[i].TemplateID = templateID
	}

	// Set template ID on all phases
	for i := range phases {
		phases[i].ID = uuid.New().String()
		phases[i].TemplateID = templateID
	}

	return &JobTemplate{
		ID:          templateID,
		Name:        name,
		Description: description,
		Items:       items,
		Phases:      phases,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// AddItem adds an item to the template
func (t *JobTemplate) AddItem(item TemplateItem) {
	item.ID = uuid.New().String()
	item.TemplateID = t.ID
	t.Items = append(t.Items, item)
	t.UpdatedAt = time.Now()
}

// RemoveItem removes an item from the template
func (t *JobTemplate) RemoveItem(itemID string) error {
	if len(t.Items) <= 1 {
		return errors.New("cannot remove last item from template")
	}

	found := false
	newItems := make([]TemplateItem, 0, len(t.Items)-1)
	
	for _, item := range t.Items {
		if item.ItemID != itemID {
			newItems = append(newItems, item)
		} else {
			found = true
		}
	}

	if !found {
		return errors.New("item not found in template")
	}

	t.Items = newItems
	t.UpdatedAt = time.Now()
	return nil
}

// Deactivate marks the template as inactive
func (t *JobTemplate) Deactivate() {
	t.IsActive = false
	t.UpdatedAt = time.Now()
}

// Activate marks the template as active
func (t *JobTemplate) Activate() {
	t.IsActive = true
	t.UpdatedAt = time.Now()
}