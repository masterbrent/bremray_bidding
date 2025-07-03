package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/masterbrent/electrical-bidding-app/internal/models"
)

type itemRepository struct {
	db *sql.DB
}

// NewItemRepository creates a new item repository
func NewItemRepository(db *sql.DB) ItemRepository {
	return &itemRepository{db: db}
}

// Create inserts a new item into the database
func (r *itemRepository) Create(ctx context.Context, item *models.Item) error {
	if item == nil {
		return ErrInvalidInput
	}

	query := `
		INSERT INTO items (id, name, unit, unit_price, category, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		item.ID,
		item.Name,
		item.Unit,
		item.UnitPrice,
		item.Category,
		item.CreatedAt,
		item.UpdatedAt,
	)

	if err != nil {
		// Check for duplicate key error
		if isDuplicateKeyError(err) {
			return ErrDuplicate
		}
		return NewRepositoryError("failed to create item", err)
	}

	return nil
}

// GetByID retrieves an item by its ID
func (r *itemRepository) GetByID(ctx context.Context, id string) (*models.Item, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	query := `
		SELECT id, name, unit, unit_price, category, created_at, updated_at
		FROM items
		WHERE id = $1
	`

	var item models.Item
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Unit,
		&item.UnitPrice,
		&item.Category,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, NewRepositoryError("failed to get item", err)
	}

	return &item, nil
}

// List retrieves items with optional filtering
func (r *itemRepository) List(ctx context.Context, filter map[string]interface{}) ([]*models.Item, error) {
	query := `SELECT id, name, unit, unit_price, category, created_at, updated_at FROM items`
	args := []interface{}{}
	
	// Build WHERE clause from filter
	whereConditions := []string{}
	argIndex := 1
	
	if filter != nil {
		if category, ok := filter["category"].(string); ok && category != "" {
			whereConditions = append(whereConditions, fmt.Sprintf("category = $%d", argIndex))
			args = append(args, category)
			argIndex++
		}
	}
	
	if len(whereConditions) > 0 {
		query += " WHERE " + whereConditions[0]
		for i := 1; i < len(whereConditions); i++ {
			query += " AND " + whereConditions[i]
		}
	}
	
	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, NewRepositoryError("failed to list items", err)
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Unit,
			&item.UnitPrice,
			&item.Category,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, NewRepositoryError("failed to scan item", err)
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, NewRepositoryError("error iterating items", err)
	}

	return items, nil
}

// Update updates an existing item
func (r *itemRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	if id == "" {
		return ErrInvalidInput
	}

	// First, get the existing item
	item, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Apply updates to the model (with validation)
	if err := item.Update(updates); err != nil {
		return NewRepositoryError("validation failed", err)
	}

	// Update in database
	query := `
		UPDATE items
		SET name = $2, unit = $3, unit_price = $4, category = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		item.ID,
		item.Name,
		item.Unit,
		item.UnitPrice,
		item.Category,
		item.UpdatedAt,
	)

	if err != nil {
		return NewRepositoryError("failed to update item", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return NewRepositoryError("failed to check update result", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete removes an item from the database
func (r *itemRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidInput
	}

	query := `DELETE FROM items WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return NewRepositoryError("failed to delete item", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return NewRepositoryError("failed to check delete result", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// Helper function to check for duplicate key errors
func isDuplicateKeyError(err error) bool {
	// This is PostgreSQL specific - you might need to adjust for other databases
	return err != nil && err.Error() == "pq: duplicate key value violates unique constraint"
}