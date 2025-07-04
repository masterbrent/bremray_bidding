package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCompany(t *testing.T) {
	tests := []struct {
		name      string
		compName  string
		email     string
		expectErr bool
	}{
		{
			name:      "Valid company",
			compName:  "Bremray Electrical",
			email:     "info@bremray.com",
			expectErr: false,
		},
		{
			name:      "Empty name",
			compName:  "",
			email:     "info@bremray.com",
			expectErr: true,
		},
		{
			name:      "Empty email",
			compName:  "Bremray Electrical",
			email:     "",
			expectErr: true,
		},
		{
			name:      "Invalid email",
			compName:  "Bremray Electrical",
			email:     "invalid-email",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			company, err := NewCompany(tt.compName, tt.email)
			
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, company)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, company)
				assert.NotEmpty(t, company.ID)
				assert.Equal(t, tt.compName, company.Name)
				assert.Equal(t, tt.email, company.Email)
				assert.NotZero(t, company.CreatedAt)
				assert.NotZero(t, company.UpdatedAt)
			}
		})
	}
}

func TestCompany_Update(t *testing.T) {
	company, err := NewCompany("Bremray Electrical", "info@bremray.com")
	assert.NoError(t, err)

	originalUpdatedAt := company.UpdatedAt

	// Update with valid data
	err = company.Update(
		"Bremray Electrical LLC",
		"contact@bremray.com",
		"123 Main St",
		"Anytown",
		"CA",
		"12345",
		"555-123-4567",
		"LIC123456",
		"https://bremray.com",
	)
	
	assert.NoError(t, err)
	assert.Equal(t, "Bremray Electrical LLC", company.Name)
	assert.Equal(t, "contact@bremray.com", company.Email)
	assert.Equal(t, "123 Main St", company.Address)
	assert.Equal(t, "Anytown", company.City)
	assert.Equal(t, "CA", company.State)
	assert.Equal(t, "12345", company.Zip)
	assert.Equal(t, "555-123-4567", company.Phone)
	assert.Equal(t, "LIC123456", company.License)
	assert.Equal(t, "https://bremray.com", company.Website)
	assert.NotEqual(t, originalUpdatedAt, company.UpdatedAt)

	// Update with invalid email
	err = company.Update("", "invalid-email", "", "", "", "", "", "", "")
	assert.Error(t, err)
}

func TestCompany_Logo(t *testing.T) {
	company, err := NewCompany("Bremray Electrical", "info@bremray.com")
	assert.NoError(t, err)
	assert.Nil(t, company.Logo)

	// Set logo
	logoURL := "https://example.com/logo.png"
	company.SetLogo(logoURL)
	assert.NotNil(t, company.Logo)
	assert.Equal(t, logoURL, *company.Logo)

	// Remove logo
	company.RemoveLogo()
	assert.Nil(t, company.Logo)
}

func TestCompany_Validate(t *testing.T) {
	tests := []struct {
		name      string
		company   *Company
		expectErr bool
	}{
		{
			name: "Valid company",
			company: &Company{
				Name:  "Bremray Electrical",
				Email: "info@bremray.com",
			},
			expectErr: false,
		},
		{
			name: "Empty name",
			company: &Company{
				Name:  "",
				Email: "info@bremray.com",
			},
			expectErr: true,
		},
		{
			name: "Invalid email",
			company: &Company{
				Name:  "Bremray Electrical",
				Email: "invalid",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.company.Validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}