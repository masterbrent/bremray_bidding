package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/masterbrent/electrical-bidding-app/internal/models"
)

// CustomerRepository defines the interface for customer database operations
type CustomerRepository interface {
	Create(ctx context.Context, customer *models.Customer) error
	GetByID(ctx context.Context, id string) (*models.Customer, error)
	GetByEmail(ctx context.Context, email string) (*models.Customer, error)
	Update(ctx context.Context, customer *models.Customer) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*models.Customer, error)
}

type customerRepository struct {
	db *sql.DB
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *models.Customer) error {
	query := `
		INSERT INTO customers (id, name, email, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		customer.ID, customer.Name, customer.Email, customer.Phone,
	)
	
	return err
}

func (r *customerRepository) GetByID(ctx context.Context, id string) (*models.Customer, error) {
	query := `
		SELECT id, name, email, phone
		FROM customers
		WHERE id = $1
	`
	
	customer := &models.Customer{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&customer.ID, &customer.Name, &customer.Email, &customer.Phone,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("customer not found")
	}
	
	return customer, err
}

func (r *customerRepository) GetByEmail(ctx context.Context, email string) (*models.Customer, error) {
	query := `
		SELECT id, name, email, phone
		FROM customers
		WHERE email = $1
	`
	
	customer := &models.Customer{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&customer.ID, &customer.Name, &customer.Email, &customer.Phone,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("customer not found")
	}
	
	return customer, err
}

func (r *customerRepository) Update(ctx context.Context, customer *models.Customer) error {
	query := `
		UPDATE customers SET
			name = $2, email = $3, phone = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	result, err := r.db.ExecContext(ctx, query,
		customer.ID, customer.Name, customer.Email, customer.Phone,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("customer not found")
	}
	
	return nil
}

func (r *customerRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM customers WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("customer not found")
	}
	
	return nil
}

func (r *customerRepository) List(ctx context.Context, limit, offset int) ([]*models.Customer, error) {
	query := `
		SELECT id, name, email, phone
		FROM customers
		ORDER BY name
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	customers := make([]*models.Customer, 0)
	for rows.Next() {
		customer := &models.Customer{}
		err := rows.Scan(
			&customer.ID, &customer.Name, &customer.Email, &customer.Phone,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	
	return customers, nil
}