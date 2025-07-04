package models

import (
	"testing"
)

func TestNewJobTemplate(t *testing.T) {
	tests := []struct {
		name        string
		templateName string
		description string
		items       []TemplateItem
		wantErr     bool
		errMsg      string
	}{
		{
			name:         "valid template",
			templateName: "Basic Outlet Installation",
			description:  "Install standard electrical outlet",
			items: []TemplateItem{
				{ItemID: "item1", DefaultQuantity: 2},
				{ItemID: "item2", DefaultQuantity: 1},
			},
			wantErr: false,
		},
		{
			name:         "empty name",
			templateName: "",
			description:  "Some description",
			items:        []TemplateItem{},
			wantErr:      true,
			errMsg:       "template name is required",
		},
		{
			name:         "no items",
			templateName: "Empty Template",
			description:  "Template with no items",
			items:        []TemplateItem{},
			wantErr:      true,
			errMsg:       "template must have at least one item",
		},
		{
			name:         "nil items",
			templateName: "Nil Items Template",
			description:  "Template with nil items",
			items:        nil,
			wantErr:      true,
			errMsg:       "template must have at least one item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template, err := NewJobTemplate(tt.templateName, tt.description, tt.items)

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

				if template == nil {
					t.Fatal("expected template but got nil")
				}

				if template.Name != tt.templateName {
					t.Errorf("expected name %q but got %q", tt.templateName, template.Name)
				}

				if template.Description != tt.description {
					t.Errorf("expected description %q but got %q", tt.description, template.Description)
				}

				if len(template.Items) != len(tt.items) {
					t.Errorf("expected %d items but got %d", len(tt.items), len(template.Items))
				}

				if template.ID == "" {
					t.Error("expected ID to be generated")
				}

				if !template.IsActive {
					t.Error("expected template to be active by default")
				}
			}
		})
	}
}

func TestJobTemplate_AddItem(t *testing.T) {
	template, _ := NewJobTemplate("Test Template", "Description", []TemplateItem{
		{ItemID: "item1", DefaultQuantity: 1},
	})

	newItem := TemplateItem{
		ItemID:          "item2",
		DefaultQuantity: 3,
	}

	template.AddItem(newItem)

	if len(template.Items) != 2 {
		t.Errorf("expected 2 items but got %d", len(template.Items))
	}

	lastItem := template.Items[len(template.Items)-1]
	if lastItem.ItemID != newItem.ItemID {
		t.Errorf("expected item ID %q but got %q", newItem.ItemID, lastItem.ItemID)
	}
}

func TestJobTemplate_RemoveItem(t *testing.T) {
	template, _ := NewJobTemplate("Test Template", "Description", []TemplateItem{
		{ItemID: "item1", DefaultQuantity: 1},
		{ItemID: "item2", DefaultQuantity: 2},
		{ItemID: "item3", DefaultQuantity: 3},
	})

	tests := []struct {
		name    string
		itemID  string
		wantErr bool
		errMsg  string
		expectedCount int
	}{
		{
			name:          "remove existing item",
			itemID:        "item2",
			wantErr:       false,
			expectedCount: 2,
		},
		{
			name:          "remove non-existent item",
			itemID:        "item999",
			wantErr:       true,
			errMsg:        "item not found in template",
			expectedCount: 2, // Count should remain same
		},
		{
			name:          "cannot remove last item",
			itemID:        "item1",
			wantErr:       true,
			errMsg:        "cannot remove last item from template",
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		if tt.name == "cannot remove last item" {
			// Reset to single item for this test
			template.Items = []TemplateItem{{ItemID: "item1", DefaultQuantity: 1}}
		}
		
		err := template.RemoveItem(tt.itemID)

		if tt.wantErr {
			if err == nil {
				t.Errorf("%s: expected error but got none", tt.name)
			} else if err.Error() != tt.errMsg {
				t.Errorf("%s: expected error message %q but got %q", tt.name, tt.errMsg, err.Error())
			}
		} else {
			if err != nil {
				t.Errorf("%s: unexpected error: %v", tt.name, err)
			}
		}

		if len(template.Items) != tt.expectedCount {
			t.Errorf("%s: expected %d items but got %d", tt.name, tt.expectedCount, len(template.Items))
		}
	}
}

func TestJobTemplate_Deactivate(t *testing.T) {
	template, _ := NewJobTemplate("Test Template", "Description", []TemplateItem{
		{ItemID: "item1", DefaultQuantity: 1},
	})

	if !template.IsActive {
		t.Error("expected template to be active initially")
	}

	template.Deactivate()

	if template.IsActive {
		t.Error("expected template to be inactive after deactivation")
	}
}

func TestJobTemplate_Activate(t *testing.T) {
	template, _ := NewJobTemplate("Test Template", "Description", []TemplateItem{
		{ItemID: "item1", DefaultQuantity: 1},
	})

	template.Deactivate()
	
	if template.IsActive {
		t.Error("expected template to be inactive after deactivation")
	}

	template.Activate()

	if !template.IsActive {
		t.Error("expected template to be active after activation")
	}
}