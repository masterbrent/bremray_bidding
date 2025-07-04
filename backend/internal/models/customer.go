package models

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
)

// Customer represents a customer who can have jobs
type Customer struct {
	ID    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Phone string `json:"phone" db:"phone"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewCustomer creates a new Customer with validation
func NewCustomer(name, email, phone string) (*Customer, error) {
	if name == "" {
		return nil, errors.New("customer name is required")
	}
	if email == "" {
		return nil, errors.New("customer email is required")
	}
	if !emailRegex.MatchString(email) {
		return nil, errors.New("invalid email format")
	}

	return &Customer{
		ID:    uuid.New().String(),
		Name:  name,
		Email: email,
		Phone: phone,
	}, nil
}