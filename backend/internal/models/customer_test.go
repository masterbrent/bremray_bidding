package models

import (
	"testing"
)

func TestNewCustomer(t *testing.T) {
	tests := []struct {
		name      string
		custName  string
		email     string
		phone     string
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "valid customer",
			custName: "John Doe",
			email:    "john@example.com",
			phone:    "555-1234",
			wantErr:  false,
		},
		{
			name:     "empty name",
			custName: "",
			email:    "john@example.com",
			phone:    "555-1234",
			wantErr:  true,
			errMsg:   "customer name is required",
		},
		{
			name:     "empty email",
			custName: "John Doe",
			email:    "",
			phone:    "555-1234",
			wantErr:  true,
			errMsg:   "customer email is required",
		},
		{
			name:     "invalid email",
			custName: "John Doe",
			email:    "not-an-email",
			phone:    "555-1234",
			wantErr:  true,
			errMsg:   "invalid email format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer, err := NewCustomer(tt.custName, tt.email, tt.phone)

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

				if customer == nil {
					t.Fatal("expected customer but got nil")
				}

				if customer.Name != tt.custName {
					t.Errorf("expected name %q but got %q", tt.custName, customer.Name)
				}

				if customer.Email != tt.email {
					t.Errorf("expected email %q but got %q", tt.email, customer.Email)
				}

				if customer.Phone != tt.phone {
					t.Errorf("expected phone %q but got %q", tt.phone, customer.Phone)
				}

				if customer.ID == "" {
					t.Error("expected ID to be generated")
				}
			}
		})
	}
}