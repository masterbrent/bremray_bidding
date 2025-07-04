package models

import (
	"errors"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

// Company represents the company settings
type Company struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Logo      *string   `json:"logo" db:"logo"`
	Address   string    `json:"address" db:"address"`
	City      string    `json:"city" db:"city"`
	State     string    `json:"state" db:"state"`
	Zip       string    `json:"zip" db:"zip"`
	Phone     string    `json:"phone" db:"phone"`
	Email     string    `json:"email" db:"email"`
	License   string    `json:"license" db:"license"`
	Website   string    `json:"website" db:"website"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// NewCompany creates a new company instance
func NewCompany(name, email string) (*Company, error) {
	if err := validateCompany(name, email); err != nil {
		return nil, err
	}

	return &Company{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// Update updates company fields
func (c *Company) Update(name, email, address, city, state, zip, phone, license, website string) error {
	if name != "" {
		c.Name = name
	}
	
	if email != "" {
		if _, err := mail.ParseAddress(email); err != nil {
			return errors.New("invalid email format")
		}
		c.Email = email
	}
	
	if address != "" {
		c.Address = address
	}
	
	if city != "" {
		c.City = city
	}
	
	if state != "" {
		c.State = state
	}
	
	if zip != "" {
		c.Zip = zip
	}
	
	if phone != "" {
		c.Phone = phone
	}
	
	if license != "" {
		c.License = license
	}
	
	if website != "" {
		c.Website = website
	}
	
	c.UpdatedAt = time.Now()
	return nil
}

// SetLogo sets the company logo URL
func (c *Company) SetLogo(logoURL string) {
	c.Logo = &logoURL
	c.UpdatedAt = time.Now()
}

// RemoveLogo removes the company logo
func (c *Company) RemoveLogo() {
	c.Logo = nil
	c.UpdatedAt = time.Now()
}

// Validate validates the company data
func (c *Company) Validate() error {
	return validateCompany(c.Name, c.Email)
}

// validateCompany validates company fields
func validateCompany(name, email string) error {
	if name == "" {
		return errors.New("company name is required")
	}
	
	if email == "" {
		return errors.New("company email is required")
	}
	
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("invalid email format")
	}
	
	return nil
}