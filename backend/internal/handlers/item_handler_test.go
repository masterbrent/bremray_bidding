package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/masterbrent/electrical-bidding-app/internal/models"
	"github.com/masterbrent/electrical-bidding-app/internal/repository"
)

// MockItemService implements the item service interface for testing
type MockItemService struct {
	items    map[string]*models.Item
	err      error
	lastCall string
}

func NewMockItemService() *MockItemService {
	return &MockItemService{
		items: make(map[string]*models.Item),
	}
}

func (m *MockItemService) Create(ctx context.Context, name, unit string, unitPrice float64, category string) (*models.Item, error) {
	m.lastCall = "Create"
	if m.err != nil {
		return nil, m.err
	}
	
	item, err := models.NewItem(name, unit, unitPrice, category)
	if err != nil {
		return nil, err
	}
	
	m.items[item.ID] = item
	return item, nil
}

func (m *MockItemService) GetByID(ctx context.Context, id string) (*models.Item, error) {
	m.lastCall = "GetByID"
	if m.err != nil {
		return nil, m.err
	}
	
	item, exists := m.items[id]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return item, nil
}

func (m *MockItemService) List(ctx context.Context) ([]*models.Item, error) {
	m.lastCall = "List"
	if m.err != nil {
		return nil, m.err
	}
	
	items := make([]*models.Item, 0, len(m.items))
	for _, item := range m.items {
		items = append(items, item)
	}
	return items, nil
}

func (m *MockItemService) Update(ctx context.Context, id string, updates map[string]interface{}) (*models.Item, error) {
	m.lastCall = "Update"
	if m.err != nil {
		return nil, m.err
	}
	
	item, exists := m.items[id]
	if !exists {
		return nil, repository.ErrNotFound
	}
	
	// Create a copy to avoid modifying the original if validation fails
	itemCopy := *item
	if err := itemCopy.Update(updates); err != nil {
		return nil, err
	}
	
	// Only update the stored item if validation passed
	m.items[id] = &itemCopy
	return &itemCopy, nil
}

func (m *MockItemService) Delete(ctx context.Context, id string) error {
	m.lastCall = "Delete"
	if m.err != nil {
		return m.err
	}
	
	if _, exists := m.items[id]; !exists {
		return repository.ErrNotFound
	}
	
	delete(m.items, id)
	return nil
}

func TestItemHandler_Create(t *testing.T) {
	tests := []struct {
		name         string
		body         interface{}
		expectedCode int
		setupMock    func(*MockItemService)
	}{
		{
			name: "valid item creation",
			body: map[string]interface{}{
				"name":      "Test Item",
				"unit":      "each",
				"unitPrice": 15.50,
				"category":  "Electrical",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "missing required field",
			body: map[string]interface{}{
				"unit":      "each",
				"unitPrice": 15.50,
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid price",
			body: map[string]interface{}{
				"name":      "Test Item",
				"unit":      "each",
				"unitPrice": -10.00,
				"category":  "Electrical",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid JSON",
			body:         "invalid json",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			service := NewMockItemService()
			if tt.setupMock != nil {
				tt.setupMock(service)
			}
			
			handler := NewItemHandler(service)
			
			// Create request
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			// Record response
			w := httptest.NewRecorder()
			
			// Execute
			handler.Create(w, req)
			
			// Assert
			if w.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, w.Code)
			}
			
			// If successful, check response
			if w.Code == http.StatusCreated {
				var response map[string]interface{}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				
				if response["id"] == "" {
					t.Error("expected ID in response")
				}
			}
		})
	}
}

func TestItemHandler_GetByID(t *testing.T) {
	testItem := &models.Item{
		ID:        "test-id",
		Name:      "Test Item",
		Unit:      "each",
		UnitPrice: 15.50,
		Category:  "Electrical",
	}

	tests := []struct {
		name         string
		id           string
		expectedCode int
		setupMock    func(*MockItemService)
	}{
		{
			name:         "existing item",
			id:           "test-id",
			expectedCode: http.StatusOK,
			setupMock: func(m *MockItemService) {
				m.items["test-id"] = testItem
			},
		},
		{
			name:         "non-existing item",
			id:           "non-existing",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "empty id",
			id:           "",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			service := NewMockItemService()
			if tt.setupMock != nil {
				tt.setupMock(service)
			}
			
			handler := NewItemHandler(service)
			
			// Create request with gorilla/mux
			req := httptest.NewRequest(http.MethodGet, "/api/items/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			
			// Record response
			w := httptest.NewRecorder()
			
			// Execute
			handler.GetByID(w, req)
			
			// Assert
			if w.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, w.Code)
			}
			
			// If successful, check response
			if w.Code == http.StatusOK {
				var response models.Item
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				
				if response.ID != testItem.ID {
					t.Errorf("expected ID %s, got %s", testItem.ID, response.ID)
				}
			}
		})
	}
}

func TestItemHandler_List(t *testing.T) {
	tests := []struct {
		name         string
		expectedCode int
		setupMock    func(*MockItemService)
		expectedLen  int
	}{
		{
			name:         "empty list",
			expectedCode: http.StatusOK,
			expectedLen:  0,
		},
		{
			name:         "list with items",
			expectedCode: http.StatusOK,
			expectedLen:  2,
			setupMock: func(m *MockItemService) {
				item1, _ := models.NewItem("Item 1", "each", 10.00, "Test")
				item2, _ := models.NewItem("Item 2", "box", 20.00, "Test")
				m.items[item1.ID] = item1
				m.items[item2.ID] = item2
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			service := NewMockItemService()
			if tt.setupMock != nil {
				tt.setupMock(service)
			}
			
			handler := NewItemHandler(service)
			
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/api/items", nil)
			
			// Record response
			w := httptest.NewRecorder()
			
			// Execute
			handler.List(w, req)
			
			// Assert
			if w.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, w.Code)
			}
			
			// Check response
			var response []models.Item
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			
			if len(response) != tt.expectedLen {
				t.Errorf("expected %d items, got %d", tt.expectedLen, len(response))
			}
		})
	}
}

func TestItemHandler_Update(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		body         interface{}
		expectedCode int
		setupMock    func(*MockItemService)
	}{
		{
			name: "successful update",
			id:   "test-id",
			body: map[string]interface{}{
				"name":       "Updated Item",
				"unit_price": 25.00,
			},
			expectedCode: http.StatusOK,
			setupMock: func(m *MockItemService) {
				item, _ := models.NewItem("Original Item", "each", 15.00, "Test")
				item.ID = "test-id"
				m.items["test-id"] = item
			},
		},
		{
			name: "non-existing item",
			id:   "non-existing",
			body: map[string]interface{}{
				"name": "Updated",
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "invalid update",
			id:   "test-id",
			body: map[string]interface{}{
				"unit_price": -10.00,
			},
			expectedCode: http.StatusBadRequest,
			setupMock: func(m *MockItemService) {
				item, _ := models.NewItem("Item", "each", 15.00, "Test")
				item.ID = "test-id"
				m.items["test-id"] = item
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			service := NewMockItemService()
			if tt.setupMock != nil {
				tt.setupMock(service)
			}
			
			handler := NewItemHandler(service)
			
			// Create request
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/api/items/"+tt.id, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			
			// Record response
			w := httptest.NewRecorder()
			
			// Execute
			handler.Update(w, req)
			
			// Assert
			if w.Code != tt.expectedCode {
				body := w.Body.String()
				t.Errorf("expected status code %d, got %d. Response: %s", tt.expectedCode, w.Code, body)
			}
		})
	}
}

func TestItemHandler_Delete(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		expectedCode int
		setupMock    func(*MockItemService)
	}{
		{
			name:         "successful delete",
			id:           "test-id",
			expectedCode: http.StatusNoContent,
			setupMock: func(m *MockItemService) {
				item, _ := models.NewItem("Item", "each", 15.00, "Test")
				item.ID = "test-id"
				m.items["test-id"] = item
			},
		},
		{
			name:         "non-existing item",
			id:           "non-existing",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "empty id",
			id:           "",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			service := NewMockItemService()
			if tt.setupMock != nil {
				tt.setupMock(service)
			}
			
			handler := NewItemHandler(service)
			
			// Create request
			req := httptest.NewRequest(http.MethodDelete, "/api/items/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			
			// Record response
			w := httptest.NewRecorder()
			
			// Execute
			handler.Delete(w, req)
			
			// Assert
			if w.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, w.Code)
			}
		})
	}
}