package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/masterbrent/electrical-bidding-app/internal/models"
)

// CompanyRepository defines the interface for company operations
type CompanyRepository interface {
	Get(ctx context.Context) (*models.Company, error)
	Update(ctx context.Context, company *models.Company) error
}

// companyRepository implements CompanyRepository
type companyRepository struct {
	db *sql.DB
}

// NewCompanyRepository creates a new company repository
func NewCompanyRepository(db *sql.DB) CompanyRepository {
	return &companyRepository{db: db}
}

// Get retrieves the company settings (always returns the single company record)
func (r *companyRepository) Get(ctx context.Context) (*models.Company, error) {
	query := `
		SELECT id, name, logo, address, city, state, zip, phone, email, license, website, created_at, updated_at
		FROM companies
		WHERE id = 'default'
		LIMIT 1
	`

	company := &models.Company{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&company.ID,
		&company.Name,
		&company.Logo,
		&company.Address,
		&company.City,
		&company.State,
		&company.Zip,
		&company.Phone,
		&company.Email,
		&company.License,
		&company.Website,
		&company.CreatedAt,
		&company.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("company settings not found")
		}
		return nil, err
	}

	return company, nil
}

// Update updates the company settings
func (r *companyRepository) Update(ctx context.Context, company *models.Company) error {
	query := `
		UPDATE companies
		SET name = $1, logo = $2, address = $3, city = $4, state = $5, 
		    zip = $6, phone = $7, email = $8, license = $9, website = $10,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = 'default'
	`

	result, err := r.db.ExecContext(ctx, query,
		company.Name,
		company.Logo,
		company.Address,
		company.City,
		company.State,
		company.Zip,
		company.Phone,
		company.Email,
		company.License,
		company.Website,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("company settings not found")
	}

	return nil
}