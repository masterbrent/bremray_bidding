package repository

import (
	"context"
	"testing"

	"github.com/masterbrent/electrical-bidding-app/internal/models"
)

// MockItemRepository is a mock implementation for testing
type MockItemRepository struct {
	items  map[string]*models.Item
	err    error
	calls  map[string]int
}

func NewMockItemRepository() *MockItemRepository {
	return &MockItemRepository{
		items: make(map[string]*models.Item),
		calls: make(map[string]int),
	}
}

func (m *MockItemRepository) SetError(err error) {
	m.err = err
}

func (m *MockItemRepository) GetCallCount(method string) int {
	return m.calls[method]
}

func TestItemRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		item    *models.Item
		wantErr bool
		setup   func(*MockItemRepository)
	}{
		{
			name: "successful create",
			item: &models.Item{
				ID:        "test-id",
				Name:      "Test Item",
				Unit:      "each",
				UnitPrice: 10.00,
				Category:  "Test",
			},
			wantErr: false,
		},
		{
			name: "nil item",
			item: nil,
			wantErr: true,
		},
		{
			name: "duplicate ID",
			item: &models.Item{
				ID:   "existing-id",
				Name: "New Item",
			},
			wantErr: true,
			setup: func(m *MockItemRepository) {
				m.items["existing-id"] = &models.Item{ID: "existing-id"}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockItemRepository()
			if tt.setup != nil {
				tt.setup(repo)
			}

			itemRepo := NewItemRepository(nil) // We'll implement this
			
			// For now, let's define the interface
			if tt.item != nil && !tt.wantErr {
				// Mock successful creation
				repo.items[tt.item.ID] = tt.item
			}
		})
	}
}

func TestItemRepository_GetByID(t *testing.T) {
	testItem := &models.Item{
		ID:        "test-id",
		Name:      "Test Item",
		Unit:      "each",
		UnitPrice: 15.50,
		Category:  "Electrical",
	}

	tests := []struct {
		name    string
		id      string
		want    *models.Item
		wantErr bool
		setup   func(*MockItemRepository)
	}{
		{
			name: "existing item",
			id:   "test-id",
			want: testItem,
			wantErr: false,
			setup: func(m *MockItemRepository) {
				m.items["test-id"] = testItem
			},
		},
		{
			name:    "non-existing item",
			id:      "non-existing",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty id",
			id:      "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockItemRepository()
			if tt.setup != nil {
				tt.setup(repo)
			}

			// Test retrieval
			got, exists := repo.items[tt.id]
			if tt.wantErr && exists {
				t.Errorf("expected error but found item")
			}
			if !tt.wantErr && !exists {
				t.Errorf("expected item but not found")
			}
			if tt.want != nil && got != nil {
				if got.ID != tt.want.ID {
					t.Errorf("expected ID %s but got %s", tt.want.ID, got.ID)
				}
			}
		})
	}
}

func TestItemRepository_List(t *testing.T) {
	items := []*models.Item{
		{ID: "1", Name: "Item 1", Category: "Electrical"},
		{ID: "2", Name: "Item 2", Category: "Wire"},
		{ID: "3", Name: "Item 3", Category: "Electrical"},
	}

	tests := []struct {
		name     string
		filter   map[string]interface{}
		wantLen  int
		setup    func(*MockItemRepository)
	}{
		{
			name:    "list all",
			filter:  nil,
			wantLen: 3,
			setup: func(m *MockItemRepository) {
				for _, item := range items {
					m.items[item.ID] = item
				}
			},
		},
		{
			name:    "filter by category",
			filter:  map[string]interface{}{"category": "Electrical"},
			wantLen: 2,
			setup: func(m *MockItemRepository) {
				for _, item := range items {
					m.items[item.ID] = item
				}
			},
		},
		{
			name:    "empty repository",
			filter:  nil,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockItemRepository()
			if tt.setup != nil {
				tt.setup(repo)
			}

			// Count items matching filter
			count := 0
			for _, item := range repo.items {
				if tt.filter == nil {
					count++
				} else if cat, ok := tt.filter["category"]; ok && item.Category == cat {
					count++
				}
			}

			if count != tt.wantLen {
				t.Errorf("expected %d items but got %d", tt.wantLen, count)
			}
		})
	}
}

func TestItemRepository_Update(t *testing.T) {
	originalItem := &models.Item{
		ID:        "test-id",
		Name:      "Original Name",
		Unit:      "each",
		UnitPrice: 10.00,
		Category:  "Test",
	}

	tests := []struct {
		name    string
		id      string
		updates map[string]interface{}
		wantErr bool
		setup   func(*MockItemRepository)
	}{
		{
			name: "successful update",
			id:   "test-id",
			updates: map[string]interface{}{
				"name":       "Updated Name",
				"unit_price": 15.00,
			},
			wantErr: false,
			setup: func(m *MockItemRepository) {
				m.items["test-id"] = &models.Item{
					ID:        originalItem.ID,
					Name:      originalItem.Name,
					Unit:      originalItem.Unit,
					UnitPrice: originalItem.UnitPrice,
					Category:  originalItem.Category,
				}
			},
		},
		{
			name:    "non-existing item",
			id:      "non-existing",
			updates: map[string]interface{}{"name": "New Name"},
			wantErr: true,
		},
		{
			name:    "empty id",
			id:      "",
			updates: map[string]interface{}{"name": "New Name"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockItemRepository()
			if tt.setup != nil {
				tt.setup(repo)
			}

			// Mock update
			if item, exists := repo.items[tt.id]; exists && !tt.wantErr {
				err := item.Update(tt.updates)
				if err != nil && !tt.wantErr {
					t.Errorf("unexpected error: %v", err)
				}
			} else if !exists && !tt.wantErr {
				t.Errorf("expected item to exist for update")
			}
		})
	}
}

func TestItemRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
		setup   func(*MockItemRepository)
		verify  func(*MockItemRepository, *testing.T)
	}{
		{
			name:    "successful delete",
			id:      "test-id",
			wantErr: false,
			setup: func(m *MockItemRepository) {
				m.items["test-id"] = &models.Item{ID: "test-id"}
			},
			verify: func(m *MockItemRepository, t *testing.T) {
				if _, exists := m.items["test-id"]; exists {
					t.Error("expected item to be deleted")
				}
			},
		},
		{
			name:    "non-existing item",
			id:      "non-existing",
			wantErr: true,
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockItemRepository()
			if tt.setup != nil {
				tt.setup(repo)
			}

			// Mock delete
			if tt.id != "" {
				if _, exists := repo.items[tt.id]; exists {
					delete(repo.items, tt.id)
				} else if !tt.wantErr {
					t.Error("expected item to exist for deletion")
				}
			}

			if tt.verify != nil {
				tt.verify(repo, t)
			}
		})
	}
}