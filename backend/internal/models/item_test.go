package models

import (
	"testing"
	"time"
)

func TestNewItem(t *testing.T) {
	tests := []struct {
		name      string
		itemName  string
		unit      string
		unitPrice float64
		category  string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid item",
			itemName:  "Outlet",
			unit:      "each",
			unitPrice: 15.50,
			category:  "Electrical",
			wantErr:   false,
		},
		{
			name:      "empty name",
			itemName:  "",
			unit:      "each",
			unitPrice: 15.50,
			category:  "Electrical",
			wantErr:   true,
			errMsg:    "item name is required",
		},
		{
			name:      "empty unit",
			itemName:  "Outlet",
			unit:      "",
			unitPrice: 15.50,
			category:  "Electrical",
			wantErr:   true,
			errMsg:    "unit is required",
		},
		{
			name:      "negative price",
			itemName:  "Outlet",
			unit:      "each",
			unitPrice: -10.00,
			category:  "Electrical",
			wantErr:   true,
			errMsg:    "unit price must be positive",
		},
		{
			name:      "zero price allowed",
			itemName:  "Labor",
			unit:      "hour",
			unitPrice: 0,
			category:  "Labor",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := NewItem(tt.itemName, tt.unit, tt.unitPrice, tt.category)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.errMsg {
					t.Errorf("expected error message %q but got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				
				if item == nil {
					t.Fatal("expected item but got nil")
				}
				
				if item.Name != tt.itemName {
					t.Errorf("expected name %q but got %q", tt.itemName, item.Name)
				}
				
				if item.Unit != tt.unit {
					t.Errorf("expected unit %q but got %q", tt.unit, item.Unit)
				}
				
				if item.UnitPrice != tt.unitPrice {
					t.Errorf("expected unit price %f but got %f", tt.unitPrice, item.UnitPrice)
				}
				
				if item.Category != tt.category {
					t.Errorf("expected category %q but got %q", tt.category, item.Category)
				}
				
				if item.ID == "" {
					t.Error("expected ID to be generated")
				}
				
				if item.CreatedAt.IsZero() {
					t.Error("expected CreatedAt to be set")
				}
				
				if item.UpdatedAt.IsZero() {
					t.Error("expected UpdatedAt to be set")
				}
			}
		})
	}
}

func TestItem_Update(t *testing.T) {
	item, err := NewItem("Outlet", "each", 15.50, "Electrical")
	if err != nil {
		t.Fatalf("failed to create item: %v", err)
	}
	
	originalUpdatedAt := item.UpdatedAt
	time.Sleep(10 * time.Millisecond) // Ensure time difference
	
	tests := []struct {
		name      string
		updates   map[string]interface{}
		wantErr   bool
		errMsg    string
	}{
		{
			name: "update name",
			updates: map[string]interface{}{
				"name": "Switch",
			},
			wantErr: false,
		},
		{
			name: "update price",
			updates: map[string]interface{}{
				"unit_price": 20.00,
			},
			wantErr: false,
		},
		{
			name: "update multiple fields",
			updates: map[string]interface{}{
				"name":     "Circuit Breaker",
				"unit":     "box",
				"category": "Panel",
			},
			wantErr: false,
		},
		{
			name: "invalid price",
			updates: map[string]interface{}{
				"unit_price": -5.00,
			},
			wantErr: true,
			errMsg:  "unit price must be positive",
		},
		{
			name: "empty name",
			updates: map[string]interface{}{
				"name": "",
			},
			wantErr: true,
			errMsg:  "item name is required",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy for this test
			testItem := *item
			
			err := testItem.Update(tt.updates)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if err.Error() != tt.errMsg {
					t.Errorf("expected error message %q but got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				
				// Check updates were applied
				if name, ok := tt.updates["name"]; ok {
					if testItem.Name != name.(string) {
						t.Errorf("expected name %q but got %q", name, testItem.Name)
					}
				}
				
				if price, ok := tt.updates["unit_price"]; ok {
					if testItem.UnitPrice != price.(float64) {
						t.Errorf("expected price %f but got %f", price, testItem.UnitPrice)
					}
				}
				
				// UpdatedAt should be newer
				if !testItem.UpdatedAt.After(originalUpdatedAt) {
					t.Error("expected UpdatedAt to be updated")
				}
			}
		})
	}
}