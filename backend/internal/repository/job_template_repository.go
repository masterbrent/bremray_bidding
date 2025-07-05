package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/masterbrent/electrical-bidding-app/internal/models"
)

// JobTemplateRepository defines the interface for job template database operations
type JobTemplateRepository interface {
	Create(ctx context.Context, template *models.JobTemplate) error
	GetByID(ctx context.Context, id string) (*models.JobTemplate, error)
	Update(ctx context.Context, template *models.JobTemplate) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*models.JobTemplate, error)
	ListActive(ctx context.Context) ([]*models.JobTemplate, error)
	
	// Template items operations
	AddTemplateItem(ctx context.Context, templateID string, item *models.TemplateItem) error
	GetTemplateItems(ctx context.Context, templateID string) ([]models.TemplateItem, error)
	UpdateTemplateItem(ctx context.Context, item *models.TemplateItem) error
	RemoveTemplateItem(ctx context.Context, templateID, itemID string) error
}

type jobTemplateRepository struct {
	db *sql.DB
}

// NewJobTemplateRepository creates a new job template repository
func NewJobTemplateRepository(db *sql.DB) JobTemplateRepository {
	return &jobTemplateRepository{db: db}
}

func (r *jobTemplateRepository) Create(ctx context.Context, template *models.JobTemplate) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Insert template
	query := `
		INSERT INTO job_templates (id, name, description, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	_, err = tx.ExecContext(ctx, query,
		template.ID, template.Name, template.Description,
		template.IsActive, template.CreatedAt, template.UpdatedAt,
	)
	if err != nil {
		return err
	}
	
	// Insert template items
	for _, item := range template.Items {
		itemQuery := `
			INSERT INTO template_items (id, template_id, item_id, default_quantity, created_at, updated_at)
			VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`
		_, err = tx.ExecContext(ctx, itemQuery,
			item.ID, template.ID, item.ItemID, item.DefaultQuantity,
		)
		if err != nil {
			return err
		}
	}
	
	// Insert template phases
	for _, phase := range template.Phases {
		phaseQuery := `
			INSERT INTO template_phases (id, template_id, name, phase_order, description)
			VALUES ($1, $2, $3, $4, $5)
		`
		_, err = tx.ExecContext(ctx, phaseQuery,
			phase.ID, template.ID, phase.Name, phase.Order, phase.Description,
		)
		if err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

func (r *jobTemplateRepository) GetByID(ctx context.Context, id string) (*models.JobTemplate, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM job_templates
		WHERE id = $1
	`
	
	template := &models.JobTemplate{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&template.ID, &template.Name, &template.Description,
		&template.IsActive, &template.CreatedAt, &template.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("template not found")
	}
	if err != nil {
		return nil, err
	}
	
	// Load template items
	template.Items, err = r.GetTemplateItems(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Load template phases
	template.Phases, err = r.GetTemplatePhases(ctx, id)
	if err != nil {
		log.Printf("Error loading phases for template %s: %v", id, err)
		return nil, err
	}
	log.Printf("Loaded %d phases for template %s", len(template.Phases), id)
	
	return template, nil
}

func (r *jobTemplateRepository) Update(ctx context.Context, template *models.JobTemplate) error {
	log.Printf("Repository Update called for template %s, phases: %d", template.ID, len(template.Phases))
	
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Update template
	query := `
		UPDATE job_templates SET
			name = $2, description = $3, is_active = $4, updated_at = $5
		WHERE id = $1
	`
	
	result, err := tx.ExecContext(ctx, query,
		template.ID, template.Name, template.Description,
		template.IsActive, template.UpdatedAt,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("template not found")
	}
	
	// Always update phases when Phases field is not nil (even if empty)
	if template.Phases != nil {
		log.Printf("Updating phases for template %s (count: %d)", template.ID, len(template.Phases))
		
		// Delete existing phases
		_, err = tx.ExecContext(ctx, `DELETE FROM template_phases WHERE template_id = $1`, template.ID)
		if err != nil {
			log.Printf("Error deleting phases: %v", err)
			return err
		}
		
		// Insert new phases if any
		for _, phase := range template.Phases {
			insertQuery := `
				INSERT INTO template_phases (id, template_id, name, phase_order, description)
				VALUES ($1, $2, $3, $4, $5)
			`
			log.Printf("Inserting phase: ID=%s, Name=%s, Order=%d", phase.ID, phase.Name, phase.Order)
			_, err = tx.ExecContext(ctx, insertQuery,
				phase.ID, template.ID, phase.Name, phase.Order, phase.Description,
			)
			if err != nil {
				log.Printf("Error inserting phase: %v", err)
				return err
			}
		}
	} else {
		log.Printf("Phases field is nil for template %s, not updating phases", template.ID)
	}
	
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}
	
	log.Printf("Template %s updated successfully", template.ID)
	return nil
}

func (r *jobTemplateRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM job_templates WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("template not found")
	}
	
	return nil
}

func (r *jobTemplateRepository) List(ctx context.Context, limit, offset int) ([]*models.JobTemplate, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM job_templates
		ORDER BY name
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanTemplates(ctx, rows)
}

func (r *jobTemplateRepository) ListActive(ctx context.Context) ([]*models.JobTemplate, error) {
	query := `
		SELECT id, name, description, is_active, created_at, updated_at
		FROM job_templates
		WHERE is_active = true
		ORDER BY name
	`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return r.scanTemplates(ctx, rows)
}

// Template Items operations
func (r *jobTemplateRepository) AddTemplateItem(ctx context.Context, templateID string, item *models.TemplateItem) error {
	query := `
		INSERT INTO template_items (id, template_id, item_id, default_quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		item.ID, templateID, item.ItemID, item.DefaultQuantity,
	)
	
	return err
}

func (r *jobTemplateRepository) GetTemplateItems(ctx context.Context, templateID string) ([]models.TemplateItem, error) {
	query := `
		SELECT id, template_id, item_id, default_quantity
		FROM template_items
		WHERE template_id = $1
		ORDER BY id
	`
	
	rows, err := r.db.QueryContext(ctx, query, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	items := make([]models.TemplateItem, 0)
	for rows.Next() {
		var item models.TemplateItem
		err := rows.Scan(
			&item.ID, &item.TemplateID, &item.ItemID, &item.DefaultQuantity,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	
	return items, nil
}

func (r *jobTemplateRepository) UpdateTemplateItem(ctx context.Context, item *models.TemplateItem) error {
	query := `
		UPDATE template_items SET
			default_quantity = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	
	result, err := r.db.ExecContext(ctx, query, item.ID, item.DefaultQuantity)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("template item not found")
	}
	
	return nil
}

func (r *jobTemplateRepository) RemoveTemplateItem(ctx context.Context, templateID, itemID string) error {
	query := `DELETE FROM template_items WHERE template_id = $1 AND item_id = $2`
	
	result, err := r.db.ExecContext(ctx, query, templateID, itemID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("template item not found")
	}
	
	return nil
}

// Helper function to scan template rows
func (r *jobTemplateRepository) scanTemplates(ctx context.Context, rows *sql.Rows) ([]*models.JobTemplate, error) {
	templates := make([]*models.JobTemplate, 0)
	for rows.Next() {
		template := &models.JobTemplate{}
		err := rows.Scan(
			&template.ID, &template.Name, &template.Description,
			&template.IsActive, &template.CreatedAt, &template.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		// Load template items
		template.Items, err = r.GetTemplateItems(ctx, template.ID)
		if err != nil {
			return nil, err
		}
		
		// Load template phases
		template.Phases, err = r.GetTemplatePhases(ctx, template.ID)
		if err != nil {
			return nil, err
		}
		
		templates = append(templates, template)
	}
	
	return templates, nil
}
// GetTemplatePhases retrieves all phases for a template
func (r *jobTemplateRepository) GetTemplatePhases(ctx context.Context, templateID string) ([]models.TemplatePhase, error) {
	log.Printf("GetTemplatePhases called for template %s", templateID)
	
	query := `
		SELECT id, template_id, name, phase_order, description
		FROM template_phases
		WHERE template_id = $1
		ORDER BY phase_order
	`
	
	rows, err := r.db.QueryContext(ctx, query, templateID)
	if err != nil {
		log.Printf("Error querying template phases: %v", err)
		return nil, err
	}
	defer rows.Close()
	
	phases := make([]models.TemplatePhase, 0)
	for rows.Next() {
		var phase models.TemplatePhase
		var description sql.NullString
		err := rows.Scan(
			&phase.ID, &phase.TemplateID, &phase.Name, 
			&phase.Order, &description,
		)
		if err != nil {
			log.Printf("Error scanning phase row: %v", err)
			return nil, err
		}
		if description.Valid {
			phase.Description = description.String
		}
		phases = append(phases, phase)
	}
	
	log.Printf("Found %d phases for template %s", len(phases), templateID)
	return phases, nil
}
